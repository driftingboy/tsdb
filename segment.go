package tsdb

type Segment interface {
	// Write operations
	//
	// insertRows is a goroutine safe way to insert data points into itself.
	// If data points older than its min timestamp were given, they won't be
	// ingested, instead, gave back as a first returned value.
	insertRows(rows []*Sample) error
	// clean removes everything managed by this partition.
	clean() error

	// Read operations
	//
	// selectDataPoints gives back certain metric's data points within the given range.
	// seriesId generate by metric string and labels
	selectDataPoints(seriesId string, start, end int64) ([]*DataPoint, error)
	// minTimestamp returns the minimum Unix timestamp in milliseconds.
	minTimestamp() int64
	// maxTimestamp returns the maximum Unix timestamp in milliseconds.
	maxTimestamp() int64
	// size returns the number of data points the partition holds.
	size() int64
	// active means not only writable but having the qualities to be the head partition.
	active() bool
	// expired means it should get removed.
	expired() bool
}
