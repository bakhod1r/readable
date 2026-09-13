# Contributing to readable

Thanks for helping. readable is small on purpose; please read this before
opening a pull request.

## Principles

- **Zero dependencies.** Standard library only. `go.mod` has no `require`.
- **Predictable output.** Same input, same string, on every platform and Go
  version. Output formats are documented in GoDoc with real examples and are
  treated as API.
- **No business logic.** Formatting only: no locale databases, no exchange
  rates, no timezone lookup, no parsing of user input into domain objects.
- **Pure functions.** No I/O, no hidden clock reads except in the documented
  `time.Now` convenience wrappers, which always have a `...From` variant.
- **No mutable globals.** Package-level tables are read-only; everything is
  safe for concurrent use without locks.
- **Integer math for exactness.** Integer inputs are scaled and rounded with
  integer arithmetic; floats are used only where the input is a float.
- **Fail closed** in masking functions: when in doubt, reveal less.

## Development workflow

Requirements: Go 1.24+, GNU make, [golangci-lint] v2.x.

```sh
make check        # gofmt check, go vet, golangci-lint, tests with -race
make cover        # coverage.html and the 100% gate
make fuzz         # every fuzz target, FUZZTIME=30s each
make bench        # benchmarks with -benchmem
make vuln         # govulncheck
make tidy
```

1. Open an issue first for new API so scope can be agreed.
2. Branch from `main`, keep the change focused.
3. `make check` and `make cover` must pass before you push.

## Test requirements

- 100% statement coverage is enforced in CI.
- Table-driven tests covering normal, boundary (zero, unit thresholds,
  rounding that promotes to the next unit), negative and extreme values
  (`math.MinInt64`, `math.MaxInt64`, `MaxUint64`, NaN, ±Inf, invalid UTF-8).
- A runnable `Example` for every exported function; its `// Output:` is the
  documentation.
- A benchmark for every new formatter; avoid regressions in allocations.
- A fuzz target, or an extension of one, for functions taking arbitrary
  strings or the full numeric range. Fuzz properties should check invariants
  (no panic, valid UTF-8, masking never reveals more than documented).
- Commit fuzz regression inputs under `testdata/fuzz`.

## API naming rules

- The base function takes the natural Go type and uses defaults:
  `Number(int64)`, `Bytes(uint64)`, `Duration(time.Duration)`.
- Variants use suffixes: `...WithPrecision`, `...WithOptions`, `...WithLabel`,
  `...From` (explicit `now`), unit families `...IEC` / `...SI`.
- Options go in a `...Options` struct whose zero value means the defaults.
- Return `string`; return `(string, error)` only when a result can be
  undefined, using the sentinel errors (`ErrUndefined`).
- Do not panic on any input.
- Names use the domain term exactly (`MaskCard`, not `HideCC`).

## Commit style

[Conventional Commits](https://www.conventionalcommits.org/):

```
feat(money): add MoneyCompact
fix(bytes): round 1023.999 KB up to 1 MB
docs: clarify Truncate is not for secrets
test(mask): fuzz MaskEmail for leaks
```

Use `feat!:` or a `BREAKING CHANGE:` footer for output or API changes.

## Versioning

[Semantic Versioning](https://semver.org/). Output strings are part of the API.

- **v0.x**: minor releases (`v0.2.0`) may contain breaking changes, always
  listed in `CHANGELOG.md`. Patch releases never break.
- **v1.0.0 onwards**: no breaking changes to signatures or documented output
  without a new major version. Bug fixes that change output which contradicted
  the documentation are allowed in minor releases and called out.

Every user-visible change gets a line under `[Unreleased]` in `CHANGELOG.md`.

[golangci-lint]: https://golangci-lint.run/
