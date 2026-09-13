# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). Output
strings are part of the public API.

## [Unreleased]

### Added

- Numbers: `Number`, `NumberWithPrecision`, `NumberWithOptions`, `NumberFloat`,
  `NumberWords`, `Ordinal`, `Count`, `CountPlural`.
- Bytes and throughput formatters.
- Duration, natural duration, date and relative time formatters.
- Money, percent and progress formatters.
- Text: lists and pluralization.
- Sensitive values and IDs: `Mask`, `MaskEmail`, `MaskPhone`, `MaskCard`,
  `MaskToken`, `MaskIP`, `ID`, `Truncate`, `Hash`, `ShortUUID`.
- Logs: `LogDuration`, `LogBytes`, `LogNumber`.
- CI: gofmt, `go mod tidy`, vet, staticcheck, race tests, fuzzing, benchmarks.
