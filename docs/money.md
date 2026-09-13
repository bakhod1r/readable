---
title: Money
layout: default
nav_order: 6
---

# Money
{: .no_toc }

Amounts are `int64` **minor units** (cents for USD), never floats.

| Minor digits | Currencies |
|---|---|
| 2 | most ISO 4217 codes: USD, EUR, GBP, CHF, RUB, KZT, CNY, … |
| 3 | BHD, IQD, JOD, KWD, LYD, OMR, TND |
| 4 | CLF, UYW |
| 0 | JPY, KRW, VND, UZS, … and any unknown code |

Unknown codes are accepted; pick the digits explicitly with `MoneyWithPrecision`.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## Money

```go
func Money(amount int64, currency string) string
```

Money formats amount, expressed in the currency's minor units,
with thousands separators followed by the currency code. Using integer
minor units avoids floating-point error. Codes are matched ASCII
case-insensitively but printed as given.

Active ISO 4217 codes use their standard number of minor digits (USD and
CHF have 2, JPY has 0, BHD has 3), except UZS, which has 0. Any other code
(for example "USDT") is accepted and treated as having 0 minor digits; use
MoneyWithPrecision to choose explicitly.

```go
Money(150000000, "UZS") // "150,000,000 UZS"
Money(150050, "USD")    // "1,500.50 USD"
```

## MoneySymbol

```go
func MoneySymbol(amount int64, currency string) string
```

MoneySymbol is like Money but prefixes the currency symbol when one is known
("$1,500.50"). Currencies without a symbol fall back to the code.

## MoneyWithPrecision

```go
func MoneyWithPrecision(amount int64, currency string, minorDigits int) string
```

MoneyWithPrecision is like Money but amount has exactly minorDigits minor
digits (clamped to [0, MaxPrecision]), whatever the currency.

```go
MoneyWithPrecision(1234567, "USDT", 4) // "123.4567 USDT"
```

## MoneyCompact

```go
func MoneyCompact(amount int64, currency string) string
```

MoneyCompact formats amount, expressed in the currency's minor units, with
the major value compacted like Number (K, M, B, T, at most DefaultPrecision
trimmed fraction digits, rounded half up using integer math). A known symbol
is prefixed; otherwise the code is appended. Major values below 1000 are
formatted like MoneySymbol.

```go
MoneyCompact(150000000, "USD") // "$1.5M"
MoneyCompact(1500000, "UZS")   // "1.5M UZS"
MoneyCompact(99999, "USD")     // "$999.99"
```

## MoneyAccounting

```go
func MoneyAccounting(amount int64, currency string) string
```

MoneyAccounting is like Money but renders negative amounts in accounting
style, wrapped in parentheses without a minus sign.

```go
MoneyAccounting(-1500000, "UZS") // "(1,500,000 UZS)"
MoneyAccounting(1500000, "UZS")  // "1,500,000 UZS"
```

## MoneyChange

```go
func MoneyChange(amount int64, currency string) string
```

MoneyChange is like Money but always shows the direction of change: positive
amounts get "+", negative amounts get the typographic minus sign U+2212
("−"), and zero has no sign.

```go
MoneyChange(1500000, "UZS")  // "+1,500,000 UZS"
MoneyChange(-1500000, "UZS") // "−1,500,000 UZS"
```
