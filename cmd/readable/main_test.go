package main

import (
	"errors"
	"os"
	"strings"
	"testing"
	"testing/iotest"
)

func TestRun(t *testing.T) {
	tests := []struct {
		args        []string
		stdin       string
		out, errOut string
		code        int
	}{
		{[]string{"bytes", "1536"}, "", "1.5 KB\n", "", 0},
		{[]string{"number", "1234567"}, "", "1.23M\n", "", 0},
		{[]string{"duration", "90m"}, "", "1h 30m\n", "", 0},
		{[]string{"parse-bytes", "1.5", "GiB"}, "", "1610612736\n", "", 0},
		{[]string{"parse-number", "12.5K"}, "", "12500\n", "", 0},
		{[]string{"parse-duration", "2 days, 3h"}, "", "51h0m0s\n", "", 0},
		{[]string{"roman", "2026"}, "", "MMXXVI\n", "", 0},
		{[]string{"roman", "0"}, "", "", "roman: \"0\" is not an integer in 1-3999\n", 1},
		{[]string{"slug", "Hello, World!"}, "", "hello-world\n", "", 0},
		{[]string{"humanize", "created_at"}, "", "Created at\n", "", 0},
		{[]string{"redact"}, "mail john.doe@gmail.com\nok\n", "mail j***@gmail.com\nok\n", "", 0},
		{[]string{"bytes"}, "1024\nnope\n", "1 KB\n", "readable: ParseBytes \"nope\": invalid syntax\n", 1},
		{nil, "", "", "usage: readable <bytes|duration|humanize|number|parse-bytes|parse-duration|parse-number|redact|roman|slug> [value...]\n", 2},
		{[]string{"what"}, "", "", "", 2},
	}
	for _, tt := range tests {
		var out, errOut strings.Builder
		code := run(tt.args, strings.NewReader(tt.stdin), &out, &errOut)
		if code != tt.code || out.String() != tt.out || (tt.errOut != "" || code != 2) && errOut.String() != tt.errOut {
			t.Errorf("run(%q) = %d, %q, %q; want %d, %q, %q", tt.args, code, out.String(), errOut.String(), tt.code, tt.out, tt.errOut)
		}
	}
}

func TestRunReadError(t *testing.T) {
	var out, errOut strings.Builder
	if code := run([]string{"slug"}, iotest.ErrReader(errors.New("boom")), &out, &errOut); code != 1 || errOut.String() != "boom\n" {
		t.Errorf("code %d, stderr %q", code, errOut.String())
	}
}

func TestMain(m *testing.M) {
	got := -1
	exit = func(code int) { got = code }
	os.Args = []string{"readable"}
	main()
	if got != 2 {
		panic("main did not exit with usage code")
	}
	exit = os.Exit
	os.Exit(m.Run())
}
