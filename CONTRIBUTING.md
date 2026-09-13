# Contributing to readable

Thanks for taking the time to contribute. This is a small, focused library, so
the bar is simple: keep it **zero-dependency**, keep it **pure**, and keep the
public API **stable**.

## Ground rules

1. **Stdlib only.** `readable` must never require a third-party module.
   `go.mod` has no `require` block, and CI fails if `go mod tidy` changes it.
   Test helpers are no exception.
2. **Pure and deterministic.** Formatters take a value and return a string. No
   global state, no hidden `time.Now()` (use the `...From(now, t)` variants),
   no locale lookups, no I/O. Every function must be safe for concurrent use.
3. **Never panic.** Any input — `MinInt64`, `NaN`, empty strings, invalid
   UTF-8 — returns a sensible string. New formatters get a fuzz target in
   `fuzz_test.go`.
4. **Backwards compatibility.** Output strings are part of the API. Changing
   what an existing call prints is a breaking change; new behaviour goes behind
   a new function or an options struct field, not a changed default.
5. **Masking never leaks.** `Mask*`, `Truncate`, `Hash` and `ShortUUID` must
   never reveal more characters than documented. If you touch `mask.go` or
   `id.go`, add a test proving short and malformed inputs stay hidden.

## Getting started

```bash
git clone https://github.com/bakhod1r/readable
cd readable
go test -race ./...
```

There is nothing to install — no code generation, no make targets.

## Before you open a pull request

Run what CI runs:

```bash
gofmt -l .              # must print nothing
go mod tidy && git diff --exit-code go.mod
go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@latest ./...
go test -race -cover ./...
```

For changes to a hot path, include before/after benchmark output in the pull
request description:

```bash
go test -run '^$' -bench . -benchmem -count=5 .
```

## Tests

Tests live next to the code they cover: `package readable` for internals,
`package readable_test` for the public API. Table-driven tests are the norm.
New exported functions get a runnable `Example` in one of the
`examples_*_test.go` files — it doubles as documentation, so keep its
`// Output:` accurate. A bug fix needs a test that fails before the fix and
passes after it.

## Commit messages

Conventional Commits, lowercase scope in parentheses:

```
fix(mask): hide all digits of short phone numbers
feat(money): add Compact option for large amounts
docs: document rounding rules
```

Types in use: `feat`, `fix`, `perf`, `docs`, `test`, `refactor`, `chore`.

## Reporting bugs

Open an issue with:

- the Go version and OS,
- the exact call and input value,
- what you expected and what you got.

For anything security-sensitive, follow [SECURITY.md](SECURITY.md) instead of
opening a public issue.

## Proposing features

Open an issue first and describe the value you need to format. Formatters that
work for any locale-free input are the easiest to land; features that need
locale data or external dependencies usually aren't accepted.

By contributing you agree that your work is licensed under the
[MIT License](LICENSE).
