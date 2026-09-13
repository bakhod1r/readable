---
title: Parsing & templates
layout: default
nav_order: 10
---

# Parsing & templates
{: .no_toc }

Parse functions read formatted values back and return errors wrapping `ErrSyntax`, `ErrUnit` or `ErrRange`. `FuncMap` exposes the formatters to `text/template` and `html/template`.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## ErrSyntax

```go
var (
	// ErrSyntax reports input that is not a number followed by a unit.
	ErrSyntax = errors.New("invalid syntax")
	// ErrUnit reports a missing or unknown unit.
	ErrUnit = errors.New("unknown unit")
	// ErrRange reports a value that does not fit the result type.
	ErrRange = errors.New("value out of range")
)
```

Errors returned (wrapped) by the Parse functions. Test with errors.Is.

## FuncMap

```go
func FuncMap() map[string]any
```

FuncMap returns the package's formatters by lowercase-first name for use
with text/template and html/template:

```go
tmpl := template.New("").Funcs(readable.FuncMap())
// {{ bytes .Size }} · {{ duration .Elapsed }} · {{ money .Amount "USD" }}
```

The map is typed map[string]any so it converts to either package's FuncMap
without importing them. Each call returns a new map. Only single-result
formatters are included; PercentChange and the Parse functions are not.

## Format

```go
func Format(v any) map[string]string
```

Format formats the exported fields of a struct (or pointer to struct) into a
map keyed by field name, for admin pages, CLIs and API views. The `readable`
struct tag picks the formatter:

```go
bytes      Bytes            (unsigned or signed integers)
duration   Duration         (time.Duration or integer nanoseconds)
number     Number           (integers)
percent    Percent          (floats)
money=USD  Money            (integer minor units)
mask=K     MaskEmail, MaskPhone, MaskCard, MaskToken, MaskIP or Mask(v, 0, 0)
           for K = email, phone, card, token, ip or all (strings)
-          field omitted
```

Untagged fields and fields whose tag does not fit their type use fmt.Sprint,
except that an unknown mask kind masks everything. A nil pointer or
non-struct returns nil.

```go
type User struct {
	Email string `readable:"mask=email"`
	Quota uint64 `readable:"bytes"`
}
Format(User{"john.doe@gmail.com", 1536}) // map[Email:j***@gmail.com Quota:1.5 KB]
```

## RedactAttr

```go
func RedactAttr(_ []string, a slog.Attr) slog.Attr
```

RedactAttr is a slog.HandlerOptions.ReplaceAttr function that masks
sensitive attributes before they are written:

  - keys containing password, passwd, secret, token, apikey, authorization,
```go
cookie or privatekey (case-insensitive, ignoring "_" and "-"): "****"
```

  - keys email, phone, card or pan, and ip: MaskEmail, MaskPhone, MaskCard
```
and MaskIP
```

  - every other string value: Redact

Non-string values under sensitive keys are replaced too; other non-string
values are kept.

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout,
	&slog.HandlerOptions{ReplaceAttr: readable.RedactAttr}))
```
