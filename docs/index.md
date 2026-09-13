---
title: Home
layout: default
nav_order: 1
---

# readable
{: .no_toc }

**Turn raw values into clean human-readable strings — zero dependencies, pure stdlib.**
{: .fs-6 .fw-300 }

[Get started](getting-started){: .btn .btn-primary .mr-2 }
[GitHub](https://github.com/bakhod1r/readable){: .btn }

---

```go
readable.Number(1_500_000)          // 1.5M
readable.Bytes(1_048_576)           // 1 MB
readable.Percent(0.1534)            // 15.34%
readable.Duration(3*time.Hour + 25*time.Minute) // 3h 25m
readable.Money(1_500_000, "UZS")    // 1,500,000 UZS
readable.RelativeTimeFrom(now, t)   // 2 hours ago
readable.List([]string{"Go", "Redis", "Kafka"}) // Go, Redis, and Kafka
readable.MaskEmail("john.doe@gmail.com")        // j***@gmail.com
```

## Why readable

- **Zero dependencies.** Standard library only.
- **Exact rounding.** Integer formatters use integer math, half away from zero.
  No `999.99999K`, no float drift; `999_999` becomes `1M`.
- **Predictable output.** Pure functions, stable formats, full `int64` range
  (including `math.MinInt64`), NaN/Inf handled.
- **Fail-closed masking.** Malformed or short sensitive input is fully masked,
  never partially revealed. See [Masking & IDs](masking).
- **Fast.** Hot paths do one allocation (the returned string).
- **Concurrency-safe.** No global mutable state.
- **Tested.** 100% test coverage, fuzz targets, runnable examples.

## Non-goals

- **Localization.** Output is English with `,` and `.` separators.
- **Decimal or accounting arithmetic.** `readable` formats amounts; it does not compute them.
- **Parsing.** Strings go out, not in.
- **Business logic.** Currency conversion, time zones and policy stay in your code.

## License

[MIT](https://github.com/bakhod1r/readable/blob/main/LICENSE) © [bakhod1r](https://github.com/bakhod1r)
