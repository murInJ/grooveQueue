package core

import (
	"github.com/murinj/grooveQueue/utils"
	"runtime"
	"sync/atomic"
)

type GrooveQueue struct {
	padding0   [56]byte //64-(8+4)
	mask       uint32
	padding1   [56]byte
	len        uint32
	padding2   [56]byte
	cap        uint32
	padding3   [56]byte
	head       uint32
	padding4   [56]byte
	tail       uint32
	padding5   [56]byte
	writeState uint32
	padding6   [56]byte
	buffer     []interface{}
}

func NewGrooveQueue(queueCap uint32) *GrooveQueue {
	capi := utils.NextPowerOf2(queueCap)
	return &GrooveQueue{
		mask:       capi - 1,
		cap:        capi,
		head:       0,
		tail:       0,
		writeState: 0,
		buffer:     make([]interface{}, capi),
	}
}

func (q *GrooveQueue) EnQueue(data interface{}) bool {
	for {
		tail := atomic.LoadUint32(&q.tail)
		nextTail := (tail + 1) & q.mask
		if nextTail == atomic.LoadUint32(&q.head) {
			return false
		}
		if atomic.CompareAndSwapUint32(&q.tail, tail, nextTail) {
			q.buffer[tail] = data
			atomic.StoreUint32(&q.writeState, nextTail)
			return true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *GrooveQueue) DeQueue() (interface{}, bool) {
	for {
		head := atomic.LoadUint32(&q.head)
		nextHead := (head + 1) & q.mask
		if head == atomic.LoadUint32(&q.tail) {
			return nil, false
		}
		value := q.buffer[head]
		if atomic.LoadUint32(&q.writeState) != head && atomic.CompareAndSwapUint32(&q.head, head, nextHead) {
			return value, true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *GrooveQueue) Len() uint32 {
	head := atomic.LoadUint32(&q.head)
	tail := atomic.LoadUint32(&q.tail)
	if head <= tail {
		return tail - head
	} else {
		return head - tail
	}
}
