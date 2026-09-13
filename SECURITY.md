# Security Policy

## Supported versions

`readable` has not tagged a release yet. Until v1, only the latest commit on
`main` is supported; after v1, fixes land on the latest minor release.

## Reporting a vulnerability

**Do not open a public issue.**

Report privately through GitHub's
[private vulnerability reporting](https://github.com/bakhod1r/readable/security/advisories/new),
or by email to **bakhodiryashinmansur@gmail.com** with `[readable security]` in
the subject.

Please include:

- the affected version or commit,
- a minimal reproduction (the call and the input value),
- the impact you believe it has.

You can expect an acknowledgement within 72 hours and an assessment within
7 days. If the report is confirmed, a patch release and a GitHub Security
Advisory follow; you will be credited unless you ask otherwise.

## Scope

`readable` only formats values into strings, so the security-relevant surface
is narrow but real. In scope:

- **Masking leaks** — `Mask`, `MaskEmail`, `MaskPhone`, `MaskCard`,
  `MaskToken`, `MaskIP`, `Truncate`, `Hash` or `ShortUUID` revealing more of
  the input than documented, for any input (short, empty, malformed, invalid
  UTF-8, multi-byte).
- **Panics** — any input that makes an exported function panic, which can
  crash a logging or request path.
- **Denial of service** — an input that makes a function allocate or loop
  without bound.

Out of scope:

- Sensitive values a caller logs without passing them through a `Mask*`
  function.
- Masking that is weaker than a caller wants but matches the documented
  behaviour — open a feature request instead.
- Anything requiring an attacker who already controls the process.

## Handling sensitive values safely

- Mask before logging, not after: `log.Printf("card=%s", readable.MaskCard(pan))`.
- Masking is for display, not protection. A masked value is not encrypted or
  hashed — never store it as a substitute for the real secret.
