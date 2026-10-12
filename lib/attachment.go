package lib

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

// Attachment is a file to send with a request, with the file name and media type to send it under.
// Assign it to a request's file parameter, such as BundleCreateParams.WatermarkAttachmentFile. It is an
// io.Writer, so content can be copied into it, or it can be built with NewAttachment.
type Attachment struct {
	bytes.Buffer
	FileName    string `json:"-" url:"-"`
	ContentType string `json:"-" url:"-"`
}

// NewAttachment returns an Attachment holding a copy of content.
func NewAttachment(fileName string, contentType string, content []byte) *Attachment {
	attachment := &Attachment{FileName: fileName, ContentType: contentType}
	attachment.Write(content)
	return attachment
}

// MultipartBody encodes the request body as multipart/form-data when the parameters set a file parameter,
// a field declared as io.Writer. Each file is sent as a part holding its bytes, and every other parameter
// as a form field encoded as it would be in a query string. The body is held in memory, so the request
// can be sent again on retry. When no file parameter is set it returns a nil body and the request is sent
// as JSON as before.
//
// A file must also be readable, such as an *Attachment, *bytes.Buffer or *os.File. Content is read from
// its current position, a buffer's unread bytes are not consumed, and the value is never closed.
func MultipartBody(values Values) (io.Reader, string, error) {
	var params Params
	switch v := values.(type) {
	case Params:
		params = v
	case *Params:
		if v == nil {
			return nil, "", nil
		}
		params = *v
	default:
		return nil, "", nil
	}

	files := fileFields(params.Params)
	if len(files) == 0 {
		return nil, "", nil
	}
	form, err := params.ToValues()
	if err != nil {
		return nil, "", err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	keys := make([]string, 0, len(form))
	for key := range form {
		if !isFileKey(key, files) {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		for _, value := range form[key] {
			if err := writer.WriteField(key, value); err != nil {
				return nil, "", err
			}
		}
	}
	for _, file := range files {
		if err := writeFilePart(writer, file); err != nil {
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return bytes.NewReader(body.Bytes()), writer.FormDataContentType(), nil
}

var writerType = reflect.TypeOf((*io.Writer)(nil)).Elem()

type fileField struct {
	name  string
	value interface{}
}

// fileFields returns the set file parameters: the request struct's io.Writer fields, which is how the
// generator declares a schema file parameter.
func fileFields(params interface{}) []fileField {
	v := reflect.ValueOf(params)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	var files []fileField
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)
		if !field.IsExported() || field.Type != writerType || value.IsNil() {
			continue
		}
		if elem := value.Elem(); elem.Kind() == reflect.Pointer && elem.IsNil() {
			continue
		}
		files = append(files, fileField{name: formName(field), value: value.Interface()})
	}
	return files
}

func formName(field reflect.StructField) string {
	for _, key := range []string{"json", "url"} {
		if name, _, _ := strings.Cut(field.Tag.Get(key), ","); name != "" && name != "-" {
			return name
		}
	}
	return field.Name
}

func isFileKey(key string, files []fileField) bool {
	for _, file := range files {
		if key == file.name || strings.HasPrefix(key, file.name+"[") {
			return true
		}
	}
	return false
}

func writeFilePart(writer *multipart.Writer, file fileField) error {
	content, err := fileContent(file)
	if err != nil {
		return err
	}

	fileName, contentType := file.name, "application/octet-stream"
	if named, ok := file.value.(interface{ Name() string }); ok && named.Name() != "" {
		fileName = filepath.Base(named.Name())
	}
	if attachment, ok := file.value.(*Attachment); ok {
		if attachment.FileName != "" {
			fileName = attachment.FileName
		}
		if attachment.ContentType != "" {
			contentType = attachment.ContentType
		}
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, quoteEscaper.Replace(file.name), quoteEscaper.Replace(fileName)))
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}
	_, err = part.Write(content)
	return err
}

func fileContent(file fileField) ([]byte, error) {
	switch value := file.value.(type) {
	case interface{ Bytes() []byte }:
		return value.Bytes(), nil
	case io.Reader:
		return io.ReadAll(value)
	}
	return nil, fmt.Errorf("%s: the file is an io.Writer that cannot be read; use a *lib.Attachment, *bytes.Buffer or *os.File", file.name)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// ErrAttachmentTooLarge is returned, wrapped, when the reader of an attachment
// made with NewJSONAttachmentUpTo has more bytes than its limit.
var ErrAttachmentTooLarge = errors.New("attachment is larger than its limit")

// ErrAttachmentSizeMismatch is returned, wrapped, when the reader of an
// attachment made with NewJSONAttachment ends before or continues after its size.
var ErrAttachmentSizeMismatch = errors.New("attachment does not have its declared size")

// JSONAttachment is file content for a request parameter that takes a file, such
// as files_sdk.UserUpdateParams.AvatarFileAttachment. Unlike an Attachment,
// which a parameter's io.Writer field sends as a multipart/form-data part, the
// request sends it in the same JSON body as the other parameters, as a file
// object:
//
//	{"encoded_content": "<base64 of the content>", "filename": "...", "type": "..."}
//
// Before sending, the SDK reads the content once, from the reader's current
// position, and keeps it encoded: in memory, or in a temporary file in
// os.TempDir when its size or limit encodes to more than 1 MiB. The file is
// created with os.CreateTemp, which on Unix lets only its owner read it, and
// is removed when the call returns. If the HTTP transport is still sending
// the request body then, the file stays open until the transport closes the
// body, and on Windows, where an open file cannot be removed, it is removed
// only then. Retries resend the kept copy, so every attempt sends the same
// bytes. The SDK does not close or seek the reader, and a JSONAttachment can be
// sent only once.
//
// When the reader fails, the request's error names the parameter and wraps
// the reader's error without repeating its text, which can include a signed
// URL of the content's source. Use errors.Is, errors.As or errors.Unwrap to
// inspect the reader's error.
//
// Reading stops with the request context's error (see files_sdk.WithContext)
// at the next read after the context is done. A Read that is already blocked
// is not interrupted; to stop one, close or unblock the reader.
//
// A JSONAttachment can also be a top-level value of the map given to an
// UpdateWithMap method, keyed by the parameter name.
type JSONAttachment struct {
	reader      io.Reader
	limit       int64
	exact       bool
	filename    string
	contentType string
	sent        atomic.Bool
}

// NewJSONAttachment returns an attachment of exactly size bytes read from r, such
// as an *os.File and its Stat size. The request fails, before anything is
// sent, if r ends before size bytes or has more. filename and contentType,
// such as "avatar.png" and "image/png", are sent as given; the server
// validates them.
func NewJSONAttachment(r io.Reader, size int64, filename string, contentType string) *JSONAttachment {
	return &JSONAttachment{reader: r, limit: size, exact: true, filename: filename, contentType: contentType}
}

// NewJSONAttachmentUpTo returns an attachment of everything r has, up to maxSize
// bytes, for a reader whose length is not known in advance, such as a
// response body. The request fails with ErrAttachmentTooLarge, before
// anything is sent, if r has more than maxSize bytes. To find that out the
// SDK reads at most one byte past maxSize.
func NewJSONAttachmentUpTo(r io.Reader, maxSize int64, filename string, contentType string) *JSONAttachment {
	return &JSONAttachment{reader: r, limit: maxSize, filename: filename, contentType: contentType}
}

// Content up to this encoded size is kept in memory; larger content is kept
// in a temporary file.
const attachmentMemoryLimit = 1 << 20

// The largest content whose base64 encoding fits in an int64.
const maxEncodableSize = math.MaxInt64 / 4 * 3

// Consecutive empty reads after which a reader is taken to be stuck, as
// bufio does.
const maxEmptyReads = 100

func encodedSize(size int64) int64 {
	return (size + 2) / 3 * 4
}

type namedAttachment struct {
	key string // JSON key of the file parameter
	// Go field for the parameter's io.Writer, such as AvatarFile for
	// AvatarFileAttachment, or "" for a map value.
	writerField string
	attachment  *JSONAttachment
}

func (a namedAttachment) errorf(format string, args ...any) error {
	return fmt.Errorf("attachment for %s: "+format, append([]any{a.key}, args...)...)
}

// check reports a problem with the attachment that is known before its reader
// is read, so that no reader is read when a request cannot be sent.
func (a namedAttachment) check() error {
	switch {
	case a.attachment.reader == nil:
		return a.errorf("reader is nil")
	case a.attachment.limit < 0:
		return a.errorf("size must not be negative")
	case a.attachment.limit > maxEncodableSize:
		return a.errorf("size %d is too large to encode", a.attachment.limit)
	case !utf8.ValidString(a.attachment.filename):
		// JSON text is UTF-8; encoding would replace the invalid bytes.
		return a.errorf("filename is not valid UTF-8")
	case !utf8.ValidString(a.attachment.contentType):
		return a.errorf("type is not valid UTF-8")
	case a.attachment.sent.Load():
		return a.alreadySentError()
	}
	return nil
}

func (a namedAttachment) alreadySentError() error {
	return a.errorf("already sent; make a new JSONAttachment for each request")
}

// readerError is a failure of an attachment's reader. Its message leaves out
// the reader's error, whose text can include a signed URL of the content's
// source; errors.Is and errors.As reach that error, err.
type readerError struct {
	key string
	err error
}

func (e *readerError) Error() string {
	return fmt.Sprintf("attachment for %s: reading: the reader returned an error", e.key)
}

func (e *readerError) Unwrap() error {
	return e.err
}

// findAttachments returns the attachments of a request struct, from its
// fields tagged `attachment`, or of a map, from its values, in key order.
func findAttachments(params interface{}) []namedAttachment {
	value := reflect.ValueOf(params)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	var attachments []namedAttachment
	switch value.Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			key, tagged := field.Tag.Lookup("attachment")
			if !tagged || !field.IsExported() {
				continue
			}
			if attachment, ok := value.Field(i).Interface().(*JSONAttachment); ok && attachment != nil {
				writerField := strings.TrimSuffix(field.Name, "Attachment")
				attachments = append(attachments, namedAttachment{key: key, writerField: writerField, attachment: attachment})
			}
		}
	case reflect.Map:
		if m, ok := value.Interface().(map[string]interface{}); ok {
			for key, v := range m {
				if attachment, ok := v.(*JSONAttachment); ok && attachment != nil {
					attachments = append(attachments, namedAttachment{key: key, attachment: attachment})
				}
			}
		}
	}
	sort.Slice(attachments, func(i, j int) bool { return attachments[i].key < attachments[j].key })
	return attachments
}

// HasJSONAttachments reports whether p sends any JSONAttachment. files_sdk.Call
// sends such a request with PrepareJSON instead of ToJSON.
func (p Params) HasJSONAttachments() bool {
	return len(findAttachments(p.Params)) > 0
}

// withoutAttachments returns p with a map's attachment values removed. A
// struct's attachment fields are already left out of its JSON and query
// values by their tags.
func (p Params) withoutAttachments() Params {
	m, ok := p.Params.(map[string]interface{})
	if !ok {
		return p
	}
	rest := make(map[string]interface{}, len(m))
	for key, v := range m {
		if _, isAttachment := v.(*JSONAttachment); !isAttachment {
			rest[key] = v
		}
	}
	return Params{Params: rest}
}

// PrepareJSON reads the attachments of p and returns the whole JSON request
// body: the other parameters, encoded as ToJSON would encode them, with a
// file object for each attachment. Nothing is read from any reader until the
// other parameters and every attachment have been checked. ctx stops the
// reading at the next read after it is done. The caller must Close the
// result.
func (p Params) PrepareJSON(ctx context.Context) (*PreparedJSON, error) {
	attachments := findAttachments(p.Params)
	companions, err := p.withoutAttachments().encodeJSON()
	if err != nil {
		return nil, err
	}
	members, err := jsonObjectMembers(companions)
	if err != nil {
		return nil, err
	}
	if err := checkAttachments(attachments, members, int64(len(companions))); err != nil {
		return nil, err
	}

	body := &PreparedJSON{}
	body.appendBytes([]byte("{"))
	body.appendBytes(companions[1 : len(companions)-1])
	for i, a := range attachments {
		if i > 0 || len(members) > 0 {
			body.appendBytes([]byte(","))
		}
		if err := body.appendAttachment(ctx, a); err != nil {
			// The reader's error comes first; the staged content is only
			// SDK-owned temporary data.
			_ = body.Close()
			return nil, err
		}
	}
	body.appendBytes([]byte("}"))
	return body, nil
}

func jsonObjectMembers(object []byte) (map[string]json.RawMessage, error) {
	if !bytes.HasPrefix(object, []byte("{")) {
		return nil, fmt.Errorf("request parameters must encode as a JSON object to send attachments")
	}
	var members map[string]json.RawMessage
	if err := json.Unmarshal(object, &members); err != nil {
		return nil, err
	}
	return members, nil
}

func checkAttachments(attachments []namedAttachment, members map[string]json.RawMessage, companionsLength int64) error {
	// The body is the other parameters plus at most the encoded limits.
	length := companionsLength
	for _, a := range attachments {
		if err := a.check(); err != nil {
			return err
		}
		// Only a struct can also set the key, through the parameter's io.Writer
		// field; a map key has one value.
		if _, set := members[a.key]; set {
			return fmt.Errorf("%[1]s and %[1]sAttachment are both set; leave %[1]s nil to send the attachment", a.writerField)
		}
		encoded := encodedSize(a.attachment.limit)
		if length > math.MaxInt64-encoded {
			return fmt.Errorf("attachments are too large to send in one request")
		}
		length += encoded
	}
	return nil
}

// PreparedJSON is a JSON request body that can be read any number of times,
// returned by Params.PrepareJSON for requests with attachments, and by
// CopyBody. Its content does not change.
type PreparedJSON struct {
	parts  []preparedPart
	length int64

	mu     sync.Mutex
	files  []*stagedFile
	open   int  // readers from Open not yet closed
	closed bool // Close was called
}

// stagedFile is a temporary file that holds part of a body.
type stagedFile struct {
	*os.File
	removed bool // its name is already removed; the open file is still readable
}

type preparedPart struct {
	data []byte
	file *os.File
	size int64
}

// Len returns the length of the body in bytes.
func (b *PreparedJSON) Len() int64 {
	return b.length
}

// NewReader returns a reader over the whole body. Each reader has its own
// position, so readers can be used at the same time.
func (b *PreparedJSON) NewReader() io.Reader {
	readers := make([]io.Reader, 0, len(b.parts))
	for _, part := range b.parts {
		if part.file != nil {
			readers = append(readers, io.NewSectionReader(part.file, 0, part.size))
		} else {
			readers = append(readers, bytes.NewReader(part.data))
		}
	}
	return io.MultiReader(readers...)
}

// Open returns a reader over the whole body, like NewReader, that keeps the
// body's temporary files open until it is closed, even after Close. It is for
// an HTTP transport, which can finish sending a request body and close it
// after the call that sent the request returns. Closing the reader more than
// once is safe.
func (b *PreparedJSON) Open() io.ReadCloser {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.open++
	return &openReader{Reader: b.NewReader(), body: b}
}

// Close removes the temporary files that hold the body. While readers from
// Open are still open, the files stay open for them: Close removes their
// names at once where the system allows that for an open file (not on
// Windows), and the last of the readers to close closes the files and removes
// any that remain. Readers from NewReader fail once the files are closed.
// Close can be called more than once.
func (b *PreparedJSON) Close() error {
	b.mu.Lock()
	b.closed = true
	if b.open > 0 {
		for _, file := range b.files {
			if !file.removed && os.Remove(file.Name()) == nil {
				file.removed = true
			}
		}
	}
	return b.releaseFilesLocked()
}

func (b *PreparedJSON) release() error {
	b.mu.Lock()
	b.open--
	return b.releaseFilesLocked()
}

// releaseFilesLocked closes the files, and removes those whose names remain,
// once they are neither needed nor held open; it unlocks b.mu.
func (b *PreparedJSON) releaseFilesLocked() error {
	var files []*stagedFile
	if b.closed && b.open == 0 {
		files, b.files = b.files, nil
	}
	b.mu.Unlock()
	var errs []error
	for _, file := range files {
		errs = append(errs, file.Close())
		if !file.removed {
			errs = append(errs, os.Remove(file.Name()))
		}
	}
	return errors.Join(errs...)
}

type openReader struct {
	io.Reader
	body *PreparedJSON
	once sync.Once
}

func (r *openReader) Close() error {
	var err error
	r.once.Do(func() { err = r.body.release() })
	return err
}

// CopyBody reads r to its end and returns a copy of what it read, kept as
// attachment content is: in memory up to 1 MiB, and otherwise in a temporary
// file that Close removes. It reads at most 32 KiB at a time, and ctx stops
// the reading at the next read after it is done. It returns r's error
// unchanged; it keeps nothing when it fails.
func CopyBody(ctx context.Context, r io.Reader) (*PreparedJSON, error) {
	body := &PreparedJSON{}
	var buffered *bufio.Writer
	var file *os.File
	var inMemory []byte
	buffer := make([]byte, 32<<10)
	emptyReads := 0
	fail := func(err error) (*PreparedJSON, error) {
		_ = body.Close()
		return nil, err
	}
	for {
		if err := ctx.Err(); err != nil {
			return fail(err)
		}
		n, err := r.Read(buffer)
		if n > 0 {
			emptyReads = 0
			if buffered == nil && len(inMemory)+n > attachmentMemoryLimit {
				// os.CreateTemp creates the file with mode 0600.
				created, createErr := os.CreateTemp("", "files-sdk-request-body-*")
				if createErr != nil {
					return fail(fmt.Errorf("copying the request body: creating a temporary file: %w", createErr))
				}
				file = created
				body.files = append(body.files, &stagedFile{File: file})
				buffered = bufio.NewWriterSize(file, 64<<10)
				if _, writeErr := buffered.Write(inMemory); writeErr != nil {
					return fail(fmt.Errorf("copying the request body: writing a temporary file: %w", writeErr))
				}
				inMemory = nil
			}
			if buffered != nil {
				if _, writeErr := buffered.Write(buffer[:n]); writeErr != nil {
					return fail(fmt.Errorf("copying the request body: writing a temporary file: %w", writeErr))
				}
			} else {
				inMemory = append(inMemory, buffer[:n]...)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fail(err)
		}
		if n == 0 {
			if emptyReads++; emptyReads == maxEmptyReads {
				return fail(io.ErrNoProgress)
			}
		}
	}
	if buffered == nil {
		body.appendBytes(inMemory)
		return body, nil
	}
	if err := buffered.Flush(); err != nil {
		return fail(fmt.Errorf("copying the request body: writing a temporary file: %w", err))
	}
	size, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return fail(fmt.Errorf("copying the request body: writing a temporary file: %w", err))
	}
	body.parts = append(body.parts, preparedPart{file: file, size: size})
	body.length = size
	return body, nil
}

func (b *PreparedJSON) appendBytes(data []byte) {
	if len(data) == 0 {
		return
	}
	b.length += int64(len(data))
	if last := len(b.parts) - 1; last >= 0 && b.parts[last].file == nil {
		b.parts[last].data = append(b.parts[last].data, data...)
		return
	}
	b.parts = append(b.parts, preparedPart{data: append([]byte(nil), data...)})
}

func (b *PreparedJSON) appendAttachment(ctx context.Context, a namedAttachment) error {
	key, err := json.Marshal(a.key)
	if err != nil {
		return err
	}
	filename, err := json.Marshal(a.attachment.filename)
	if err != nil {
		return err
	}
	contentType, err := json.Marshal(a.attachment.contentType)
	if err != nil {
		return err
	}
	b.appendBytes(key)
	b.appendBytes([]byte(`:{"encoded_content":"`))
	if err := b.appendEncodedContent(ctx, a); err != nil {
		return err
	}
	b.appendBytes([]byte(`","filename":`))
	b.appendBytes(filename)
	b.appendBytes([]byte(`,"type":`))
	b.appendBytes(contentType)
	b.appendBytes([]byte("}"))
	return nil
}

// appendEncodedContent reads the attachment and appends its base64 encoding,
// kept in memory or in a temporary file depending on its limit.
func (b *PreparedJSON) appendEncodedContent(ctx context.Context, a namedAttachment) error {
	if !a.attachment.sent.CompareAndSwap(false, true) {
		return a.alreadySentError()
	}
	if encodedSize(a.attachment.limit) <= attachmentMemoryLimit {
		var encoded bytes.Buffer
		if err := a.encode(ctx, &encoded); err != nil {
			return err
		}
		b.appendBytes(encoded.Bytes())
		return nil
	}

	// os.CreateTemp creates the file with mode 0600.
	file, err := os.CreateTemp("", "files-sdk-attachment-*")
	if err != nil {
		return a.errorf("creating a temporary file: %w", err)
	}
	b.files = append(b.files, &stagedFile{File: file})
	buffered := bufio.NewWriterSize(file, 64<<10)
	if err := a.encode(ctx, buffered); err != nil {
		return err
	}
	if err := buffered.Flush(); err != nil {
		return a.errorf("writing a temporary file: %w", err)
	}
	size, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return a.errorf("writing a temporary file: %w", err)
	}
	b.parts = append(b.parts, preparedPart{file: file, size: size})
	b.length += size
	return nil
}

// encode reads the attachment's content into dst as base64.
func (a namedAttachment) encode(ctx context.Context, dst io.Writer) error {
	encoder := base64.NewEncoder(base64.StdEncoding, dst)
	if err := a.copyContent(ctx, encoder); err != nil {
		return err
	}
	// Close writes the final partial block; it does not close dst.
	if err := encoder.Close(); err != nil {
		return a.errorf("encoding: %w", err)
	}
	return nil
}

// copyContent copies the reader to dst, reading at most one byte past the
// attachment's limit and checking ctx before every read.
func (a namedAttachment) copyContent(ctx context.Context, dst io.Writer) error {
	limit := a.attachment.limit
	buffer := make([]byte, 32<<10)
	var read int64
	emptyReads := 0
	for {
		if err := ctx.Err(); err != nil {
			return a.errorf("%w", err)
		}
		chunk := buffer
		if remaining := limit - read + 1; remaining < int64(len(chunk)) {
			chunk = chunk[:remaining]
		}
		n, err := a.attachment.reader.Read(chunk)
		if n > 0 {
			emptyReads = 0
			read += int64(n)
			if read > limit {
				if a.attachment.exact {
					return a.errorf("reader has more than its size of %d bytes: %w", limit, ErrAttachmentSizeMismatch)
				}
				return a.errorf("reader has more than %d bytes: %w", limit, ErrAttachmentTooLarge)
			}
			if _, err := dst.Write(chunk[:n]); err != nil {
				return a.errorf("writing: %w", err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return &readerError{key: a.key, err: err}
		}
		if n == 0 {
			if emptyReads++; emptyReads == maxEmptyReads {
				return a.errorf("reading: %w", io.ErrNoProgress)
			}
		}
	}
	if a.attachment.exact && read != limit {
		return a.errorf("reader ended after %d of its %d bytes: %w", read, limit, ErrAttachmentSizeMismatch)
	}
	return nil
}
