package readable

import (
	"regexp"
	"strings"
)

// redactRule replaces every match of re with mask(match).
type redactRule struct {
	re   *regexp.Regexp
	mask func(string) string
}

var redactRules = []redactRule{
	// key=value secrets in URLs, query strings and config dumps.
	{regexp.MustCompile(`(?i)\b(?:password|passwd|pwd|secret|token|api[_-]?key|access[_-]?token|auth)=[^\s&"']+`),
		func(m string) string { return m[:strings.IndexByte(m, '=')+1] + maskStars4 }},
	{regexp.MustCompile(`(?i)\bbearer\s+[A-Za-z0-9._~+/-]+=*`),
		func(m string) string { return m[:len("bearer")] + " " + maskStars4 }},
	// JSON Web Tokens.
	{regexp.MustCompile(`\beyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]*`), MaskToken},
	// Vendor API keys.
	{regexp.MustCompile(`\b(?:(?:sk|pk|rk)_(?:live|test)_[A-Za-z0-9]{8,}|ghp_[A-Za-z0-9]{20,}|github_pat_[A-Za-z0-9_]{20,}|xox[abpr]-[A-Za-z0-9-]{10,}|AKIA[0-9A-Z]{16}\b)`), MaskToken},
	{regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`), MaskEmail},
	// Card numbers must pass the Luhn check; other long numbers are kept.
	{regexp.MustCompile(`\b\d(?:[ -]?\d){11,}\b`), func(m string) string {
		if n := digitCount(m); n <= panMaxDigits && luhn(m) {
			return MaskCard(m)
		}
		return m
	}},
	// International phone numbers.
	{regexp.MustCompile(`\+\d[\d ().-]{6,}\d`), MaskPhone},
	{regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`), MaskIP},
}

// Redact finds and masks sensitive values inside free text such as log
// lines and error messages:
//
//   - key=value secrets (password, secret, token, api_key, access_token,
//     auth) and "Bearer" credentials: the value becomes "****"
//   - JWTs and vendor API keys (sk_live_, pk_test_, ghp_, github_pat_,
//     xoxb-, AKIA...): MaskToken
//   - email addresses: MaskEmail
//   - card numbers of 12–19 digits that pass the Luhn check: MaskCard
//   - phone numbers starting with "+": MaskPhone
//   - IPv4 addresses: MaskIP
//
// Detection is pattern based and best effort: secrets without a recognisable
// shape are not found. Use the Mask functions on known fields.
//
//	Redact("login john.doe@gmail.com from 192.168.1.42")
//	// "login j***@gmail.com from 192.168.*.*"
func Redact(text string) string {
	for _, r := range redactRules {
		text = r.re.ReplaceAllStringFunc(text, r.mask)
	}
	return text
}

// luhn reports whether the digits of s (ignoring other bytes) pass the
// Luhn checksum.
func luhn(s string) bool {
	sum, double := 0, false
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			continue
		}
		d := int(c - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}

func digitCount(s string) int {
	n := 0
	for i := range len(s) {
		if s[i] >= '0' && s[i] <= '9' {
			n++
		}
	}
	return n
}
