---
title: Logs
layout: default
nav_order: 10
---

# Logs
{: .no_toc }

Compact, space-free, machine-stable forms for structured logs.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## LogDuration

```go
func LogDuration(d time.Duration) string
```

LogDuration formats d for machine-readable logs. It is exactly
time.Duration.String: stable, no spaces, full precision.

```go
LogDuration(1532 * time.Millisecond) // "1.532s"
LogDuration(time.Hour)               // "1h0m0s"
```

## LogBytes

```go
func LogBytes(n uint64) string
```

LogBytes formats n bytes for logs with IEC units, up to DefaultPrecision
trimmed fraction digits and no space.

```go
LogBytes(12345678) // "11.77MiB"
```

## LogNumber

```go
func LogNumber(n int64) string
```

LogNumber formats n as a plain base-10 integer with no grouping.

```go
LogNumber(1234567) // "1234567"
```
