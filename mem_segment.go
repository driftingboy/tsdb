package tsdb

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	//"github.com/dgryski/go-tsz"
)

// A memoryPartition implements a partition to store data points on heap.
// TSDB will eliminate the oldest memSegment when it reaches the configured memory limit.
// The memSegment's memory usage reaches the upper limit and will be flushed to the disk.
// It is concurrency safe.
type memSegment struct {
	// A hash map from metric name to memoryMetric.
	seriesSet sync.Map // Sync map 适合读远远多于写的场景, series（指标） 在运行时新增的频率非常小
	// The number of data points
	numPoints int64
	// minT is immutable.
	minT int64
	maxT int64

	// Write ahead log.
	// wal wal
	// The timestamp range of partitions after which they get persisted
	// segmentDuration int64
	// timestampPrecision TimestampPrecision
	once sync.Once
}

// TODO newMemoryPartition 提供初始化 series（指标）的 option
func newMemorySegment() Segment {
	return &memSegment{
		minT: math.MinInt64,
		maxT: math.MaxInt64,
	}
}

func (ms *memSegment) insertRows(samples []*Sample) error {
	for _, s := range samples {
		series := ms.getOrCreateSeriesByID(s.Id)
		series.insertPoint(&s.DataPoints)
	}
	return nil
}

func (ms *memSegment) clean() error {
	return nil
}

func (ms *memSegment) selectDataPoints(metricId string, start int64, end int64) ([]*DataPoint, error) {
	series := ms.getOrCreateSeriesByID(metricId)
	return series.selectPoints(start, end), nil
}

func (ms *memSegment) minTimestamp() int64 {
	return ms.minT
}

func (ms *memSegment) maxTimestamp() int64 {
	return ms.maxT
}

func (ms *memSegment) size() int64 {
	return ms.numPoints
}

func (ms *memSegment) active() bool {
	return true
}

func (ms *memSegment) expired() bool {
	return false
}

func (ms *memSegment) getOrCreateSeriesByID(id string) *memorySeries {
	m, _ := ms.seriesSet.LoadOrStore(id,
		&memorySeries{
			id:             id,
			points:         make([]*DataPoint, 0, 1000),
			disOrderPoints: make([]*DataPoint, 0),
		},
	)
	return m.(*memorySeries)
}

// 内存中的时间序列
type memorySeries struct {
	mu sync.RWMutex

	id           string
	len          int64
	minTimestamp int64
	maxTimestamp int64
	// points must kept in order
	points         []*DataPoint
	disOrderPoints []*DataPoint
	// 压缩存储
	// block *tsz.Series
}

func (m *memorySeries) insertPoint(point *DataPoint) {
	len := atomic.LoadInt64(&m.len)
	// TODO: 互斥锁的优化
	// 方案1. 分片锁,保证前面的数据查询新能不受影响（tsdb是顺序的不会修改之前的数据）
	// 方案2. copy on write:
	/*
		m.points := make([]*DataPoint, 1000)
		for i := 0; i < 1000; i++ {
			m.points[i] = point
		}
	*/
	m.mu.Lock()
	defer m.mu.Unlock()

	isFirstInsert := len == 0
	if isFirstInsert {
		m.points = append(m.points, point)
		atomic.StoreInt64(&m.minTimestamp, point.Ts)
		atomic.StoreInt64(&m.maxTimestamp, point.Ts)
		atomic.AddInt64(&m.len, 1)
		return
	}

	order := m.points[len-1].Ts < point.Ts
	if order {
		m.points = append(m.points, point)
		atomic.StoreInt64(&m.maxTimestamp, point.Ts)
		atomic.AddInt64(&m.len, 1)
		return
	}
	m.disOrderPoints = append(m.disOrderPoints, point)
}

func (m *memorySeries) selectPoints(start, end int64) []*DataPoint {
	length := atomic.LoadInt64(&m.len)
	max := atomic.LoadInt64(&m.maxTimestamp)
	min := atomic.LoadInt64(&m.minTimestamp)
	if unreachable := end < min; unreachable {
		return []*DataPoint{}
	}
	if startBeforeInitialTime := start < min; startBeforeInitialTime {
		start = min
	}
	if endAfterLastTime := end > max; endAfterLastTime {
		end = max
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	// TODO 实现 sortSet 接口（1，arry 2. avl)
	// get first index equal start
	starIndex := sort.Search(int(length), func(i int) bool {
		return m.points[i].Ts >= start
	})
	// get last index equal end
	endIndex := sort.Search(int(length), func(i int) bool {
		return m.points[i].Ts >= end+1
	}) - 1

	result := make([]*DataPoint, endIndex-starIndex+1)
	copy(result, m.points[starIndex:endIndex+1])
	return result
}
