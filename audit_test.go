package readable

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

// Errors are the most common way a secret reaches a log: a driver error that
// quotes the DSN. RedactAttr must not pass them through untouched.
func TestRedactAttrRedactsErrors(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: RedactAttr}))
	log.Error("connect", "err", errors.New("dial postgres://app:hunter2@db.example.com/prod: refused"))
	if strings.Contains(buf.String(), "hunter2") {
		t.Fatalf("error text leaked the password: %s", buf.String())
	}
}

type stringerSecret struct{}

func (stringerSecret) String() string { return "token for alice@example.com" }

func TestRedactAttrRedactsStringers(t *testing.T) {
	a := RedactAttr(nil, slog.Any("who", stringerSecret{}))
	if strings.Contains(a.Value.String(), "alice@example.com") {
		t.Fatalf("Stringer leaked an email: %s", a.Value.String())
	}
}

type maskedNonString struct {
	PAN   int64  `readable:"mask=card"`
	Token []byte `readable:"mask=token"`
}

// A mask tag on a type it cannot mask must fail closed, not print the value.
func TestMaskTagOnNonStringFailsClosed(t *testing.T) {
	out := Format(maskedNonString{PAN: 4111111111111111, Token: []byte("sk_live_abcdef")})
	for k, v := range out {
		if strings.Contains(v, "4111111111111111") || strings.Contains(v, "115") {
			t.Fatalf("%s printed in the clear: %s", k, v)
		}
	}
}

// "1,5" is a decimal comma in half the locales Format writes; reading it as
// the grouped number 15 is silently ten times wrong.
func TestParseRejectsAmbiguousComma(t *testing.T) {
	if v, err := ParseBytes("1,5 KB"); err == nil {
		t.Fatalf("ParseBytes(1,5 KB) = %d, want an error", v)
	}
	if v, err := ParseBytes("1,500 B"); err != nil || v != 1500 {
		t.Fatalf("ParseBytes(1,500 B) = %d, %v; want 1500", v, err)
	}
}
