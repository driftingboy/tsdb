package tsdb

import (
	"github/driftingboy/tsdb/logger"
	"os"
)

type segmentConf struct {
	capacity int64
}

type TSDB struct {
	sl SegmentList
	// v1 同步插入（更新为异步worker插入）
	// workers
	sc     *segmentConf
	logger logger.Logger
}

func OpenTSDB(sc *segmentConf) *TSDB {
	return &TSDB{
		sl:     newSegmentDLinkedList(sc.capacity),
		sc:     sc,
		logger: logger.WithPrefix(logger.NewStdLogger(os.Stdout, 4096), logger.DefaultCaller, logger.DefaultTimer),
	}
}

func (ts *TSDB) InsertRows(rows []*Sample) error {
	if len(rows) == 0 {
		return nil
	}
	segment := ts.GetAvailableHead()
	return segment.insertRows(rows)
}

func (ts *TSDB) Query() {}

func (ts *TSDB) GetAvailableHead() (segment Segment) {
	segment = ts.sl.getHead()
	if segment.active() {
		return
	}

	segment = newMemorySegment(ts.sc.capacity)
	ts.sl.insert(segment)
	go func() {
		if err := ts.flushToDisk(); err != nil {
			ts.logger.Log(logger.Error, "err", err)
		}
	}()
	return
}

// flushPartitions persists all in-memory partitions ready to persisted.
// For the in-memory mode, just removes it from the partition list.
func (ts *TSDB) flushToDisk() error {
	return nil
}
