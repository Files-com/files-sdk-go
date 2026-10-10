package lib

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
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
