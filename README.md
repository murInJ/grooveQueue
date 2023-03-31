# grooveQueue
go lockless queue

## now surport
- Entry and exit of individual data
- Supports multiple reads and writes
## Latest Performance Tests
### go benchmark
> goos: windows
>
>goarch: amd64
>
>pkg: github.com/murinj/grooveQueue/test
>
>cpu: AMD Ryzen 5 4600H with Radeon Graphics
>
>BenchmarkQueue_ser
>
>BenchmarkQueue_ser-12            2587262               481.7 ns/op
>
>BenchmarkQueue_par
>
>BenchmarkQueue_par-12            4043463               319.3 ns/op
>
>PASS
