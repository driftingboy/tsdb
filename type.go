package tsdb

import (
	"sync"

	"github.com/cespare/xxhash"
)

// DataPoint: 随时间收集的数据点 <收集时间，当前数据>
type DataPoint struct {
	Ts int64
	V  float64
}

// Metric: 收集数据的指标，指标名称+标签集，如： name{label_key0:label_value0, label_key1:label_value1}
type Metric struct {
	Name     string
	LabelSet []Label
}

var labelBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 1024)
	},
}

type Label struct {
	Name  string
	Value string
}

// LabelSet 表示 Label 组合
type LabelSet []Label

func (ls LabelSet) Hash() uint64 {
	b := labelBufPool.Get().([]byte)
	defer func() {
		// reset buf
		b = b[:0]
		labelBufPool.Put(b)
	}()

	for _, l := range ls {
		b = append(b, l.Name...)
		b = append(b, l.Value...)
	}

	return xxhash.Sum64(b)
}

// 一个数据样本, 一个数据点所包含的所有数据
// (Metric+LabelSet => y, Point.Timestamp => x) = Point.value
type Sample struct {
	Id string
	// Metric     string
	// Labels     LabelSet
	DataPoint
}

// func (s Sample) ID() string {
// 	// return s.Id
// 	return fmt.Sprintf("%d-%d", xxhash.Sum64([]byte(s.Metric)), s.Labels.Hash())
// }
