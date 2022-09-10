package diskqueue

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiskQueue_Len(t *testing.T) {
	q, err := newDiskQueue()
	assert.Nil(t, err)
	defer q.Close()

	assert.Equal(t, int64(0), q.Len())
}

func TestDiskQueue_Push(t *testing.T) {
	q, err := newDiskQueue()
	assert.Nil(t, err)
	defer q.Close()

	lines := strings.Fields("hello world this is diskqueue")
	for _, line := range lines {
		assert.Nil(t, q.Push([]byte(line)))
	}
}

func TestDiskQueue_Pop(t *testing.T) {
	q, err := newDiskQueue()
	assert.Nil(t, err)
	defer q.Close()

	assert.Nil(t, q.Clear())

	lines := strings.Fields("hello world this is diskqueue")
	for _, line := range lines {
		assert.Nil(t, q.Push([]byte(line)))
	}

	for i := 0; ; i++ {
		select {
		case line := <-q.Pop():
			assert.Equal(t, lines[i], string(line))
		case <-time.After(time.Second * 1):
			return
		}
	}
}

func TestDiskQueue_Clear(t *testing.T) {
	q, err := newDiskQueue()
	assert.Nil(t, err)
	defer q.Close()

	assert.Nil(t, q.Clear())
}

func TestDiskQueue_Close(t *testing.T) {
	q, err := newDiskQueue()
	assert.Nil(t, err)

	assert.Nil(t, q.Close())
}

func TestDiskQueue_DeleteAndExit(t *testing.T) {
	q, err := newDiskQueue()
	assert.Nil(t, err)

	assert.Nil(t, q.DeleteAndExit())
}

func newDiskQueue() (DiskQueue, error) {
	opts := []Option{
		WithName("test"),
		WithDataPath("/tmp/diskqueue"),
		WithMaxBytesPerFile(20),
		WithMinMsgSize(0),
		WithMaxMsgSize(10),
		WithSyncEvery(1),
		WithSyncTimeout(time.Second * 1),
	}
	return NewDiskQueue(opts...)
}
