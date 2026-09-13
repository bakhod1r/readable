# Security Policy

## Supported versions

| Version | Supported |
| ------- | --------- |
| 0.1.x   | Yes       |
| < 0.1   | No        |

Before v1, only the latest minor release receives fixes.

## Reporting a vulnerability

Please **do not** open a public issue. Report privately through
[GitHub Security Advisories](https://github.com/bakhod1r/readable/security/advisories/new)
("Report a vulnerability" on the Security tab).

Include the function, the exact input, the output you got, and why it is a
security problem.

## Scope

In scope:

- **Masking leaks.** `Mask`, `MaskCard`, `MaskEmail`, `MaskPhone`, `MaskIP`,
  `MaskToken` document what they reveal. Any input that makes them reveal more
  than documented (extra digits, length of hidden parts where a fixed
  placeholder is promised, secret bytes via invalid UTF-8) is a vulnerability.
- **Denial of service.** Panics, unbounded allocation, or super-linear time on
  pathological input (huge strings, extreme integers, NaN/Inf, invalid UTF-8).

Out of scope:

- `Truncate`, `Hash` and `ShortUUID` are for readability and deliberately
  reveal characters; they are not masking functions.
- Misuse such as logging the unmasked value next to the masked one.

## Response timeline

- Acknowledgement within 3 business days.
- Initial assessment within 7 days.
- Fix and advisory for confirmed issues within 30 days, coordinated with the
  reporter. Credit is given unless you ask otherwise.
