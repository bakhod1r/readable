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
