package tsdb

import "sync"

type SegmentList interface {
	insert(s Segment)
	delete(s Segment) error
	replace(new, old Segment) error
	getHead() Segment
	GetIterator()
}

type SegmentIterator interface {
	HasNext() bool
	Next()
	CurrentSegment() Segment
}

// 双向链表
type SegmentDLinkedList struct {
	mux sync.Mutex

	length int64
	head   *SegmentNode
	tail   *SegmentNode
}

func newSegmentDLinkedList(segmentSize int64) *SegmentDLinkedList {
	segment := newMemorySegment(segmentSize)
	return &SegmentDLinkedList{
		head: &SegmentNode{val: segment},
		tail: &SegmentNode{val: segment},
	}
}

type SegmentNode struct {
	val  Segment
	next *SegmentNode
	pre  *SegmentNode
}

func (sl *SegmentDLinkedList) insert(s Segment) {
	sn := &SegmentNode{val: s, pre: sl.head}
	sl.head.next = sn
	sl.head = sn
}

func (sl *SegmentDLinkedList) delete(s Segment) error {
	panic("not implemented") // TODO: Implement
}

func (sl *SegmentDLinkedList) replace(new Segment, old Segment) error {
	panic("not implemented") // TODO: Implement
}

func (sl *SegmentDLinkedList) getHead() Segment {
	return sl.head.val
}

func (sl *SegmentDLinkedList) GetIterator() {
	panic("not implemented") // TODO: Implement
}
