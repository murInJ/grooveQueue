package test

import (
	"github.com/murinj/grooveQueue/core"
	"math/rand"
	"testing"
)

func BenchmarkQueue_ser(b *testing.B) {
	q := core.NewGrooveQueue(1048576) //1048576
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := rand.Int() % 2
		if op == 0 {
			data := rand.Int()
			b.StartTimer()
			q.EnQueue(data)
			b.StopTimer()
		} else {
			b.StartTimer()
			q.DeQueue()
			b.StopTimer()
		}
	}
}

func BenchmarkQueue_par(b *testing.B) {
	q := core.NewGrooveQueue(1048576) //1048576
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			op := rand.Int() % 2
			if op == 0 {
				data := rand.Int()
				b.StartTimer()
				q.EnQueue(data)
				b.StopTimer()
			} else {
				b.StartTimer()
				q.DeQueue()
				b.StopTimer()
			}
		}
	})
}

func BenchmarkQueueBatch_100_ser(b *testing.B) {
	q := core.NewGrooveQueue(1048576)
	num := 100
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := rand.Int() % 2
		if op == 0 {
			data := make([]interface{}, num)
			for i := 0; i < num; i++ {
				data[i] = rand.Int()
			}
			b.StartTimer()
			q.EnQueueBatch(data)
			b.StopTimer()
		} else {
			b.StartTimer()
			q.DeQueueBatch(uint32(num))
			b.StopTimer()
		}
	}
}

func BenchmarkQueueBatch_100_par(b *testing.B) {
	q := core.NewGrooveQueue(1048576)
	num := 100
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			op := rand.Int() % 2
			if op == 0 {
				data := make([]interface{}, num)
				for i := 0; i < num; i++ {
					data[i] = rand.Int()
				}
				b.StartTimer()
				q.EnQueueBatch(data)
				b.StopTimer()
			} else {
				b.StartTimer()
				q.DeQueueBatch(uint32(num))
				b.StopTimer()
			}
		}
	})
}
