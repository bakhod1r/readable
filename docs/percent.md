---
title: Percent, progress & rates
layout: default
nav_order: 7
---

# Percent, progress & rates
{: .no_toc }

Ratios, percentage changes, progress bars and event rates.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## Percent

```go
func Percent(ratio float64) string
```

Percent formats a ratio as a percentage with up to DefaultPrecision fraction
digits, rounded to nearest and trimmed of trailing zeros: 0.15 means 15%.
Results that round to zero lose their sign.

```go
Percent(0.1534) // "15.34%"
Percent(1)      // "100%"
```

NaN and infinities format as "NaN%", "+Inf%" and "-Inf%".

## PercentWithPrecision

```go
func PercentWithPrecision(ratio float64, precision int) string
```

PercentWithPrecision is like Percent with at most precision fraction digits
(clamped to [0, MaxPrecision]).

```go
PercentWithPrecision(0.123456, 0) // "12%"
PercentWithPrecision(0.123456, 1) // "12.3%"
```

## PercentChange

```go
func PercentChange(from, to float64) (string, error)
```

PercentChange formats the relative change from from to to, measured against
the magnitude of from, like Percent.

```go
PercentChange(100, 120) // "20%", nil
PercentChange(100, 80)  // "-20%", nil
```

It returns ErrUndefined if from is zero, any input is NaN or infinite,
or the change overflows.

## ErrUndefined

```go
var ErrUndefined = errors.New("readable: undefined result")
```

ErrUndefined is returned when a result has no meaningful value, such as a
percentage change from zero. Test for it with errors.Is.

## Progress

```go
func Progress(done, total int64) string
```

Progress formats done out of total as a whole percentage, rounded down so
"100%" appears only when done reaches total. done is clamped to [0, total];
a non-positive total yields "0%".

```go
Progress(73, 100)   // "73%"
Progress(999, 1000) // "99%"
```

## ProgressBar

```go
func ProgressBar(done, total int64, width int) string
```

ProgressBar renders a text progress bar of width cells followed by Progress.
Filled cells (U+2588) are floor(width*done/total); the rest are U+2591.
A non-positive width means 20; widths above 200 are capped.

```go
ProgressBar(73, 100, 20) // "██████████████░░░░░░ 73%"
```

## Rate

```go
func Rate(count float64, period time.Duration) string
```

Rate formats count events observed over period as a per-second rate like
NumberFloat followed by "/s".

```go
Rate(125000, time.Second) // "125K/s"
Rate(300, time.Minute)    // "5/s"
```

A non-positive period yields "NaN/s".

## RateWithLabel

```go
func RateWithLabel(count float64, period time.Duration, label string) string
```

RateWithLabel is like Rate with label and a space between the value and
"/s". An empty label is omitted.

```go
RateWithLabel(1200, time.Second, "req") // "1.2K req/s"
```

## PerMinute

```go
func PerMinute(count float64, label string) string
```

PerMinute formats count like NumberFloat followed by an optional label and
"/min".

```go
PerMinute(1200, "req") // "1.2K req/min"
PerMinute(1200, "")    // "1.2K/min"
```

## RequestRate

```go
func RequestRate(perSecond float64) string
```

RequestRate formats a request rate like NumberFloat followed by " req/s".

```go
RequestRate(15234) // "15.23K req/s"
```
