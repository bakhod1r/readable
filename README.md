# readable

Small, fast, dependency-free human-readable formatting for Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/bakhod1r/readable.svg)](https://pkg.go.dev/github.com/bakhod1r/readable)
[![CI](https://github.com/bakhod1r/readable/actions/workflows/ci.yml/badge.svg)](https://github.com/bakhod1r/readable/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bakhod1r/readable)](https://goreportcard.com/report/github.com/bakhod1r/readable)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

`readable` turns machine values into short strings for dashboards, CLIs, logs
and notifications: `1.5M`, `1.5 KB`, `3h 25m`, `$1,500.50`, `2 hours ago`,
`Go, Redis, and Kafka`, `j***@gmail.com`.

**📖 Documentation: [bakhod1r.github.io/readable](https://bakhod1r.github.io/readable/)**

## Install

```sh
go get github.com/bakhod1r/readable
```

Requires Go 1.24 or newer. Zero dependencies.

## Quick start

```go
readable.Number(1_500_000)                      // 1.5M
readable.Bytes(1_048_576)                       // 1 MB
readable.Duration(3*time.Hour + 25*time.Minute) // 3h 25m
readable.Money(150050, "USD")                   // 1,500.50 USD
readable.Percent(0.1534)                        // 15.34%
readable.RelativeTimeFrom(now, t)               // 2 hours ago
readable.List([]string{"Go", "Redis", "Kafka"}) // Go, Redis, and Kafka
readable.MaskEmail("john.doe@gmail.com")        // j***@gmail.com
```

## What's inside

| Area | Functions |
|---|---|
| [Numbers](https://bakhod1r.github.io/readable/numbers/) | `Number`, `NumberFloat`, `NumberWords`, `Ordinal`, `Count`, `Plural` |
| [Bytes & throughput](https://bakhod1r.github.io/readable/bytes/) | `Bytes`, `BytesIEC`, `BytesSI`, `FileSize`, `Throughput`, `Bandwidth` |
| [Durations & time](https://bakhod1r.github.io/readable/time/) | `Duration`, `DurationNatural`, `Latency`, `RelativeTime`, `Date`, `TimeRange` |
| [Money](https://bakhod1r.github.io/readable/money/) | `Money`, `MoneySymbol`, `MoneyCompact`, `MoneyAccounting`, `MoneyChange` |
| [Percent & rates](https://bakhod1r.github.io/readable/percent/) | `Percent`, `PercentChange`, `Progress`, `ProgressBar`, `Rate` |
| [Text](https://bakhod1r.github.io/readable/text/) | `Humanize`, `Enum`, `List`, `Bool` |
| [Masking & IDs](https://bakhod1r.github.io/readable/masking/) | `MaskEmail`, `MaskPhone`, `MaskCard`, `MaskToken`, `MaskIP`, `ID`, `Truncate` |

## Why readable

- **Exact rounding** with integer math — `999_999` becomes `1M`, never `1000K`.
- **Fail-closed masking** — malformed input is fully masked, never partially revealed.
- **Predictable** — pure, concurrency-safe, full `int64` range, NaN/Inf handled.
- **Tested** — 100% coverage, fuzz targets, runnable examples.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md), [CHANGELOG.md](CHANGELOG.md) and [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE)
