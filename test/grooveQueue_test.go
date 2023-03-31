package test

import (
	"github.com/murinj/grooveQueue/core"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkQueue_ser(b *testing.B) {
	q := core.NewGrooveQueue(1048576)
	store := make([]int, 0, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := rand.Int() % 2
		if op == 0 {
			data := rand.Int()
			b.StartTimer()
			ok := q.EnQueue(data)
			b.StopTimer()
			if !ok {
				//b.Log("en failed")
			} else {
				store = append(store, data)
			}
		} else {
			b.StartTimer()
			data, ok := q.DeQueue()
			b.StopTimer()
			if !ok {
				//b.Log("de failed")
			} else {
				if data != store[0] {
					b.Fatal("not equal")
				}
				store = store[1:]
			}
		}
	}
}

func BenchmarkQueue_par(b *testing.B) {
	q := core.NewGrooveQueue(1048576)
	store := sync.Map{}
	cnt := int64(0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			op := rand.Int() % 2
			if op == 0 {
				data := atomic.LoadInt64(&cnt)
				atomic.AddInt64(&data, 1)
				b.StartTimer()
				ok := q.EnQueue(data)
				b.StopTimer()
				if !ok {
					//b.Log("en failed")
				} else {
					store.Store(data, data)
				}
			} else {
				b.StartTimer()
				data, ok := q.DeQueue()
				b.StopTimer()
				if !ok {
					//b.Log("de failed")
				} else {
					obj, ok1 := store.LoadAndDelete(data)
					if obj != nil {
						if !ok1 {
							b.Fatal("not equal2:", data, obj)
						}
						if data != obj {
							b.Fatal("not equal:", data, obj)
						}
					}
				}
			}
		}
	})
}
