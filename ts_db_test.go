package tsdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mini use case
func Test_UseCase(t *testing.T) {
	// insert
	sconf := &segmentConf{capacity: 5}
	db := OpenTSDB(sconf)
	testData := []*Sample{
		{Id: "req-total<node=127.0.0.1>", DataPoint: DataPoint{Ts: 1632557712, V: 10}},
		{Id: "req-total<node=127.0.0.1>", DataPoint: DataPoint{Ts: 1632557713, V: 2}},
		{Id: "req-total<node=127.0.0.1>", DataPoint: DataPoint{Ts: 1632557714, V: 3}},
		{Id: "req-total<node=127.0.0.1>", DataPoint: DataPoint{Ts: 1632557715, V: 4}},
		{Id: "req-total<node=127.0.0.1>", DataPoint: DataPoint{Ts: 1632557716, V: 5}},
		{Id: "req-total<node=127.0.0.1>", DataPoint: DataPoint{Ts: 1632557717, V: 11}},
	}
	err := db.InsertRows(testData)
	assert.NoError(t, err)
	assert.Equal(t, db.sl.getHead().size(), int64(6))

	// test new segment
	testData1 := []*Sample{
		{Id: "req-total<node=192.168.0.32>", DataPoint: DataPoint{Ts: 1632557712, V: 10}},
		{Id: "req-total<node=192.168.0.32>", DataPoint: DataPoint{Ts: 1632557713, V: 2}},
	}
	err = db.InsertRows(testData1)
	assert.NoError(t, err)
	assert.Equal(t, db.sl.getHead().size(), int64(2))
	// query
	// iter := db.sl.GetIterator()
	// if iter.HasNext() {
	// 	value := iter.Value()
	// 	t.Log(value.Size(), value.MaxT())
	// 	iter.Next()
	// }
}

// test concurrency
