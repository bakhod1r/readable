---
title: Durations & time
layout: default
nav_order: 5
---

# Durations & time
{: .no_toc }

Durations at different lengths, latency, relative times, dates and ranges.

`RelativeTime`, `Date`, `Time` and `TimeRange` read `time.Now`. Each has a `...From(now, ...)` twin that takes an explicit `now` — use it in tests.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## Duration

```go
func Duration(d time.Duration) string
```

Duration formats d as space-separated days, hours, minutes and seconds,
omitting zero components. Durations of one second or more truncate the
sub-second remainder; shorter non-zero ones use a single ms, µs or ns
component. Negative durations are prefixed with "-"; math.MinInt64 is
handled without overflow.

```go
Duration(3*time.Hour + 25*time.Minute) // "3h 25m"
Duration(90 * time.Second)             // "1m 30s"
Duration(125 * time.Millisecond)       // "125ms"
Duration(0)                            // "0s"
```

## DurationShort

```go
func DurationShort(d time.Duration) string
```

DurationShort is equivalent to Duration, provided for symmetry with
DurationNatural and DurationApprox.

```go
DurationShort(90 * time.Second) // "1m 30s"
```

## DurationLong

```go
func DurationLong(d time.Duration) string
```

DurationLong is like Duration but uses full unit names, pluralised unless
the count is 1, separated by ", ". Zero is "0 seconds".

```go
DurationLong(3*time.Hour + 25*time.Minute + 12*time.Second)
// "3 hours, 25 minutes, 12 seconds"
```

## DurationWithOptions

```go
func DurationWithOptions(d time.Duration, opts DurationOptions) string
```

DurationWithOptions is like Duration, or DurationLong when opts.Long is set,
limited to opts.Units components.

```go
DurationWithOptions(3*24*time.Hour+12*time.Hour+32*time.Minute, DurationOptions{Units: 2}) // "3d 12h"
DurationWithOptions(time.Hour+5*time.Second, DurationOptions{Units: 2})                    // "1h"
DurationWithOptions(3*time.Hour+25*time.Minute+12*time.Second, DurationOptions{Units: 2, Long: true})
// "3 hours, 25 minutes"
```

## DurationOptions

```go
type DurationOptions struct {
	// Units is the maximum number of components shown, counted from the
	// largest non-zero one; smaller components are truncated, never rounded.
	// Zero components inside that window still count. Units <= 0 shows all
	// components. It has no effect below one second.
	Units int
	// Long uses full unit names separated by ", ", as in DurationLong.
	Long bool
}
```

DurationOptions controls DurationWithOptions.

## DurationNatural

```go
func DurationNatural(d time.Duration) string
```

DurationNatural formats d as an English phrase with full unit names, joining
the last two components with "and". Days, hours, minutes and seconds are
used, zero components are omitted and the sub-second remainder is truncated.
Durations below one second use a single millisecond, microsecond or
nanosecond component, like DurationLong. Negative durations are prefixed
with "-"; math.MinInt64 is handled without overflow.

```go
DurationNatural(49*time.Hour + 32*time.Minute) // "2 days, 1 hour and 32 minutes"
DurationNatural(time.Hour + 5*time.Second)     // "1 hour and 5 seconds"
DurationNatural(125 * time.Millisecond)        // "125 milliseconds"
DurationNatural(0)                             // "0 seconds"
```

## DurationApprox

```go
func DurationApprox(d time.Duration) string
```

DurationApprox describes the magnitude of d as a rough English phrase.
The sign of d is ignored. The unit is chosen from the exact magnitude
and the count is rounded to the nearest whole unit, halves rounding up;
a count that rounds up to a full next unit is reported in that unit.

```go
< 1 minute   "less than a minute"
< 1 hour     "about N minutes"  (59m30s is "about 1 hour")
< 1 day      "about N hours"
< 30 days    "about N days"
< 365 days   "about N months"   (30-day months)
otherwise    "about N years"    (365-day years)
```

Examples:

```go
DurationApprox(100 * time.Minute) // "about 2 hours"
DurationApprox(49 * time.Hour)    // "about 2 days"
DurationApprox(30 * time.Second)  // "less than a minute"
```

## Latency

```go
func Latency(d time.Duration) string
```

Latency formats d using the single most appropriate unit (ns, µs, ms, s,
m or h) with at most one fraction digit. It suits logs and metrics.

```go
Latency(350 * time.Microsecond)  // "350µs"
Latency(42 * time.Millisecond)   // "42ms"
Latency(1400 * time.Millisecond) // "1.4s"
Latency(120 * time.Second)       // "2m"
```

## RelativeTime

```go
func RelativeTime(t time.Time) string
```

RelativeTime describes t relative to the current time, e.g. "2 hours ago".
It calls time.Now; prefer RelativeTimeFrom for deterministic output.

## RelativeTimeFrom

```go
func RelativeTimeFrom(now, t time.Time) string
```

RelativeTimeFrom describes t relative to now: "in N units" when t is
after now and "N units ago" otherwise. Counts are truncated towards zero,
the sub-second part of the difference is ignored and units are pluralised
unless the count is 1.

```go
|diff| < 5s     "just now"
< 60s           "N seconds ago"   / "in N seconds"
< 60m           "N minutes ago"
< 24h           "N hours ago"
< 7d            "N days ago"
< 30d           "N weeks ago"
< 365d          "N months ago"    (30-day months)
otherwise       "N years ago"     (365-day years)
```

The difference is computed with time.Time.Sub, which saturates at
math.MaxInt64 and math.MinInt64 nanoseconds, so any gap of more than about
292 years reports "292 years ago" or "in 292 years". Time zones do not
affect the result; only the instants are compared.

## Date

```go
func Date(t time.Time) string
```

Date describes the calendar day of t relative to the current time. It calls
time.Now; prefer DateFrom for deterministic output.

## DateFrom

```go
func DateFrom(now, t time.Time) string
```

DateFrom describes the calendar day of t relative to now. Calendar days are
compared in t's location: now is converted to t.Location() first. Any year,
including years before 1970 and after 9999, is supported.

```
same day          "Today"
previous day      "Yesterday"
next day          "Tomorrow"
same year         "Sep 13"
otherwise         "Sep 13, 2025"
```

## Time

```go
Package time provides functionality for measuring and displaying time.

The calendrical calculations always assume a Gregorian calendar, with no leap
seconds.

# Monotonic Clocks

Operating systems provide both a “wall clock,” which is subject to changes
for clock synchronization, and a “monotonic clock,” which is not. The general
rule is that the wall clock is for telling time and the monotonic clock is for
measuring time. Rather than split the API, in this package the Time returned by
time.Now contains both a wall clock reading and a monotonic clock reading; later
time-telling operations use the wall clock reading, but later time-measuring
operations, specifically comparisons and subtractions, use the monotonic clock
reading.

For example, this code always computes a positive elapsed time of approximately
20 milliseconds, even if the wall clock is changed during the operation being
timed:
```

start := time.Now()
... operation that takes 20 milliseconds ...
t := time.Now()
elapsed := t.Sub(start)

Other idioms, such as time.Since(start), time.Until(deadline), and
time.Now().Before(deadline), are similarly robust against wall clock resets.

The rest of this section gives the precise details of how operations use
monotonic clocks, but understanding those details is not required to use this
package.

The Time returned by time.Now contains a monotonic clock reading. If Time t
has a monotonic clock reading, t.Add adds the same duration to both the wall
clock and monotonic clock readings to compute the result. Because t.AddDate(y,
m, d), t.Round(d), and t.Truncate(d) are wall time computations, they always
strip any monotonic clock reading from their results. Because t.In, t.Local,
and t.UTC are used for their effect on the interpretation of the wall time,
they also strip any monotonic clock reading from their results. The canonical
way to strip a monotonic clock reading is to use t = t.Round(0).

If Times t and u both contain monotonic clock readings, the operations
t.After(u), t.Before(u), t.Equal(u), t.Compare(u), and t.Sub(u) are carried out
using the monotonic clock readings alone, ignoring the wall clock readings.
If either t or u contains no monotonic clock reading, these operations fall back
to using the wall clock readings.

On some systems the monotonic clock will stop if the computer goes to sleep.
On such a system, t.Sub(u) may not accurately reflect the actual time that
passed between t and u. The same applies to other functions and methods that
subtract times, such as Since, Until, Time.Before, Time.After, Time.Add,
Time.Equal and Time.Compare. In some cases, you may need to strip the monotonic
clock to get accurate results.

Because the monotonic clock reading has no meaning outside the current process,
the serialized forms generated by t.GobEncode, t.MarshalBinary, t.MarshalJSON,
and t.MarshalText omit the monotonic clock reading, and t.Format provides
no format for it. Similarly, the constructors time.Date, time.Parse,
time.ParseInLocation, and time.Unix, as well as the unmarshalers t.GobDecode,
t.UnmarshalBinary. t.UnmarshalJSON, and t.UnmarshalText always create times with
no monotonic clock reading.

The monotonic clock reading exists only in Time values. It is not a part of
Duration values or the Unix times returned by t.Unix and friends.

Note that the Go == operator compares not just the time instant but also the
Location and the monotonic clock reading. See the documentation for the Time
type for a discussion of equality testing for Time values.

For debugging, the result of t.String does include the monotonic clock
reading if present. If t != u because of different monotonic clock readings,
that difference will be visible when printing t.String() and u.String().

# Timer Resolution

Timer resolution varies depending on the Go runtime, the operating system
and the underlying hardware. On Unix, the resolution is ~1ms. On Windows
version 1803 and newer, the resolution is ~0.5ms. On older Windows versions,
the default resolution is ~16ms, but a higher resolution may be requested using
golang.org/x/sys/windows.TimeBeginPeriod.

const Layout = "01/02 03:04:05PM '06 -0700" ...
const Nanosecond Duration = 1 ...
func After(d Duration) <-chan Time
func Sleep(d Duration)
func Tick(d Duration) <-chan Time
type Duration int64
func ParseDuration(s string) (Duration, error)
func Since(t Time) Duration
func Until(t Time) Duration
type Location struct{ ... }
var Local *Location = &localLoc
var UTC *Location = &utcLoc
func FixedZone(name string, offset int) *Location
func LoadLocation(name string) (*Location, error)
func LoadLocationFromTZData(name string, data []byte) (*Location, error)
type Month int
const January Month = 1 + iota ...
type ParseError struct{ ... }
type Ticker struct{ ... }
func NewTicker(d Duration) *Ticker
type Time struct{ ... }
func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time
func Now() Time
func Parse(layout, value string) (Time, error)
func ParseInLocation(layout, value string, loc *Location) (Time, error)
func Unix(sec int64, nsec int64) Time
func UnixMicro(usec int64) Time
func UnixMilli(msec int64) Time
type Timer struct{ ... }
func AfterFunc(d Duration, f func()) *Timer
func NewTimer(d Duration) *Timer
type Weekday int
const Sunday Weekday = iota ...

## TimeFrom

```go
func TimeFrom(now, t time.Time) string
```

TimeFrom returns RelativeTimeFrom(now, t), a middle dot and t's date and
24-hour time in t's location. The year is included only when it differs from
now's year in that location.

```
"2 hours ago · Sep 13, 07:20"
"1 year ago · Sep 13, 2025, 07:20"
```

## TimeRange

```go
func TimeRange(start, end time.Time) string
```

TimeRange formats the interval between start and end relative to the current
time. It calls time.Now; prefer TimeRangeFrom for deterministic output.

## TimeRangeFrom

```go
func TimeRangeFrom(now, start, end time.Time) string
```

TimeRangeFrom formats the interval between start and end using a 24-hour
clock. If end is before start the two are swapped. Both ends are shown
in start's location, and "today" and "this year" are judged against now
converted to that location.

```
same day, today        "09:30–11:45"
same day               "Sep 20, 09:30–11:45"
different days         "Sep 13, 23:30 → Sep 14, 02:15"
different years        "Dec 31, 2026, 23:30 → Jan 1, 2027, 02:15"
```

The year is shown on both dates when start and end fall in different years
or when either differs from now's year. The en dash is U+2013.

## ParseDuration

```go
func ParseDuration(s string) (time.Duration, error)
```

ParseDuration parses a duration as printed by Duration, DurationLong,
DurationNatural or Latency: "3d 12h", "2h30m", "1.5ms", "350µs", "2 days,
1 hour and 32 minutes". Components are a number and a unit (ns, µs/us, ms,
s, m, h, d or their English names, case-insensitive), separated by optional
spaces, commas or "and". A day is 24 hours. A leading "-" negates the whole
duration. Results are rounded to the nearest nanosecond, half away from
zero.

```go
ParseDuration("3h 25m") // 3h25m0s, nil
ParseDuration("1.5ms")  // 1.5ms, nil
```

Errors wrap ErrSyntax, ErrUnit or ErrRange (outside time.Duration).

## ETA

```go
func ETA(done, total int64, elapsed time.Duration) string
```

ETA estimates the time left for a task that has done of total units after
elapsed, assuming a constant rate, as "~3m left". The estimate keeps the
largest unit only and is at least one second. It returns "done" when done
>= total > 0 and "" when there is no rate yet (done, total or elapsed not
positive).

```go
ETA(25, 100, time.Minute)   // "~3m left"
ETA(100, 100, time.Minute)  // "done"
ETA(0, 100, time.Minute)    // ""
```

## Calendar

```go
func Calendar(now, t time.Time) string
```

Calendar describes t for schedules and feeds, relative to the calendar day
of now in t's location:

```
same day          "today 14:30"
previous day      "yesterday 14:30"
next day          "tomorrow 14:30"
within 6 days     "Monday 14:30"
same year         "Sep 2, 14:30"
otherwise         "Sep 2, 2025, 14:30"
```

## Locale

```go
type Locale string
```

Locale selects the language of the localised formatters. Unknown locales
fall back to English. Use ParseLocale to map a BCP 47 tag such as "ru-RU".

const English Locale = "en" ...
func ParseLocale(tag string) (Locale, bool)
func (l Locale) DurationLong(d time.Duration) string
func (l Locale) Int(n int64) string
func (l Locale) List(items []string) string
func (l Locale) Percent(ratio float64) string
func (l Locale) Plural(n uint64, forms ...string) string
func (l Locale) RelativeTime(t time.Time) string
func (l Locale) RelativeTimeFrom(now, t time.Time) string

## ParseLocale

```go
func ParseLocale(tag string) (Locale, bool)
```

ParseLocale maps a BCP 47 language tag to a supported Locale, matching
case-insensitively and accepting '_' for '-'. Region subtags are ignored;
"uz-Cyrl" selects UzbekCyrillic. Unsupported tags return English, false.

```go
ParseLocale("ru-RU")   // Russian, true
ParseLocale("uz_Cyrl") // UzbekCyrillic, true
ParseLocale("xx")      // English, false
```

## PluralRU

```go
func PluralRU(n uint64, one, few, many string) string
```

PluralRU picks the Russian plural form for n: one (1, 21, 101...), few (2–4,
22–24...) or many (0, 5–20, 25–30...).

```go
PluralRU(21, "файл", "файла", "файлов") // "файл"
PluralRU(3, "файл", "файла", "файлов")  // "файла"
PluralRU(11, "файл", "файла", "файлов") // "файлов"
```
