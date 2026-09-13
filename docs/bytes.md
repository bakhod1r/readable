---
title: Bytes & throughput
layout: default
nav_order: 4
---

# Bytes & throughput
{: .no_toc }

Sizes in binary, IEC and SI units, and data transfer rates.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## Bytes

```go
Package bytes implements functions for the manipulation of byte slices. It is
analogous to the facilities of the strings package.

const MinRead = 512
var ErrTooLarge = errors.New("bytes.Buffer: too large")
func Clone(b []byte) []byte
func Compare(a, b []byte) int
func Contains(b, subslice []byte) bool
func ContainsAny(b []byte, chars string) bool
func ContainsFunc(b []byte, f func(rune) bool) bool
func ContainsRune(b []byte, r rune) bool
func Count(s, sep []byte) int
func Cut(s, sep []byte) (before, after []byte, found bool)
func CutPrefix(s, prefix []byte) (after []byte, found bool)
func CutSuffix(s, suffix []byte) (before []byte, found bool)
func Equal(a, b []byte) bool
func EqualFold(s, t []byte) bool
func Fields(s []byte) [][]byte
func FieldsFunc(s []byte, f func(rune) bool) [][]byte
func FieldsFuncSeq(s []byte, f func(rune) bool) iter.Seq[[]byte]
func FieldsSeq(s []byte) iter.Seq[[]byte]
func HasPrefix(s, prefix []byte) bool
func HasSuffix(s, suffix []byte) bool
func Index(s, sep []byte) int
func IndexAny(s []byte, chars string) int
func IndexByte(b []byte, c byte) int
func IndexFunc(s []byte, f func(r rune) bool) int
func IndexRune(s []byte, r rune) int
func Join(s [][]byte, sep []byte) []byte
func LastIndex(s, sep []byte) int
func LastIndexAny(s []byte, chars string) int
func LastIndexByte(s []byte, c byte) int
func LastIndexFunc(s []byte, f func(r rune) bool) int
func Lines(s []byte) iter.Seq[[]byte]
func Map(mapping func(r rune) rune, s []byte) []byte
func Repeat(b []byte, count int) []byte
func Replace(s, old, new []byte, n int) []byte
func ReplaceAll(s, old, new []byte) []byte
func Runes(s []byte) []rune
func Split(s, sep []byte) [][]byte
func SplitAfter(s, sep []byte) [][]byte
func SplitAfterN(s, sep []byte, n int) [][]byte
func SplitAfterSeq(s, sep []byte) iter.Seq[[]byte]
func SplitN(s, sep []byte, n int) [][]byte
func SplitSeq(s, sep []byte) iter.Seq[[]byte]
func Title(s []byte) []byte
func ToLower(s []byte) []byte
func ToLowerSpecial(c unicode.SpecialCase, s []byte) []byte
func ToTitle(s []byte) []byte
func ToTitleSpecial(c unicode.SpecialCase, s []byte) []byte
func ToUpper(s []byte) []byte
func ToUpperSpecial(c unicode.SpecialCase, s []byte) []byte
func ToValidUTF8(s, replacement []byte) []byte
func Trim(s []byte, cutset string) []byte
func TrimFunc(s []byte, f func(r rune) bool) []byte
func TrimLeft(s []byte, cutset string) []byte
func TrimLeftFunc(s []byte, f func(r rune) bool) []byte
func TrimPrefix(s, prefix []byte) []byte
func TrimRight(s []byte, cutset string) []byte
func TrimRightFunc(s []byte, f func(r rune) bool) []byte
func TrimSpace(s []byte) []byte
func TrimSuffix(s, suffix []byte) []byte
type Buffer struct{ ... }
```

func NewBuffer(buf []byte) *Buffer
func NewBufferString(s string) *Buffer
type Reader struct{ ... }
func NewReader(b []byte) *Reader

## BytesIEC

```go
func BytesIEC(n uint64) string
```

BytesIEC formats n bytes with 1024-based IEC units: B, KiB, MiB, GiB, TiB,
PiB, EiB.

```go
BytesIEC(1024) // "1 KiB"
```

## BytesSI

```go
func BytesSI(n uint64) string
```

BytesSI formats n bytes with 1000-based SI units: B, kB, MB, GB, TB, PB, EB.

```go
BytesSI(1000) // "1 kB"
```

## FileSize

```go
func FileSize(n int64) string
```

FileSize formats a file size (as returned by os.FileInfo.Size) for display
to end users: binary units like Bytes, but at most one fraction digit.
Negative sizes keep their sign.

```go
FileSize(5242880) // "5 MB"
FileSize(1590000) // "1.5 MB"
```

## BytesRate

```go
func BytesRate(n uint64, period time.Duration) string
```

BytesRate formats n bytes transferred over period as binary units per second
like Bytes followed by "/s". The per-second value is computed exactly in
integer arithmetic, truncated, and saturates at MaxUint64 bytes.

```go
BytesRate(5662310, time.Second) // "5.4 MB/s"
```

A non-positive period yields "NaN/s".

## Throughput

```go
func Throughput(bytesPerSecond uint64) string
```

Throughput formats a data rate in bytes per second like Bytes followed by
"/s".

```go
Throughput(125000000) // "119.21 MB/s"
```

## ThroughputIEC

```go
func ThroughputIEC(bytesPerSecond uint64) string
```

ThroughputIEC is like Throughput with IEC units (KiB, MiB, ...).

```go
ThroughputIEC(125000000) // "119.21 MiB/s"
```

## Bandwidth

```go
func Bandwidth(bitsPerSecond uint64) string
```

Bandwidth formats a network rate in bits per second with 1000-based units
bps, Kbps, Mbps, Gbps and Tbps and up to DefaultPrecision trimmed fraction
digits, rounded half up. Rates of a quadrillion bps or more stay in Tbps.

```go
Bandwidth(125000000) // "125 Mbps"
```
