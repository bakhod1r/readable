---
title: Text
layout: default
nav_order: 8
---

# Text
{: .no_toc }

Labels from identifiers, English lists and booleans.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## Humanize

```go
func Humanize(s string) string
```

Humanize turns an identifier into a sentence-case label. Words are split on
'_', '-', Unicode white space and camelCase boundaries: an upper-case rune
starts a new word when the previous rune is not upper case or the next rune
is lower case, so "HTTPServer" splits as "HTTP" and "Server". The first
word is capitalised and the rest lower-cased, except the acronyms API, CPU,
DNS, HTTP, HTTPS, ID, IP, JSON, JWT, OK, RAM, SQL, SSL, TCP, TLS, UDP, UI,
URI, URL, UUID and XML, which are upper-cased wherever they appear (matched
case-insensitively). Each run of invalid UTF-8 bytes becomes a single
U+FFFD. An input with no words returns "".

```go
Humanize("created_at")        // "Created at"
Humanize("HTTP_SERVER_ERROR") // "HTTP server error"
Humanize("userCreatedAt")     // "User created at"
Humanize("user_id")           // "User ID"
```

## Enum

```go
func Enum(s string) string
```

Enum turns a SCREAMING_SNAKE_CASE enum value into a sentence-case label.
It is equivalent to Humanize.

```go
Enum("PAYMENT_COMPLETED") // "Payment completed"
Enum("IN_PROGRESS")       // "In progress"
```

## List

```go
func List(items []string) string
```

List joins items as an English list using the Oxford comma. An empty list
returns "".

```go
List([]string{"Go"})                        // "Go"
List([]string{"Go", "PostgreSQL"})          // "Go and PostgreSQL"
List([]string{"Go", "PostgreSQL", "Redis"}) // "Go, PostgreSQL, and Redis"
```

## ListWithOptions

```go
func ListWithOptions(items []string, opts ListOptions) string
```

ListWithOptions is like List with explicit options. When Limit is positive
and smaller than len(items), only the first Limit items are shown and the
remainder is summarised as a final "N more" element, which takes part in
the Oxford-comma layout like any other item. Two elements are joined as "A
conj B" without a comma. The result is built with at most one allocation;
a single item is returned as is.

```go
ListWithOptions(five, ListOptions{Limit: 3})         // "Go, Redis, Kafka, and 2 more"
ListWithOptions(two, ListOptions{Conjunction: "or"}) // "Go or Redis"
```

## ListOptions

```go
type ListOptions struct {
	// Limit is the maximum number of items shown; the rest are summarised
	// as "N more". Zero or negative means unlimited.
	Limit int
	// Conjunction joins the last item. Empty means "and".
	Conjunction string
	// Empty is returned for an empty list. Default "".
	Empty string
}
```

ListOptions configures ListWithOptions.

## Bool

```go
func Bool(b bool) string
```

Bool returns "Yes" for true and "No" for false.

```go
Bool(true) // "Yes"
```

## BoolLabel

```go
func BoolLabel(b bool, trueLabel, falseLabel string) string
```

BoolLabel returns trueLabel if b is true, otherwise falseLabel.

```go
BoolLabel(false, "Enabled", "Disabled") // "Disabled"
```
