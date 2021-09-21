package tsdb

type TSDB struct {
	sl SegmentList

	// v1 同步插入（更新为异步worker插入）
	// workers
	logger Logger
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

	segment = newMemorySegment()
	ts.sl.insert(segment)
	go func() {
		if err := ts.flushToDisk(); err != nil {
			ts.logger.Printf("failed to flush in-memory partitions: %v", err)
		}
	}()
	return
}

// flushPartitions persists all in-memory partitions ready to persisted.
// For the in-memory mode, just removes it from the partition list.
func (ts *TSDB) flushToDisk() error {
	return nil
}
