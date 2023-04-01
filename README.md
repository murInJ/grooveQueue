# grooveQueue
go lockless queue

## now surport
- Entry and exit of individual data
- - Entry and exit of data on batch
- Supports multiple reads and writes
## Latest Performance Tests
### go benchmark
goos: windows

goarch: amd64

pkg: github.com/murinj/grooveQueue/test

cpu: AMD Ryzen 5 4600H with Radeon Graphics

BenchmarkQueue_ser

BenchmarkQueue_ser-12            2170088               520.5 ns/op

BenchmarkQueue_par

BenchmarkQueue_par-12            1653750               650.8 ns/op

BenchmarkQueueBatch_ser

BenchmarkQueueBatch_ser-12        379058              3175 ns/op

BenchmarkQueueBatch_par

BenchmarkQueueBatch_par-12        379910              3169 ns/op

PASS
