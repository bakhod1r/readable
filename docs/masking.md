---
title: Masking & IDs
layout: default
nav_order: 9
---

# Masking & IDs
{: .no_toc }

Maskers hide sensitive values for display and logging. Each documents exactly what it reveals and **fails closed**: input it cannot interpret is fully replaced by a placeholder.

| Function | Revealed | Invalid or too short |
|---|---|---|
| `MaskEmail` | domain; first rune of local part if it has ≥ 3 runes | `***@domain` or `***` |
| `MaskPhone` | ≤ half the digits, ≤ 3 trailing; digit count | all `*` or `***` |
| `MaskCard` | last 4 digits; digit count (12–19 digits only) | `****` |
| `MaskToken` | known vendor prefix; last 4 runes if secret ≥ 12 runes | `****` |
| `MaskIP` | first two IPv4 octets / IPv6 groups | `***` |
| `Mask` | `keepStart` + `keepEnd` runes; rune count | all `*` |

{: .warning }
`Truncate`, `Hash` and `ShortUUID` are for readability, not secrecy — they deliberately reveal leading and trailing characters.

<details open markdown="block">
<summary>On this page</summary>
{: .text-delta }
- TOC
{:toc}
</details>

## MaskEmail

```go
func MaskEmail(email string) string
```

MaskEmail masks the local part of an email address and keeps the domain.
The last '@' separates local part and domain.

Revealed: the whole domain, plus the first rune of the local part only when
the local part has 3 or more runes. Everything else in the local part is
replaced by a fixed "***", so its length is not leaked. A local part of 1 or
2 runes becomes "***@domain". Input without '@', with an empty local part or
domain, or whose local part starts with invalid UTF-8, returns "***".

```go
MaskEmail("john.doe@gmail.com") // "j***@gmail.com"
MaskEmail("jo@gmail.com")       // "***@gmail.com"
MaskEmail("not-an-email")       // "***"
```

## MaskPhone

```go
func MaskPhone(phone string) string
```

MaskPhone masks a phone number digit by digit. Spaces, '-', '.', '(' and ')'
are removed first; a single leading '+' is kept. Any other character makes
the whole input unmaskable and "***" is returned.

Revealed: at most half of the digits and never more than 3 trailing digits,
depending on the digit count n:

```
n       kept              shown/n
1–6     none              0
7–9     last 3            ≤ 43%
10–11   first 3 + last 2  ≤ 50%
≥ 12    first 3 + last 3  ≤ 50%
```

The output has one '*' per hidden digit, so the digit count is revealed.

```go
MaskPhone("+998 90 123-45-67") // "+998******567"
MaskPhone("(555) 123-4567")    // "555*****67"
MaskPhone("1234567")           // "****567"
MaskPhone("12345")             // "*****"
```

## MaskCard

```go
func MaskCard(pan string) string
```

MaskCard masks a payment card number (PAN). Revealed: only the last 4
digits, and the digit count. Spaces and '-' are removed first. Masked digits
are grouped in 4s from the left, followed by the last 4 digits. Input with
any other character, or with fewer than 12 or more than 19 digits, returns
"****" with nothing revealed.

```go
MaskCard("8600 1234 1234 5678") // "**** **** **** 5678"
MaskCard("378282246310005")     // "**** **** *** 0005"
```

## MaskToken

```go
func MaskToken(token string) string
```

MaskToken masks a secret token or API key.

Revealed:
  - A recognised vendor prefix: the longest leading run of segments, each
```go
1–6 lowercase ASCII letters followed by '_' or '-', that is at most 12
bytes ("sk_live_", "pk_test_", "ghp_", "github_pat_", "xoxb-"). Anything
else, including a '_' inside a random base64url token, is secret.
```

  - The last 4 runes of the secret, only when the secret (the input after
```
the vendor prefix) has at least 12 runes, so at most a third is shown.
```

The hidden part is always a fixed "****". If the secret has fewer than 12
runes, or the input is not valid UTF-8, "****" alone is returned and the
prefix is not shown either.

```go
MaskToken("sk_live_abc123456789") // "sk_live_****6789"
MaskToken("Ab3_x9kLmnopqrstuv")   // "****stuv"
MaskToken("short")                // "****"
```

## MaskIP

```go
func MaskIP(ip string) string
```

MaskIP masks an IP address. Revealed: for IPv4 (and IPv4-mapped IPv6) the
first two octets; for IPv6 the first two groups. Zones are dropped. Invalid
input returns "***".

```go
MaskIP("192.168.1.42") // "192.168.*.*"
MaskIP("2001:db8::1")  // "2001:db8:*"
```

## Mask

```go
func Mask(s string, keepStart, keepEnd int) string
```

Mask replaces the middle runes of s with '*', one per rune, keeping
keepStart leading and keepEnd trailing runes. Negative counts are treated as
0. If keepStart+keepEnd would reveal the whole string, every rune is masked
(fail closed). The output always has the same rune count as s. Invalid UTF-8
bytes count as one rune each; kept ones become U+FFFD.

```go
Mask("1234567890", 2, 2) // "12******90"
Mask("secret", 3, 3)     // "******"
```

## ID

```go
func ID(n uint64) string
```

ID formats n in groups of three digits separated by '-', for reading long
numeric identifiers aloud.

```go
ID(987654321234567) // "987-654-321-234-567"
ID(1234)            // "1-234"
ID(0)               // "0"
```

## Hash

```go
Package hash provides interfaces for hash functions.

type Cloner interface{ ... }
type Hash interface{ ... }
type Hash32 interface{ ... }
type Hash64 interface{ ... }
type XOF interface{ ... }
```



## ShortUUID

```go
func ShortUUID(u string) string
```

ShortUUID shortens a UUID to its first and last 4 runes. The input is not
validated; strings of 11 runes or fewer are returned unchanged.

```go
ShortUUID("550e8400-e29b-41d4-a716-446655440000") // "550e...0000"
```

## Truncate

```go
func Truncate(s string, head, tail int) string
```

Truncate shortens s to head leading and tail trailing runes joined by "...".
Negative counts are treated as 0. If the result would not be shorter than s
(head+tail+3 >= rune count), s is returned unchanged.

Truncate is for readability, not secrecy: it reveals head+tail runes.

```go
Truncate("a8f91234abcd", 4, 2) // "a8f9...cd"
```

## Redact

```go
func Redact(text string) string
```

Redact finds and masks sensitive values inside free text such as log lines
and error messages:

  - key=value secrets (password, secret, token, api_key, access_token,
```
auth) and "Bearer" credentials: the value becomes "****"
```

  - JWTs and vendor API keys (sk_live_, pk_test_, ghp_, github_pat_, xoxb-,
```
AKIA...): MaskToken
```

  - email addresses: MaskEmail
  - card numbers of 12–19 digits that pass the Luhn check: MaskCard
  - phone numbers starting with "+": MaskPhone
  - IPv4 addresses: MaskIP

Detection is pattern based and best effort: secrets without a recognisable
shape are not found. Use the Mask functions on known fields.

```go
Redact("login john.doe@gmail.com from 192.168.1.42")
// "login j***@gmail.com from 192.168.*.*"
```
