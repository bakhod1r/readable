// Package readable converts machine-oriented values into short,
// human-readable strings for logs, CLIs, dashboards and user interfaces.
//
// Every function takes plain Go values (int64, uint64, float64,
// [time.Duration], [time.Time], string) and returns a string. There is no
// configuration state, no locale database and no dependency outside the
// standard library. Output is English and deterministic.
//
//	readable.Number(1234567)              // "1.23M"
//	readable.Bytes(1536)                  // "1.5 KB"
//	readable.Duration(90 * time.Second)   // "1m 30s"
//	readable.Money(150050, "USD")         // "1,500.50 USD"
//	readable.MaskCard("8600123412345678") // "**** **** **** 5678"
//
// # Numbers
//
// [Number] compacts integers with K, M, B and T ("12.5K"); [NumberFloat]
// does the same for float64. [NumberWithPrecision] and [NumberWithOptions]
// control fraction digits, fixed precision and the decimal separator.
// [NumberWords], [Ordinal], [Count], [CountPlural] and [Plural] cover
// wording ("1st", "3 files"). [ID] groups long numeric identifiers in threes.
//
// # Bytes
//
// [Bytes] uses 1024-based multiples with the conventional KB/MB labels,
// [BytesIEC] uses KiB/MiB and [BytesSI] uses 1000-based kB/MB. [FileSize]
// accepts the int64 returned by os.FileInfo.Size. Transfer speeds use
// [Throughput], [ThroughputIEC], [BytesRate] and [Bandwidth] (bits per
// second).
//
// # Time
//
// [Duration], [DurationShort], [DurationLong], [DurationNatural] and
// [DurationApprox] describe a [time.Duration] at different lengths;
// [DurationWithOptions] limits the number of components. [Latency] picks a
// single unit for metrics ("1.4s").
//
// [RelativeTime], [Date], [Time] and [TimeRange] read the current clock via
// [time.Now]. Each has a From variant ([RelativeTimeFrom], [DateFrom],
// [TimeFrom], [TimeRangeFrom]) that takes an explicit now for deterministic
// output and tests.
//
// # Money
//
// Money amounts are int64 values in the currency's minor units (cents for
// USD), never floats. [Money] appends the code, [MoneySymbol] prefixes a known
// symbol, and [MoneyCompact], [MoneyAccounting] and [MoneyChange] cover
// dashboards, ledgers and deltas. Unknown currency codes are accepted with 0
// minor digits; [MoneyWithPrecision] sets the digits explicitly.
//
// Ratios and rates: [Percent], [PercentWithPrecision], [PercentChange],
// [Progress], [ProgressBar], [Rate], [RateWithLabel], [PerMinute] and
// [RequestRate].
//
// # Text
//
// [Humanize] and [Enum] turn identifiers such as "user_id" or "IN_PROGRESS"
// into labels ("User ID", "In progress"). [List] and [ListWithOptions] join
// items as an English list. [Truncate] shortens a string around "...".
// [Bool] and [BoolLabel] render booleans.
//
// # Masking
//
// [Mask], [MaskCard], [MaskEmail], [MaskPhone], [MaskIP] and [MaskToken] hide
// sensitive values for display and logging. Each documents exactly what it
// reveals, and each fails closed: input it cannot interpret is fully replaced
// by a placeholder such as "***" or "****" rather than partially shown.
//
// [Truncate], [Hash] and [ShortUUID] are for readability, not secrecy; they
// deliberately reveal leading and trailing characters.
//
// # Rounding and precision
//
// Integer formatters ([Number], [Bytes], [MoneyCompact], ...) round
// half away from zero using exact integer arithmetic, so there are no binary
// floating-point surprises. Rounding may promote a value to the next unit:
// Number(999999) is "1M", not "1000K". Values at the top of the range stay in
// the largest unit, and math.MinInt64 is handled without overflow.
//
// Float formatters ([NumberFloat], [Percent], [Rate]) use strconv's correctly
// rounded conversion of the binary float value; for example 0.1+0.2 is not
// exactly 0.3. NaN and infinities format as "NaN", "+Inf" and "-Inf", followed by any
// unit suffix ("NaN%").
//
// Precision is a maximum number of fraction digits, clamped to
// [0, [MaxPrecision]] and defaulting to [DefaultPrecision]. Trailing zeros are
// trimmed ("1.00M" becomes "1M") unless fixed precision is requested with
// [NumberOptions].
//
// # Concurrency
//
// All functions are pure and safe for concurrent use. Package-level tables
// are never modified after initialization.
//
// # Stability
//
// The package follows Semantic Versioning, and output strings are part of the
// API. Before v1.0.0, minor releases may change output or signatures; every
// such change is listed in CHANGELOG.md. From v1.0.0, documented output only
// changes in a new major version.
package readable
