package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"sync/atomic"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/samber/lo"
	"github.com/zenthangplus/goccm"
)

const (
	DownloadPartChunkSize         = int64(1024 * 1024 * 5)
	DownloadPartLimit             = 15
	DownloadPartLimitForSizeGuess = 10
)

const (
	NotRun = iota
	RunSuccess
	RunError
)

type DownloadParts struct {
	globalWait goccm.ConcurrencyManager
	context.CancelFunc
	context.Context
	queueCancel  context.CancelFunc
	queueContext context.Context
	fs.File
	fs.FileInfo
	lib.WriterAndAt
	totalWritten         int64
	parts                []*Part
	queue                chan *Part
	finishedParts        chan *Part
	firstPartSuccessChan chan int
	firstPartSuccess     int
	runningPartLimit     int
	CloseError           error
	files_sdk.Config
}

func (d *DownloadParts) Init(file fs.File, info fs.FileInfo, globalWait goccm.ConcurrencyManager, writer lib.WriterAndAt, config files_sdk.Config) *DownloadParts {
	d.File = file
	d.FileInfo = info
	d.globalWait = globalWait
	d.WriterAndAt = writer
	d.setRunningPartLimit()
	d.firstPartSuccessChan = make(chan int, 1)
	d.firstPartSuccess = NotRun
	d.Config = config
	return d
}

func (d *DownloadParts) Run(ctx context.Context) error {
	d.Context, d.CancelFunc = context.WithCancel(ctx)
	d.queueContext, d.queueCancel = context.WithCancel(d.Context)
	defer func() {
		d.Config.LogPath(
			d.FileInfo.Name(),
			map[string]interface{}{
				"message":  "Finished canceling context and closing file",
				"realSize": d.totalWritten,
			},
		)
		d.CancelFunc()
		d.CloseError = d.WriterAndAt.Close()
	}()
	d.Config.LogPath(
		d.FileInfo.Name(),
		map[string]interface{}{
			"size":      d.FileInfo.Size(),
			"SizeGuess": d.SizeGuess(),
		},
	)
	if d.FileInfo.Size() <= DownloadPartChunkSize || d.SizeGuess() {
		return d.downloadFile()
	} else {
		d.buildParts()
		d.listenOnQueue()
		d.addPartsToQueue()
		return d.waitForParts()
	}
}

func (d *DownloadParts) FinalSize() int64 {
	return d.totalWritten
}

func (d *DownloadParts) setRunningPartLimit() {
	if d.SizeGuess() {
		d.runningPartLimit = DownloadPartLimitForSizeGuess
	} else {
		d.runningPartLimit = DownloadPartLimit
	}
}

func (d *DownloadParts) waitForParts() error {
	for i := range d.parts {
		part := <-d.finishedParts
		if part.error != nil {
			d.TriggerFirstSuccess(RunError)
			return part.error
		}

		d.TriggerFirstSuccess(RunSuccess)
		d.Config.LogPath(
			d.FileInfo.Name(),
			map[string]interface{}{
				"RunningParts":  d.runningParts(),
				"limit":         d.runningPartLimit,
				"parts":         len(d.parts),
				"Written":       d.totalWritten,
				"PartFinished":  part.number,
				"partBytes":     part.bytes,
				"PartsFinished": i + 1,
			},
		)
	}

	return d.realSizeOverLap()
}

func (d *DownloadParts) realSizeOverLap() error {
	lastPart := d.parts[len(d.parts)-1]
	d.Config.LogPath(
		d.FileInfo.Name(),
		map[string]interface{}{
			"message":  "starting realSizeOverLap",
			"size":     d.FileInfo.Size(),
			"realSize": d.totalWritten,
		},
	)
	defer func() {
		d.Config.LogPath(
			d.FileInfo.Name(),
			map[string]interface{}{
				"message":  "finishing realSizeOverLap",
				"size":     d.FileInfo.Size(),
				"realSize": d.totalWritten,
			},
		)
	}()
	for {
		d.Config.LogPath(
			d.FileInfo.Name(),
			map[string]interface{}{
				"message":      "finishing realSizeOverLap",
				"SizeGuess":    d.SizeGuess(),
				"queueStopped": d.queueContext.Err() != nil,
				"lastBytes":    lastPart.bytes,
				"lastLen":      lastPart.len,
			},
		)
		if d.SizeGuess() && d.queueContext.Err() == nil && lastPart.bytes != 0 {
			d.queue = make(chan *Part, 1)
			d.finishedParts = make(chan *Part, 1)
			nextPart := &Part{number: lastPart.number + 1, OffSet: OffSet{off: lastPart.off + lastPart.bytes, len: DownloadPartChunkSize}}
			d.Config.LogPath(d.FileInfo.Name(), map[string]interface{}{"message": "Next Part for size guess", "part": nextPart.number})
			d.parts = append(d.parts, nextPart)

			go d.processPart(nextPart.Start(d.Context), true)
			select {
			case lastPart = <-d.finishedParts:
				if lastPart.error != nil {
					if lastPart.error == io.EOF {
						return nil
					}
					return lastPart.error
				}
			case lastPart = <-d.queue:
				if lastPart.error != nil {
					if lastPart.error == io.EOF {
						return nil
					}
					return lastPart.error
				}
			}
		} else {
			return nil
		}
	}
}

func (d *DownloadParts) addPartsToQueue() {
	for _, part := range d.parts {
		d.queue <- part
	}
}

func (d *DownloadParts) listenOnQueue() {
	go func() {
		for {
			select {
			case part := <-d.queue:
				if d.queueContext.Err() != nil {
					d.finishedParts <- part
					continue
				}
				if part.processing {
					panic(part)
				}
				if len(part.requests) > 3 {
					d.Config.LogPath(d.FileInfo.Name(), map[string]interface{}{"message": "Maxed out reties", "part": part.number})
					d.finishedParts <- part
				} else {
					part.Clear()
					d.waitOnPerPartLimit()
					if d.runningParts() != 0 && d.runningParts()%5 == 0 {
						time.Sleep(time.Duration(100*d.runningParts()) * time.Millisecond)
					}
					if d.anyTimeouts() {
						d.runningPartLimit = int(math.Max(float64(1), math.Ceil(float64(d.runningPartLimit/2))))
						time.Sleep(time.Duration(100*d.runningParts()) * time.Millisecond)
						d.stateLog()
					}
					d.waitOnPerPartLimit()
					if manager.Wait(d.Context, d.globalWait) {
						go d.processPart(part.Start(d.Context), false)
						if d.WaitOnFirstSuccess() == RunError {
							break
						}
					} else {
						break
					}
				}
			}
		}
	}()
}

func (d *DownloadParts) waitOnPerPartLimit() {
	t := time.Now()
	for {
		if d.runningParts() < d.runningPartLimit {
			break
		}
		if time.Now().After(t.Add(time.Second * 5)) {
			t = time.Now()
			d.stateLog()
		}
	}
}

func (d *DownloadParts) stateLog() {
	d.Config.LogPath(
		d.FileInfo.Name(),
		d.state(),
	)
}

func (d *DownloadParts) state() map[string]interface{} {
	return map[string]interface{}{
		"RunningParts": d.runningParts(),
		"limit":        d.runningPartLimit,
		"parts":        len(d.parts),
		"written":      d.totalWritten,
	}
}

func (d *DownloadParts) SizeGuess() bool {
	if untrusted, ok := d.FileInfo.(UntrustedSize); ok {
		return untrusted.UntrustedSize()
	}
	return false
}

func (d *DownloadParts) buildParts() {
	size := d.FileInfo.Size()
	if d.SizeGuess() {
		d.Config.LogPath(d.FileInfo.Name(), map[string]interface{}{
			"message": "Using size guess",
			"size":    size,
		})
	}
	for i, offset := range byteChunkSlice(size, DownloadPartChunkSize) {
		d.parts = append(d.parts, &Part{OffSet: offset, number: int64(i) + 1})
	}

	d.finishedParts = make(chan *Part, len(d.parts))
	d.queue = make(chan *Part, len(d.parts))
}

func (d *DownloadParts) runningParts() int {
	var running int
	for _, part := range d.parts {
		if part.processing {
			running += 1
		}
	}

	return running
}

func (d *DownloadParts) anyTimeouts() bool {
	for _, part := range d.parts {
		if part.error == nil {
			continue
		}
		if part.error == context.DeadlineExceeded {
			return true
		}
	}

	return false
}

func (d *DownloadParts) processPart(part *Part, UnexpectedEOF bool) {
	ranger, ok := d.File.(ReaderRange)
	if ok {
		d.processRanger(part, ranger, UnexpectedEOF)
	} else {
		d.processDefault(part, UnexpectedEOF)
	}
	d.globalWait.Done()
}

func (d *DownloadParts) partRequestTimeoutValue() time.Duration {
	if d.runningPartLimit <= 6 {
		return time.Second * 60
	} else {
		// Get stricter with more running part in order to not lock up on a slow server with many connections.
		return time.Second * time.Duration(
			math.Max(15, 60-float64(d.runningParts())),
		)
	}
}

func (d *DownloadParts) processRanger(part *Part, ranger ReaderRange, UnexpectedEOF bool) {
	withContext, ok := ranger.(WithContext)
	if ok {
		timeoutValue := d.partRequestTimeoutValue()

		d.Config.LogPath(
			d.FileInfo.Name(),
			lo.Assign(d.state(), map[string]interface{}{
				"message":      "starting readerRange",
				"TimeOutValue": timeoutValue,
			}),
		)
		partCtx, partCancel := context.WithTimeout(part.Context, timeoutValue)
		defer partCancel()
		ranger = withContext.WithContext(partCtx).(ReaderRange)
	}
	r, err := ranger.ReaderRange(part.off, part.len+part.off-1)
	if d.readInitCheck(part, err, "readRange") {
		return
	}
	defer r.Close()

	wn, err := lib.CopyAt(d.WriterAndAt, part.off, r)
	part.bytes = wn
	atomic.AddInt64(&d.totalWritten, wn)
	if d.requeueOnError(part, err, "copyAt", UnexpectedEOF) {
		return
	}
	d.verifySizeAndUpdateParts(part, wn)

	d.finishedParts <- part.Done()
}

func (d *DownloadParts) verifySizeAndUpdateParts(part *Part, wn int64) {
	if wn < part.len {
		d.Config.LogPath(
			d.FileInfo.Name(),
			map[string]interface{}{"message": "returned less than expected", "part": part.number},
		)
		d.queueCancel()
		// cancelAll greater parts
		for _, p := range d.parts[part.number:] {
			d.Config.LogPath(
				d.FileInfo.Name(),
				map[string]interface{}{"message": "canceling invalid part", "part": part.number},
			)
			p.CancelFunc()
		}
	}
}

func (d *DownloadParts) readInitCheck(part *Part, err error, op string) bool {
	if err != nil {
		part.error = &fs.PathError{
			Path: d.FileInfo.Name(),
			Err:  err,
			Op:   op,
		}
		if errors.Is(err, context.DeadlineExceeded) {
			d.queue <- part.Done()
		} else {
			d.finishedParts <- part.Done() // the request is already retried 3 times internally
		}
		return true
	}
	return false
}

func (d *DownloadParts) requeueOnError(part *Part, err error, op string, UnexpectedEOF bool) bool {
	if err != nil && !errors.Is(err, io.EOF) {
		if UnexpectedEOF && errors.Is(err, io.ErrUnexpectedEOF) {
			return false
		}
		part.error = &fs.PathError{
			Path: d.FileInfo.Name(),
			Err:  err,
			Op:   op,
		}
		d.Config.LogPath(
			d.FileInfo.Name(),
			map[string]interface{}{"message": "requeuing", "error": part.error, "part": part.number},
		)
		d.queue <- part.Done() // either timeout or stream error try part again.
		return true
	}

	return false
}

// This is just a stub to use the default fs.File. It works the same just won't provide realtime progress.
func (d *DownloadParts) processDefault(part *Part, UnexpectedEOF bool) {
	buf := make([]byte, part.len)
	readAtCloser, ok := d.File.(lib.ReaderAtCloser)
	if !ok {
		panic(fmt.Errorf("can't convert fs.File to lib.ReaderAtCloser"))
	}
	timeoutValue := d.partRequestTimeoutValue()
	d.Config.LogPath(
		d.FileInfo.Name(),
		map[string]interface{}{"message": "readAt request", "TimeOutValue": timeoutValue, "part": part.number},
	)
	partCtx, partCancel := context.WithTimeout(d.Context, timeoutValue)
	defer partCancel()
	rn, err := lib.NewReaderAt(partCtx, readAtCloser).ReadAt(buf, part.off)
	part.bytes += int64(rn)
	if d.readInitCheck(part, err, "readerAt") {
		return
	}
	defer readAtCloser.Close()
	wn, err := d.WriterAndAt.WriteAt(buf, part.off)

	atomic.AddInt64(&d.totalWritten, int64(wn))
	if d.requeueOnError(part, err, "copyAt", UnexpectedEOF) {
		return
	}
	d.verifySizeAndUpdateParts(part, int64(wn))

	d.finishedParts <- part.Done()
}

func (d *DownloadParts) downloadFile() error {
	withContext, ok := d.File.(WithContext)
	if ok {
		d.File = withContext.WithContext(d.Context).(fs.File)
	}
	n, err := io.Copy(d.WriterAndAt, d.File)
	atomic.AddInt64(&d.totalWritten, n)
	if err != nil {
		return err
	}
	err = d.File.Close()
	if err != nil {
		return err
	}
	if !d.SizeGuess() {
		return nil
	}

	realTime, ok := d.File.(RealTimeStat)
	if ok {
		d.FileInfo, err = realTime.RealTimeStat()
		if err != nil {
			return err
		}
	} else {
		d.FileInfo, err = d.File.Stat()
		if err != nil {
			return err
		}
	}
	if !d.SizeGuess() {
		if d.FileInfo.Size() != d.totalWritten {
			return fmt.Errorf("server reported size does not match downloaded file. - expected: %v, actual: %v", d.FileInfo.Size(), d.totalWritten)
		}
	}
	return nil
}

func (d *DownloadParts) WaitOnFirstSuccess() int {
	if d.firstPartSuccess > 0 {
		return d.firstPartSuccess
	}

	d.firstPartSuccess = <-d.firstPartSuccessChan
	return d.firstPartSuccess
}

func (d *DownloadParts) TriggerFirstSuccess(status int) {
	if d.firstPartSuccess > 0 {
		return
	}

	d.Config.LogPath(
		d.FileInfo.Name(),
		map[string]interface{}{"message": "TriggerFirstSuccess", "status": status},
	)
	d.firstPartSuccessChan <- status
}
