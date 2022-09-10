package diskqueue

import (
	"time"
)

type Option func(*diskQueue)

func WithName(name string) Option {
	return func(d *diskQueue) {
		if name != "" {
			d.name = name
		}
	}
}

func WithDataPath(dataPath string) Option {
	return func(d *diskQueue) {
		if dataPath != "" {
			d.dataPath = dataPath
		}
	}
}

func WithMaxBytesPerFile(max int64) Option {
	return func(d *diskQueue) {
		if max > 0 {
			d.maxBytesPerFile = max
		}
	}
}

func WithMinMsgSize(size int32) Option {
	return func(d *diskQueue) {
		if size >= 0 {
			d.minMsgSize = size
		}
	}
}

func WithMaxMsgSize(size int32) Option {
	return func(d *diskQueue) {
		if size > 0 {
			d.maxMsgSize = size
		}
	}
}

func WithSyncEvery(count int64) Option {
	return func(d *diskQueue) {
		if count > 0 {
			d.syncEvery = count
		}
	}
}

func WithSyncTimeout(timeout time.Duration) Option {
	return func(d *diskQueue) {
		if timeout > 0 {
			d.syncTimeout = timeout
		}
	}
}
