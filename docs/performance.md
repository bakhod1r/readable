---
title: Performance
layout: default
nav_order: 12
---

# Performance

Core formatters run in roughly 40–100 ns/op with a single allocation — the
returned string. Reproduce on your machine:

```sh
go test -run '^$' -bench . -benchmem
```

For hot paths use the `AppendX` variants, which write into a caller-owned
buffer and allocate nothing when it has spare capacity:

```text
BenchmarkAppendNumber     39.6 ns/op   0 B/op   0 allocs/op
BenchmarkAppendBytes      40.9 ns/op   0 B/op   0 allocs/op
BenchmarkAppendDuration   49.9 ns/op   0 B/op   0 allocs/op
BenchmarkAppendMoney      58.9 ns/op   0 B/op   0 allocs/op
BenchmarkAppendPercent   227.2 ns/op   0 B/op   0 allocs/op
```
