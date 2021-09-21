package tsdb

import (
	"os"
	"time"
)

type diskSegment struct {
	dirPath string
	meta    meta
	// file descriptor of data file
	f *os.File
	// memory-mapped file backed by f
	mappedFile []byte
	// duration to store data
	retention time.Duration
}

// meta is a mapper for a meta file, which is put for each partition.
// Note that the CreatedAt is surely timestamped by tstorage but Min/Max Timestamps are likely to do by other process.
type meta struct {
	MinTimestamp  int64                 `json:"minTimestamp"`
	MaxTimestamp  int64                 `json:"maxTimestamp"`
	NumDataPoints int                   `json:"numDataPoints"`
	Metrics       map[string]diskMetric `json:"metrics"`
	CreatedAt     time.Time             `json:"createdAt"`
}

// diskMetric holds meta data to access actual data from the memory-mapped file.
type diskMetric struct {
	Name          string `json:"name"`
	Offset        int64  `json:"offset"`
	MinTimestamp  int64  `json:"minTimestamp"`
	MaxTimestamp  int64  `json:"maxTimestamp"`
	NumDataPoints int64  `json:"numDataPoints"`
}

func (ds *diskSegment) insertRows(rows []*Sample) error {
	panic("not implemented") // TODO: Implement
}

func (ds *diskSegment) clean() error {
	panic("not implemented") // TODO: Implement
}

func (ds *diskSegment) selectDataPoints(metric string, labels []Label, start int64, end int64) ([]*DataPoint, error) {
	panic("not implemented") // TODO: Implement
}

func (ds *diskSegment) minTimestamp() int64 {
	panic("not implemented") // TODO: Implement
}

func (ds *diskSegment) maxTimestamp() int64 {
	panic("not implemented") // TODO: Implement
}

func (ds *diskSegment) size() int {
	panic("not implemented") // TODO: Implement
}

func (ds *diskSegment) active() bool {
	panic("not implemented") // TODO: Implement
}

func (ds *diskSegment) expired() bool {
	panic("not implemented") // TODO: Implement
}
