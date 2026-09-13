package readable_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/bakhod1r/readable"
)

func TestRedact(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", ""},
		{"nothing secret here, order 42", "nothing secret here, order 42"},
		{"login john.doe@gmail.com from 192.168.1.42", "login j***@gmail.com from 192.168.*.*"},
		{"card 8600 1234 1234 5678 ok", "card 8600 1234 1234 5678 ok"}, // fails Luhn
		{"card 4111 1111 1111 1111 ok", "card **** **** **** 1111 ok"},
		{"pan=4111-1111-1111-1111", "pan=**** **** **** 1111"},
		{"order 123456789012345678", "order 123456789012345678"},
		{"id 0000000000000 00000000", "id 0000000000000 00000000"},
		{"call +998 90 123-45-67 now", "call +998******567 now"},
		{"GET /cb?token=abc123&x=1", "GET /cb?token=****&x=1"},
		{"db PASSWORD=hunter2 api_key=zzz", "db PASSWORD=**** api_key=****"},
		{"Authorization: Bearer abc.def-ghi==", "Authorization: Bearer ****"},
		{"key sk_live_abc123456789", "key sk_live_****6789"},
		{"gh ghp_ABCDEFGHIJKLMNOPQRSTuvwx", "gh ghp_****uvwx"},
		{"aws AKIAABCDEFGHIJKLMNOP", "aws ****MNOP"},
		{"jwt eyJhbGciOi.eyJzdWIiOi.c2lnbmF0dXJl", "jwt ****dXJl"},
		{"version 1.2.3.4", "version 1.2.*.*"},
		{"bad ip 999.1.1.1", "bad ip ***"},
		{"at 12:30:45", "at 12:30:45"},
	}
	for _, tt := range tests {
		if got := readable.Redact(tt.in); got != tt.want {
			t.Errorf("Redact(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func TestRedactAttr(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		ReplaceAttr: func(g []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return readable.RedactAttr(g, a)
		},
	}))
	logger.Info("signup from 10.0.0.7",
		"email", "john.doe@gmail.com",
		"Phone", "+998 90 123-45-67",
		"card", "8600 1234 1234 5678",
		"ip", "192.168.1.42",
		"user_password", "hunter2",
		"X-Auth-Token", 12345,
		"api-key", "k",
		"note", "mail jo@x.io",
		"count", 3,
	)
	want := `level=INFO msg="signup from 10.0.*.*" email=j***@gmail.com Phone=+998******567 card="**** **** **** 5678" ip=192.168.*.* user_password=**** X-Auth-Token=**** api-key=**** note="mail ***@x.io" count=3`
	if got := strings.TrimSpace(buf.String()); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

func ExampleRedact() {
	fmt.Println(readable.Redact("login john.doe@gmail.com from 192.168.1.42 password=hunter2"))
	// Output: login j***@gmail.com from 192.168.*.* password=****
}

func ExampleRedactAttr() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // drop the timestamp for a stable example
			}
			return readable.RedactAttr(groups, a)
		},
	}))
	logger.Info("login", "email", "john.doe@gmail.com", "token", "sk_live_abc123456789")
	// Output: level=INFO msg=login email=j***@gmail.com token=****
}
