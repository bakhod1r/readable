---
title: Design notes
layout: default
nav_order: 11
---

# Design notes
{: .no_toc }

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## Rounding

Integer formatters (`Number`, `Bytes`, `MoneyCompact`, ...) round half away
from zero using exact integer arithmetic, so there are no binary floating-point
surprises. Rounding may promote a value to the next unit: `Number(999999)` is
`1M`, not `1000K`. Values at the top of the range stay in the largest unit, and
`math.MinInt64` is handled without overflow.

Float formatters (`NumberFloat`, `Percent`, `Rate`) use `strconv`'s correctly
rounded conversion of the binary float value; for example `0.1+0.2` is not
exactly `0.3`. NaN and infinities format as `NaN`, `+Inf` and `-Inf`, followed
by any unit suffix (`NaN%`).

## Precision

Precision is a maximum number of fraction digits, clamped to
`[0, MaxPrecision]` and defaulting to `DefaultPrecision`. Trailing zeros are
trimmed (`1.00M` becomes `1M`) unless fixed precision is requested with
`NumberOptions`.

```go
const DefaultPrecision = 2 // fraction digits for formatters without explicit precision
const MaxPrecision = 9     // larger values are clamped
```

## Errors

Functions return `string`. Only results that can be undefined return
`(string, error)`, using the sentinel `ErrUndefined`:

```go
if _, err := readable.PercentChange(0, n); errors.Is(err, readable.ErrUndefined) {
	// change from zero has no percentage
}
```

No function panics on any input.

## Time and determinism

`RelativeTime`, `Date`, `Time` and `TimeRange` read `time.Now`. Each has a
`...From(now, ...)` twin that takes an explicit `now` for deterministic output
and tests. Everything else is a pure function of its arguments.

## Concurrency

All functions are pure and safe for concurrent use. Package-level tables are
never modified after initialization.

## Non-goals

- **Localization.** Output is English with `,` and `.` separators.
- **Decimal or accounting arithmetic.** `readable` formats amounts; it does not compute them.
- **Parsing.** Strings go out, not in.
- **Business logic.** Currency conversion, time zones and policy stay in your code.

## Stability

The package follows [Semantic Versioning](https://semver.org), and output
strings are part of the API. Before v1.0.0, minor releases may change output or
signatures; every such change is listed in the [changelog](changelog). From
v1.0.0, documented output only changes in a new major version.
