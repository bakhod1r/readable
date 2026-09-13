---
title: Numbers
layout: default
nav_order: 3
---

# Numbers
{: .no_toc }

Compact integers and floats, spell scales, ordinals and counted nouns.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## Number

```go
func Number(n int64) string
```

Number formats n using compact decimal units K, M, B and T with up to
DefaultPrecision fraction digits, rounding half away from zero.

```go
Number(999)     // "999"
Number(12500)   // "12.5K"
Number(1234567) // "1.23M"
Number(999999)  // "1M" (rounding promotes to the next unit)
```

Values of a quadrillion or more stay in T ("9223372.04T" for MaxInt64).

## NumberFloat

```go
func NumberFloat(f float64) string
```

NumberFloat formats f like Number with up to DefaultPrecision fraction
digits, rounded to nearest by strconv and trimmed of trailing zeros.
Values that round to 1000 of a unit move to the next one ("1K" for 999.995).
Results that round to zero lose their sign ("0" for -0.001). Values of a
quadrillion or more stay in T with all integer digits. NaN and infinities
format as "NaN", "+Inf" and "-Inf".

```go
NumberFloat(0.5)     // "0.5"
NumberFloat(1500.25) // "1.5K"
```

## NumberWithPrecision

```go
func NumberWithPrecision(n int64, precision int) string
```

NumberWithPrecision is like Number with at most precision fraction digits.

```go
NumberWithPrecision(1234567, 1) // "1.2M"
NumberWithPrecision(1234567, 3) // "1.235M"
```

## NumberWithOptions

```go
func NumberWithOptions(n int64, opts NumberOptions) string
```

NumberWithOptions is like Number with explicit options.

## NumberOptions

```go
type NumberOptions struct {
	// Precision is the maximum number of fraction digits (clamped to
	// [0, MaxPrecision]).
	Precision int
	// FixedPrecision keeps trailing zeros ("1.50M" instead of "1.5M").
	FixedPrecision bool
	// Decimal is the decimal separator. Empty means ".".
	Decimal string
}
```

NumberOptions configures NumberWithOptions.

## NumberWords

```go
func NumberWords(n int64) string
```

NumberWords is like Number but spells the scale as a short-scale English
word: thousand, million, billion, trillion. At most DefaultPrecision
fraction digits are kept, rounded half up and trimmed.

```go
NumberWords(999)     // "999"
NumberWords(12500)   // "12.5 thousand"
NumberWords(1234567) // "1.23 million"
```

## Ordinal

```go
func Ordinal(n int64) string
```

Ordinal formats n as an English ordinal number: the suffix is "th" when the
last two digits of |n| are 11, 12 or 13, otherwise "st", "nd" or "rd" for a
last digit of 1, 2 or 3 and "th" for anything else. Negative numbers keep
their sign; math.MinInt64 is handled without overflow.

```go
Ordinal(1)   // "1st"
Ordinal(12)  // "12th"
Ordinal(23)  // "23rd"
Ordinal(-2)  // "-2nd"
Ordinal(111) // "111th"
```

## Count

```go
func Count(n int64, singular string) string
```

Count formats n, a space and the singular or regular plural form of the
noun, without digit grouping. The singular is used when |n| == 1.

```go
Count(1, "file")    // "1 file"
Count(0, "file")    // "0 files"
Count(1500, "user") // "1500 users"
```

See Plural for the pluralisation rules.

## CountPlural

```go
func CountPlural(n int64, singular, plural string) string
```

CountPlural is like Count with an explicit plural form, used verbatim when
|n| != 1.

```go
CountPlural(2, "person", "people") // "2 people"
```

## Plural

```go
func Plural(n int64, singular string) string
```

Plural returns singular when |n| == 1 and its regular English plural
otherwise. Rules, with the ending matched case-insensitively: a word ending
in s, x, z, ch or sh gets "es"; a y preceded by a consonant becomes "ies";
anything else gets "s". The suffix is upper case when the word contains an
upper-case letter and no lower-case one, so "CITY" becomes "CITIES". The
empty string stays empty. Irregular nouns are not handled; use CountPlural
for those.

```go
Plural(1, "user")  // "user"
Plural(2, "match") // "matches"
Plural(2, "city")  // "cities"
Plural(2, "BOX")   // "BOXES"
```
