package readable

import (
	"net/netip"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Fixed placeholders. Their length never depends on the input, so they do
// not leak how much was hidden.
const (
	maskFull   = "***"  // unmaskable input; also the hidden part of an email local part
	maskStars4 = "****" // hidden part of a card number or token
)

// Payment card number (PAN) length bounds per ISO/IEC 7812-1: a PAN has at
// most 19 digits; 12 is the shortest length issued in practice.
const (
	panMinDigits = 12
	panMaxDigits = 19
	panLast      = 4 // trailing PAN digits allowed by PCI DSS truncation
)

// Email: the first rune of the local part is kept only when the local part
// has at least this many runes, so at most one of three or more is shown.
const emailMinLocalRunes = 3

// Token masking limits.
const (
	tokenMaxPrefix     = 12 // max bytes of recognised vendor prefix
	tokenMaxSegment    = 6  // max lowercase letters in one prefix segment
	tokenTail          = 4  // trailing runes revealed
	tokenMinSecretTail = 12 // secret runes required before the tail is revealed
)

// Mask replaces the middle runes of s with '*', one per rune, keeping
// keepStart leading and keepEnd trailing runes. Negative counts are treated
// as 0. If keepStart+keepEnd would reveal the whole string, every rune is
// masked (fail closed). The output always has the same rune count as s.
// Invalid UTF-8 bytes count as one rune each; kept ones become U+FFFD.
//
//	Mask("1234567890", 2, 2) // "12******90"
//	Mask("secret", 3, 3)     // "******"
func Mask(s string, keepStart, keepEnd int) string {
	keepStart, keepEnd = max(keepStart, 0), max(keepEnd, 0)
	n := utf8.RuneCountInString(s)
	if keepStart >= n || keepEnd >= n || keepStart+keepEnd >= n {
		return strings.Repeat("*", n)
	}
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for _, r := range s {
		if i < keepStart || i >= n-keepEnd {
			b.WriteRune(r)
		} else {
			b.WriteByte('*')
		}
		i++
	}
	return b.String()
}

// MaskEmail masks the local part of an email address and keeps the domain.
// The last '@' separates local part and domain.
//
// Revealed: the whole domain, plus the first rune of the local part only when
// the local part has 3 or more runes. Everything else in the local part is
// replaced by a fixed "***", so its length is not leaked. A local part of 1
// or 2 runes becomes "***@domain". Input without '@', with an empty local
// part or domain, or whose local part starts with invalid UTF-8, returns
// "***".
//
//	MaskEmail("john.doe@gmail.com") // "j***@gmail.com"
//	MaskEmail("jo@gmail.com")       // "***@gmail.com"
//	MaskEmail("not-an-email")       // "***"
func MaskEmail(email string) string {
	at := strings.LastIndexByte(email, '@')
	if at <= 0 || at == len(email)-1 {
		return maskFull
	}
	r, size := utf8.DecodeRuneInString(email)
	if r == utf8.RuneError && size <= 1 {
		return maskFull
	}
	if utf8.RuneCountInString(email[:at]) < emailMinLocalRunes {
		return maskFull + email[at:]
	}
	return email[:size] + maskFull + email[at:]
}

// MaskPhone masks a phone number digit by digit. Spaces, '-', '.', '(' and ')'
// are removed first; a single leading '+' is kept. Any other character makes
// the whole input unmaskable and "***" is returned.
//
// Revealed: at most half of the digits and never more than 3 trailing
// digits, depending on the digit count n:
//
//	n       kept              shown/n
//	1–6     none              0
//	7–9     last 3            ≤ 43%
//	10–11   first 3 + last 2  ≤ 50%
//	≥ 12    first 3 + last 3  ≤ 50%
//
// The output has one '*' per hidden digit, so the digit count is revealed.
//
//	MaskPhone("+998 90 123-45-67") // "+998******567"
//	MaskPhone("(555) 123-4567")    // "555*****67"
//	MaskPhone("1234567")           // "****567"
//	MaskPhone("12345")             // "*****"
func MaskPhone(phone string) string {
	plus := strings.HasPrefix(phone, "+")
	if plus {
		phone = phone[1:]
	}
	digits := make([]byte, 0, len(phone))
	for i := range len(phone) {
		c := phone[i]
		switch {
		case c >= '0' && c <= '9':
			digits = append(digits, c)
		case c == ' ' || c == '-' || c == '.' || c == '(' || c == ')':
		default:
			return maskFull
		}
	}
	n := len(digits)
	if n == 0 {
		return maskFull
	}
	head, tail := 0, 0
	switch {
	case n >= 12:
		head, tail = 3, 3
	case n >= 10:
		head, tail = 3, 2
	case n >= 7:
		tail = 3
	}
	for i := range n - head - tail {
		digits[head+i] = '*'
	}
	if plus {
		return "+" + string(digits)
	}
	return string(digits)
}

// MaskCard masks a payment card number (PAN). Revealed: only the last 4
// digits, and the digit count. Spaces and '-' are removed first. Masked digits are grouped in 4s
// from the left, followed by the last 4 digits. Input with any other
// character, or with fewer than 12 or more than 19 digits, returns "****"
// with nothing revealed.
//
//	MaskCard("8600 1234 1234 5678") // "**** **** **** 5678"
//	MaskCard("378282246310005")     // "**** **** *** 0005"
func MaskCard(pan string) string {
	var last [panLast]byte
	n := 0
	for i := range len(pan) {
		c := pan[i]
		switch {
		case c >= '0' && c <= '9':
			last[0], last[1], last[2], last[3] = last[1], last[2], last[3], c
			n++
		case c == ' ' || c == '-':
		default:
			return maskStars4
		}
	}
	if n < panMinDigits || n > panMaxDigits {
		return maskStars4
	}
	masked := n - panLast
	var b strings.Builder
	b.Grow(n + n/4 + 1)
	for i := range masked {
		if i > 0 && i%4 == 0 {
			b.WriteByte(' ')
		}
		b.WriteByte('*')
	}
	b.WriteByte(' ')
	b.Write(last[:])
	return b.String()
}

// MaskToken masks a secret token or API key.
//
// Revealed:
//   - A recognised vendor prefix: the longest leading run of segments, each
//     1–6 lowercase ASCII letters followed by '_' or '-', that is at most 12
//     bytes ("sk_live_", "pk_test_", "ghp_", "github_pat_", "xoxb-"). Anything
//     else, including a '_' inside a random base64url token, is secret.
//   - The last 4 runes of the secret, only when the secret (the input after
//     the vendor prefix) has at least 12 runes, so at most a third is shown.
//
// The hidden part is always a fixed "****". If the secret has fewer than 12
// runes, or the input is not valid UTF-8, "****" alone is returned and the
// prefix is not shown either.
//
//	MaskToken("sk_live_abc123456789") // "sk_live_****6789"
//	MaskToken("Ab3_x9kLmnopqrstuv")   // "****stuv"
//	MaskToken("short")                // "****"
func MaskToken(token string) string {
	if !utf8.ValidString(token) {
		return maskStars4
	}
	prefix := tokenVendorPrefix(token)
	rest := token[len(prefix):]
	if utf8.RuneCountInString(rest) < tokenMinSecretTail {
		return maskStars4
	}
	end := len(rest)
	for range tokenTail {
		_, size := utf8.DecodeLastRuneInString(rest[:end])
		end -= size
	}
	return prefix + maskStars4 + rest[end:]
}

// tokenVendorPrefix returns the longest prefix of s, at most tokenMaxPrefix
// bytes, made of segments of 1–tokenMaxSegment lowercase ASCII letters each
// followed by '_' or '-'. It returns "" if there is none.
func tokenVendorPrefix(s string) string {
	best, i := 0, 0
	for {
		j := i
		for j < len(s) && j-i < tokenMaxSegment && s[j] >= 'a' && s[j] <= 'z' {
			j++
		}
		if j == i || j >= len(s) || (s[j] != '_' && s[j] != '-') || j+1 > tokenMaxPrefix {
			return s[:best]
		}
		i = j + 1
		best = i
	}
}

// MaskIP masks an IP address. Revealed: for IPv4 (and IPv4-mapped IPv6) the
// first two octets; for IPv6 the first two groups. Zones are dropped. Invalid
// input returns "***".
//
//	MaskIP("192.168.1.42") // "192.168.*.*"
//	MaskIP("2001:db8::1")  // "2001:db8:*"
func MaskIP(ip string) string {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return maskFull
	}
	a = a.Unmap()
	var buf [40]byte
	out := buf[:0]
	if a.Is4() {
		b := a.As4()
		out = strconv.AppendUint(out, uint64(b[0]), 10)
		out = append(out, '.')
		out = strconv.AppendUint(out, uint64(b[1]), 10)
		out = append(out, ".*.*"...)
		return string(out)
	}
	b := a.As16()
	out = strconv.AppendUint(out, uint64(b[0])<<8|uint64(b[1]), 16)
	out = append(out, ':')
	out = strconv.AppendUint(out, uint64(b[2])<<8|uint64(b[3]), 16)
	out = append(out, ":*"...)
	return string(out)
}
