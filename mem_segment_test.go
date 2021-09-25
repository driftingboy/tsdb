package tsdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_memSegment_insertRows(t *testing.T) {
	type result struct {
		wantErr           bool
		wantDisOrderPoint []*DataPoint
		wantDataPoint     []*DataPoint
	}

	tests := []struct {
		name    string
		segment Segment
		samples []*Sample
		result  result
	}{
		{
			name:    "test-insert-order-points",
			segment: newMemorySegment(1024),
			samples: []*Sample{
				{Id: "req_total{node=1}", DataPoint: DataPoint{Ts: 1, V: 1}},
				{Id: "req_total{node=1}", DataPoint: DataPoint{Ts: 2, V: 2}},
				{Id: "req_total{node=1}", DataPoint: DataPoint{Ts: 3, V: 1}},
			},
			result: result{wantErr: false, wantDisOrderPoint: nil, wantDataPoint: []*DataPoint{
				{Ts: 1, V: 1}, {Ts: 2, V: 2}, {Ts: 3, V: 1},
			}},
		},
		{
			name:    "test-insert-disorder-points",
			segment: newMemorySegment(1024),
			samples: []*Sample{
				{Id: "req_total{node=1}", DataPoint: DataPoint{Ts: 1, V: 1}},
				{Id: "req_total{node=1}", DataPoint: DataPoint{Ts: 3, V: 2}},
				{Id: "req_total{node=1}", DataPoint: DataPoint{Ts: 2, V: 1}},
			},
			result: result{
				wantErr:           false,
				wantDisOrderPoint: []*DataPoint{{Ts: 2, V: 1}},
				wantDataPoint:     []*DataPoint{{Ts: 1, V: 1}, {Ts: 3, V: 2}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.segment.insertRows(tt.samples)
			assert.Equal(t, tt.result.wantErr, err != nil)
			dataList, err := tt.segment.selectDataPoints("req_total{node=1}", 1, 3)
			assert.Equal(t, tt.result.wantErr, err != nil)
			assert.Equal(t, tt.result.wantDataPoint, dataList)
		})
	}
}
