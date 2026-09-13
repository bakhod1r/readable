# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
Output strings are part of the public API.

## [Unreleased]

## [0.1.0] - 2026-09-13

Initial release. Standard library only, Go 1.24+.

### Added

- **Numbers**: `Number`, `NumberFloat`, `NumberWithPrecision`,
  `NumberWithOptions` with `NumberOptions`, `NumberWords`, `Ordinal`,
  `Count`, `CountPlural`, `Plural`, constants `DefaultPrecision` and
  `MaxPrecision`.
- **Bytes & throughput**: `Bytes`, `BytesIEC`, `BytesSI`, `FileSize`,
  `BytesRate`, `Throughput`, `ThroughputIEC`, `Bandwidth`.
- **Durations & time**: `Duration`, `DurationShort`, `DurationLong`,
  `DurationWithOptions`, `DurationApprox`,
  `DurationNatural`, `Latency`, `RelativeTime`, `RelativeTimeFrom`, `Date`,
  `DateFrom`, `Time`, `TimeFrom`, `TimeRange`, `TimeRangeFrom`.
- **Money** (integer minor units): `Money`, `MoneyWithPrecision`,
  `MoneySymbol`, `MoneyCompact`, `MoneyAccounting`, `MoneyChange`.
- **Percent & progress**: `Percent`, `PercentWithPrecision`,
  `PercentChange` with `ErrUndefined`, `Progress`, `ProgressBar`, `Rate`,
  `RateWithLabel`, `PerMinute`, `RequestRate`.
- **Text**: `Humanize`, `Enum`, `List`, `ListWithOptions` with
  `ListOptions`, `Truncate`, `Bool`, `BoolLabel`.
- **Masking & IDs**: `Mask`, `MaskCard`, `MaskEmail`, `MaskPhone`, `MaskIP`,
  `MaskToken`, `ID`, `Hash`, `ShortUUID`.

[Unreleased]: https://github.com/bakhod1r/readable/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/bakhod1r/readable/releases/tag/v0.1.0
