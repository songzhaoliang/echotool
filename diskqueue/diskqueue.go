package diskqueue

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	etl "github.com/songzhaoliang/echotool"
)

var (
	ErrExited = errors.New("diskqueue has been exited")
)

type DiskQueue interface {
	Len() int64
	Push([]byte) error
	Pop() chan []byte
	Clear() error
	Close() error
	DeleteAndExit() error
}

type diskQueue struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms

	// run-time state (also persisted to disk)
	readPos      int64
	writePos     int64
	readFileNum  int64
	writeFileNum int64
	length       int64

	sync.RWMutex

	// instantiation time metadata
	name            string
	dataPath        string
	maxBytesPerFile int64 // currently this cannot change once created
	minMsgSize      int32
	maxMsgSize      int32
	syncEvery       int64         // number of writes per fsync
	syncTimeout     time.Duration // duration of time per fsync
	needSync        bool

	// keeps track of the position where we have read
	// (but not yet sent over readChan)
	nextReadPos     int64
	nextReadFileNum int64

	readFile  *os.File
	writeFile *os.File
	reader    *bufio.Reader
	writeBuf  bytes.Buffer

	readChan chan []byte

	writeChan         chan []byte
	writeResponseChan chan error

	clearChan         chan struct{}
	clearResponseChan chan error

	exitFlag     int32
	exitChan     chan struct{}
	exitSyncChan chan struct{}
}

var _ DiskQueue = (*diskQueue)(nil)

func NewDiskQueue(opts ...Option) (DiskQueue, error) {
	d := &diskQueue{
		name:              "default",
		dataPath:          os.TempDir(),
		maxBytesPerFile:   math.MaxInt64,
		minMsgSize:        0,
		maxMsgSize:        math.MaxInt32,
		readChan:          make(chan []byte),
		writeChan:         make(chan []byte),
		writeResponseChan: make(chan error),
		clearChan:         make(chan struct{}),
		clearResponseChan: make(chan error),
		exitChan:          make(chan struct{}),
		exitSyncChan:      make(chan struct{}),
		syncEvery:         20,
		syncTimeout:       time.Second * 1,
	}

	for _, opt := range opts {
		opt(d)
	}

	if err := os.MkdirAll(d.dataPath, os.ModePerm); err != nil {
		return nil, err
	}

	// no need to lock here, nothing else could possibly be touching this instance
	err := d.retrieveMetaData()
	if err != nil && !os.IsNotExist(err) {
	}

	go d.ioLoop()
	return d, nil
}

// ioLoop provides the backend for exposing a go channel (via ReadChan())
// in support of multiple concurrent queue consumers
//
// it works by looping and branching based on whether or not the queue has data
// to read and blocking until data is either read or written over the appropriate
// go channels
//
// conveniently this also means that we're asynchronously reading from the filesystem
func (d *diskQueue) ioLoop() {
	var dataRead []byte
	var err error
	var count int64
	var r chan []byte

	ticker := time.NewTicker(d.syncTimeout)

	for {
		if count == d.syncEvery {
			d.needSync = true
		}

		if d.needSync {
			if err = d.sync(); err != nil {
				etl.Error("diskqueue %s sync error - %v", d.name, err)
			}
			count = 0
		}

		if (d.readFileNum < d.writeFileNum) || (d.readPos < d.writePos) {
			if d.nextReadPos == d.readPos {
				if dataRead, err = d.readOne(); err != nil {
					etl.Error("diskqueue %s read error - %v", d.name, err)
					d.handleReadError()
					continue
				}
			}
			r = d.readChan
		} else {
			r = nil
		}

		select {
		// the Go channel spec dictates that nil channel operations (read or write)
		// in a select are skipped, we set r to d.readChan only when there is data to read
		case r <- dataRead:
			count++
			// moveForward sets needSync flag if a file is removed
			d.moveForward()
		case <-d.clearChan:
			d.clearResponseChan <- d.deleteAllFiles()
			count = 0
		case dataWrite := <-d.writeChan:
			count++
			d.writeResponseChan <- d.writeOne(dataWrite)
		case <-ticker.C:
			if count > 0 {
				d.needSync = true
			}
		case <-d.exitChan:
			goto exit
		}
	}

exit:
	etl.Info("diskqueue %s close", d.name)
	ticker.Stop()
	d.exitSyncChan <- struct{}{}
}

func (d *diskQueue) Len() int64 {
	return atomic.LoadInt64(&d.length)
}

func (d *diskQueue) Push(data []byte) error {
	d.RLock()
	defer d.RUnlock()

	if d.exitFlag == 1 {
		return ErrExited
	}

	d.writeChan <- data
	return <-d.writeResponseChan
}

func (d *diskQueue) Pop() chan []byte {
	return d.readChan
}

func (d *diskQueue) Clear() error {
	d.RLock()
	defer d.RUnlock()

	if d.exitFlag == 1 {
		return ErrExited
	}

	etl.Info("diskqueue %s clear", d.name)

	d.clearChan <- struct{}{}
	return <-d.clearResponseChan
}

func (d *diskQueue) Close() error {
	if err := d.exit(false); err != nil {
		return err
	}
	return d.sync()
}

func (d *diskQueue) DeleteAndExit() error {
	return d.exit(true)
}

func (d *diskQueue) exit(needDel bool) error {
	d.Lock()
	defer d.Unlock()

	if d.exitFlag == 1 {
		return ErrExited
	}

	d.exitFlag = 1

	if needDel {
		etl.Info("diskqueue %s delete", d.name)
	} else {
		etl.Info("diskqueue %s close", d.name)
	}

	close(d.exitChan)
	<-d.exitSyncChan

	if d.readFile != nil {
		d.readFile.Close()
		d.readFile = nil
	}

	if d.writeFile != nil {
		d.writeFile.Close()
		d.writeFile = nil
	}

	return nil
}

func (d *diskQueue) sync() error {
	if d.writeFile != nil {
		if err := d.writeFile.Sync(); err != nil {
			d.writeFile.Close()
			d.writeFile = nil
			return err
		}
	}

	if err := d.persistMetaData(); err != nil {
		return err
	}

	d.needSync = false
	return nil
}

func (d *diskQueue) retrieveMetaData() error {
	var f *os.File
	var err error

	fileName := d.metaFileName()
	f, err = os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	var length int64
	_, err = fmt.Fscanf(f, "%d\n%d,%d\n%d,%d\n",
		&length,
		&d.readFileNum, &d.readPos,
		&d.writeFileNum, &d.writePos)
	if err != nil {
		return err
	}
	atomic.StoreInt64(&d.length, length)
	d.nextReadFileNum = d.readFileNum
	d.nextReadPos = d.readPos

	return nil
}

// persistMetaData atomically writes state to the filesystem
func (d *diskQueue) persistMetaData() error {
	var f *os.File
	var err error

	fileName := d.metaFileName()
	tmpFileName := fmt.Sprintf("%s.%d.tmp", fileName, rand.Int())

	f, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "%d\n%d,%d\n%d,%d\n",
		atomic.LoadInt64(&d.length),
		d.readFileNum, d.readPos,
		d.writeFileNum, d.writePos)
	if err != nil {
		f.Close()
		return err
	}
	f.Sync()
	f.Close()

	return os.Rename(tmpFileName, fileName)
}

func (d *diskQueue) metaFileName() string {
	return fmt.Sprintf(path.Join(d.dataPath, "%s.diskqueue.meta.dat"), d.name)
}

func (d *diskQueue) fileName(fileNum int64) string {
	return fmt.Sprintf(path.Join(d.dataPath, "%s.diskqueue.%06d.dat"), d.name, fileNum)
}

// readOne performs a low level filesystem read for a single []byte
// while advancing read positions and rolling files, if necessary
func (d *diskQueue) readOne() ([]byte, error) {
	var err error
	var msgSize int32

	if d.readFile == nil {
		curFileName := d.fileName(d.readFileNum)
		d.readFile, err = os.OpenFile(curFileName, os.O_RDONLY, 0600)
		if err != nil {
			return nil, err
		}

		if d.readPos > 0 {
			_, err = d.readFile.Seek(d.readPos, 0)
			if err != nil {
				d.readFile.Close()
				d.readFile = nil
				return nil, err
			}
		}

		d.reader = bufio.NewReader(d.readFile)
	}

	err = binary.Read(d.reader, binary.BigEndian, &msgSize)
	if err != nil {
		d.readFile.Close()
		d.readFile = nil
		return nil, err
	}

	if msgSize < d.minMsgSize || msgSize > d.maxMsgSize {
		// this file is corrupt and we have no reasonable guarantee on
		// where a new message should begin
		d.readFile.Close()
		d.readFile = nil
		return nil, fmt.Errorf("invalid message read size (%d)", msgSize)
	}

	readBuf := make([]byte, msgSize)
	_, err = io.ReadFull(d.reader, readBuf)
	if err != nil {
		d.readFile.Close()
		d.readFile = nil
		return nil, err
	}

	totalBytes := int64(4 + msgSize)

	// we only advance next* because we have not yet sent this to consumers
	// (where readFileNum, readPos will actually be advanced)
	d.nextReadPos = d.readPos + totalBytes
	d.nextReadFileNum = d.readFileNum

	// TODO: each data file should embed the maxBytesPerFile
	// as the first 8 bytes (at creation time) ensuring that
	// the value can change without affecting runtime
	if d.nextReadPos > d.maxBytesPerFile {
		if d.readFile != nil {
			d.readFile.Close()
			d.readFile = nil
		}

		d.nextReadFileNum++
		d.nextReadPos = 0
	}

	return readBuf, nil
}

func (d *diskQueue) handleReadError() {
	// jump to the next read file and rename the current (bad) file
	if d.readFileNum == d.writeFileNum {
		// if you can't properly read from the current write file it's safe to
		// assume that something is fucked and we should skip the current file too
		if d.writeFile != nil {
			d.writeFile.Close()
			d.writeFile = nil
		}
		d.writeFileNum++
		d.writePos = 0
	}

	badFn := d.fileName(d.readFileNum)
	badRenameFn := badFn + ".bad"
	if err := os.Rename(badFn, badRenameFn); err != nil {
		etl.Error("diskqueue %s rename error - %v", d.name, err)
	}

	d.readFileNum++
	d.readPos = 0
	d.nextReadFileNum = d.readFileNum
	d.nextReadPos = 0

	// significant state change, schedule a sync on the next iteration
	d.needSync = true
}

// writeOne performs a low level filesystem write for a single []byte
// while advancing write positions and rolling files, if necessary
func (d *diskQueue) writeOne(data []byte) error {
	var err error

	if d.writeFile == nil {
		curFileName := d.fileName(d.writeFileNum)
		d.writeFile, err = os.OpenFile(curFileName, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return err
		}

		if d.writePos > 0 {
			_, err = d.writeFile.Seek(d.writePos, 0)
			if err != nil {
				d.writeFile.Close()
				d.writeFile = nil
				return err
			}
		}
	}

	dataLen := int32(len(data))

	if dataLen < d.minMsgSize || dataLen > d.maxMsgSize {
		return fmt.Errorf("invalid message write size (%d) maxMsgSize=%d", dataLen, d.maxMsgSize)
	}

	d.writeBuf.Reset()
	err = binary.Write(&d.writeBuf, binary.BigEndian, dataLen)
	if err != nil {
		return err
	}

	_, err = d.writeBuf.Write(data)
	if err != nil {
		return err
	}

	// only write to the file once
	_, err = d.writeFile.Write(d.writeBuf.Bytes())
	if err != nil {
		d.writeFile.Close()
		d.writeFile = nil
		return err
	}

	totalBytes := int64(4 + dataLen)
	d.writePos += totalBytes
	atomic.AddInt64(&d.length, 1)

	if d.writePos > d.maxBytesPerFile {
		d.writeFileNum++
		d.writePos = 0

		// sync every time we start writing to a new file
		if err = d.sync(); err != nil {
			etl.Error("diskqueue %s sync error - %v", d.name, err)
		}

		if d.writeFile != nil {
			d.writeFile.Close()
			d.writeFile = nil
		}
	}

	return err
}

func (d *diskQueue) deleteAllFiles() error {
	err := d.skipToNextRWFile()

	innerErr := os.Remove(d.metaFileName())
	if innerErr != nil && !os.IsNotExist(innerErr) {
		etl.Error("diskqueue %s remove metadata file error - %v", d.name, err)
		return innerErr
	}

	return err
}

func (d *diskQueue) skipToNextRWFile() error {
	var err error

	if d.readFile != nil {
		d.readFile.Close()
		d.readFile = nil
	}

	if d.writeFile != nil {
		d.writeFile.Close()
		d.writeFile = nil
	}

	for i := d.readFileNum; i <= d.writeFileNum; i++ {
		fn := d.fileName(i)
		innerErr := os.Remove(fn)
		if innerErr != nil && !os.IsNotExist(innerErr) {
			etl.Error("diskqueue %s remove data file error - %v", d.name, err)
			err = innerErr
		}
	}

	d.writeFileNum++
	d.writePos = 0
	d.readFileNum = d.writeFileNum
	d.readPos = 0
	d.nextReadFileNum = d.writeFileNum
	d.nextReadPos = 0
	atomic.StoreInt64(&d.length, 0)

	return err
}

func (d *diskQueue) moveForward() {
	oldReadFileNum := d.readFileNum
	d.readFileNum = d.nextReadFileNum
	d.readPos = d.nextReadPos
	length := atomic.AddInt64(&d.length, -1)

	// see if we need to clean up the old file
	if oldReadFileNum != d.nextReadFileNum {
		// sync every time we start reading from a new file
		d.needSync = true

		fn := d.fileName(oldReadFileNum)
		if err := os.Remove(fn); err != nil {
			etl.Error("diskqueue %s remove file %s error - %v", d.name, fn, err)
		}
	}

	d.checkTailCorruption(length)
}

func (d *diskQueue) checkTailCorruption(length int64) {
	if d.readFileNum < d.writeFileNum || d.readPos < d.writePos {
		return
	}

	// we've reached the end of the diskqueue
	// if length isn't 0 something went wrong
	if length != 0 {
		etl.Error("diskqueue %s length is not 0 and reset 0", d.name)

		// force set length 0
		atomic.StoreInt64(&d.length, 0)
		d.needSync = true
	}

	if d.readFileNum != d.writeFileNum || d.readPos != d.writePos {
		if d.readFileNum > d.writeFileNum || d.readPos > d.writePos {
			etl.Error("diskqueue %s skip to next file and reset 0", d.name)
		}

		d.skipToNextRWFile()
		d.needSync = true
	}
}
