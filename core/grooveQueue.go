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
		len:        0,
		buffer:     make([]interface{}, capi),
	}
}

func (q *GrooveQueue) EnQueue(data interface{}) bool {
	for {
		if q.cap == q.Len() {
			return false
		}
		tail := atomic.LoadUint32(&q.tail)
		nextTail := (tail + 1) & q.mask

		if atomic.CompareAndSwapUint32(&q.tail, tail, nextTail) {
			atomic.AddUint32(&q.writeState, 1)
			q.buffer[tail] = data
			atomic.AddUint32(&q.len, 1)
			atomic.AddUint32(&q.writeState, ^uint32(0))
			return true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *GrooveQueue) EnQueueBatch(data []interface{}) bool {
	n := uint32(len(data))
	for {
		if q.cap-q.Len() < n {
			return false
		}
		tail := atomic.LoadUint32(&q.tail)
		nextTail := (tail + n) & q.mask

		if atomic.CompareAndSwapUint32(&q.tail, tail, nextTail) {
			atomic.AddUint32(&q.writeState, n)
			for i := uint32(0); i < n; i++ {
				q.buffer[(tail+i)&q.mask] = data
			}
			atomic.AddUint32(&q.len, n)
			atomic.AddUint32(&q.writeState, ^(n - 1))
			return true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *GrooveQueue) DeQueue() (interface{}, bool) {
	for {
		if q.Len() == 0 {
			return nil, false
		}
		head := atomic.LoadUint32(&q.head)
		nextHead := (head + 1) & q.mask
		tail := atomic.LoadUint32(&q.tail)
		state := atomic.LoadUint32(&q.writeState)
		tailBox := (nextHead + state) & q.mask
		if nextHead <= tail {
			if tailBox > tail || tailBox < nextHead {
				runtime.Gosched()
				continue
			}
		} else {
			if tailBox > tail && tailBox < nextHead {
				runtime.Gosched()
				continue
			}
		}

		value := q.buffer[head]
		if atomic.CompareAndSwapUint32(&q.head, head, nextHead) {
			atomic.AddUint32(&q.len, ^uint32(0))
			return value, true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *GrooveQueue) DeQueueBatch(n uint32) ([]interface{}, bool) {
	for {
		if q.Len() < n {
			return nil, false
		}
		head := atomic.LoadUint32(&q.head)
		nextHead := (head + n) & q.mask
		tail := atomic.LoadUint32(&q.tail)
		state := atomic.LoadUint32(&q.writeState)
		tailBox := (nextHead + state) & q.mask
		if nextHead <= tail {
			if tailBox > tail || tailBox < nextHead {
				runtime.Gosched()
				continue
			}
		} else if nextHead > tail {
			if tailBox > tail && tailBox < nextHead {
				runtime.Gosched()
				continue
			}
		}

		value := make([]interface{}, 0, n)
		for i := uint32(0); i < n; i++ {
			value = append(value, q.buffer[(head+i)&q.mask])
		}
		if atomic.CompareAndSwapUint32(&q.head, head, nextHead) {
			atomic.AddUint32(&q.len, ^(n - 1))
			return value, true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *GrooveQueue) Len() uint32 {
	return atomic.LoadUint32(&q.len)
}
