// Command readable formats and parses values from the command line.
//
//	readable bytes 1536                    # 1.5 KB
//	readable number 1234567                # 1.23M
//	readable duration 90m                  # 1h 30m
//	readable parse-bytes "1.5 GiB"         # 1610612736
//	readable parse-duration "2 days, 3h"   # 51h0m0s
//	readable roman 2026                    # MMXXVI
//	readable slug "Hello, World!"          # hello-world
//	tail -f app.log | readable redact      # masks emails, cards, tokens...
//
// With no value argument, each line of standard input is converted.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/bakhod1r/readable"
)

var exit = os.Exit

func main() {
	exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

var commands = map[string]func(string) (string, error){
	"bytes": func(s string) (string, error) {
		n, err := readable.ParseBytes(s)
		return readable.Bytes(n), err
	},
	"number": func(s string) (string, error) {
		n, err := readable.ParseNumber(s)
		return readable.Number(n), err
	},
	"duration": func(s string) (string, error) {
		d, err := readable.ParseDuration(s)
		return readable.Duration(d), err
	},
	"parse-bytes": func(s string) (string, error) {
		n, err := readable.ParseBytes(s)
		return strconv.FormatUint(n, 10), err
	},
	"parse-number": func(s string) (string, error) {
		n, err := readable.ParseNumber(s)
		return strconv.FormatInt(n, 10), err
	},
	"parse-duration": func(s string) (string, error) {
		d, err := readable.ParseDuration(s)
		return d.String(), err
	},
	"roman": func(s string) (string, error) {
		n, err := strconv.Atoi(s)
		if r := readable.Roman(n); err == nil && r != "" {
			return r, nil
		}
		return "", fmt.Errorf("roman: %q is not an integer in 1-3999", s)
	},
	"slug":     infallible(readable.Slug),
	"humanize": infallible(readable.Humanize),
	"redact":   infallible(readable.Redact),
}

func infallible(f func(string) string) func(string) (string, error) {
	return func(s string) (string, error) { return f(s), nil }
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 || commands[args[0]] == nil {
		names := make([]string, 0, len(commands))
		for name := range commands {
			names = append(names, name)
		}
		sort.Strings(names)
		fmt.Fprintf(stderr, "usage: readable <%s> [value...]\n", strings.Join(names, "|"))
		return 2
	}
	cmd, code := commands[args[0]], 0
	convert := func(s string) {
		out, err := cmd(s)
		if err != nil {
			fmt.Fprintln(stderr, err)
			code = 1
			return
		}
		fmt.Fprintln(stdout, out)
	}
	if len(args) > 1 {
		convert(strings.Join(args[1:], " "))
		return code
	}
	sc := bufio.NewScanner(stdin)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		convert(sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return code
}
