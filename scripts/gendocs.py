"""Generate API reference pages in docs/ from GoDoc. Run: make docs"""
import os, subprocess

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

PAGES = [
    ("numbers", "Numbers", 3,
     "Compact integers and floats, spell scales, ordinals and counted nouns.",
     ["Number", "NumberFloat", "NumberWithPrecision", "NumberWithOptions", "NumberOptions",
      "NumberWords", "Ordinal", "Count", "CountPlural", "Plural"]),
    ("bytes", "Bytes & throughput", 4,
     "Sizes in binary, IEC and SI units, and data transfer rates.",
     ["Bytes", "BytesIEC", "BytesSI", "FileSize", "BytesRate", "Throughput", "ThroughputIEC", "Bandwidth"]),
    ("time", "Durations & time", 5,
     "Durations at different lengths, latency, relative times, dates and ranges.\n\n"
     "`RelativeTime`, `Date`, `Time` and `TimeRange` read `time.Now`. Each has a "
     "`...From(now, ...)` twin that takes an explicit `now` — use it in tests.",
     ["Duration", "DurationShort", "DurationLong", "DurationWithOptions", "DurationOptions",
      "DurationNatural", "DurationApprox", "Latency", "RelativeTime", "RelativeTimeFrom",
      "Date", "DateFrom", "Time", "TimeFrom", "TimeRange", "TimeRangeFrom"]),
    ("money", "Money", 6,
     "Amounts are `int64` **minor units** (cents for USD), never floats.\n\n"
     "| Minor digits | Currencies |\n|---|---|\n"
     "| 2 | most ISO 4217 codes: USD, EUR, GBP, CHF, RUB, KZT, CNY, … |\n"
     "| 3 | BHD, IQD, JOD, KWD, LYD, OMR, TND |\n| 4 | CLF, UYW |\n"
     "| 0 | JPY, KRW, VND, UZS, … and any unknown code |\n\n"
     "Unknown codes are accepted; pick the digits explicitly with `MoneyWithPrecision`.",
     ["Money", "MoneySymbol", "MoneyWithPrecision", "MoneyCompact", "MoneyAccounting", "MoneyChange"]),
    ("percent", "Percent, progress & rates", 7,
     "Ratios, percentage changes, progress bars and event rates.",
     ["Percent", "PercentWithPrecision", "PercentChange", "ErrUndefined", "Progress", "ProgressBar",
      "Rate", "RateWithLabel", "PerMinute", "RequestRate"]),
    ("text", "Text", 8,
     "Labels from identifiers, English lists and booleans.",
     ["Humanize", "Enum", "List", "ListWithOptions", "ListOptions", "Bool", "BoolLabel"]),
    ("masking", "Masking & IDs", 9,
     "Maskers hide sensitive values for display and logging. Each documents exactly what it "
     "reveals and **fails closed**: input it cannot interpret is fully replaced by a placeholder.\n\n"
     "| Function | Revealed | Invalid or too short |\n|---|---|---|\n"
     "| `MaskEmail` | domain; first rune of local part if it has ≥ 3 runes | `***@domain` or `***` |\n"
     "| `MaskPhone` | ≤ half the digits, ≤ 3 trailing; digit count | all `*` or `***` |\n"
     "| `MaskCard` | last 4 digits; digit count (12–19 digits only) | `****` |\n"
     "| `MaskToken` | known vendor prefix; last 4 runes if secret ≥ 12 runes | `****` |\n"
     "| `MaskIP` | first two IPv4 octets / IPv6 groups | `***` |\n"
     "| `Mask` | `keepStart` + `keepEnd` runes; rune count | all `*` |\n\n"
     "{: .warning }\n`Truncate`, `Hash` and `ShortUUID` are for readability, not secrecy — "
     "they deliberately reveal leading and trailing characters.",
     ["MaskEmail", "MaskPhone", "MaskCard", "MaskToken", "MaskIP", "Mask",
      "ID", "Hash", "ShortUUID", "Truncate"]),
]


def godoc(name):
    out = subprocess.run(["go", "doc", "-u=false", name], cwd=ROOT, capture_output=True, text=True, check=True).stdout
    lines = [l for l in out.rstrip("\n").split("\n") if not l.startswith("package ")]
    # Signature: lines up to first indented doc line (types/vars span several lines).
    sig, i = [], 0
    while i < len(lines) and not lines[i].startswith("    "):
        sig.append(lines[i]); i += 1
    body = [l[4:] if l.startswith("    ") else l for l in lines[i:]]
    return "\n".join(sig).strip(), body


def render_body(body):
    md, code = [], []

    def flush():
        if code:
            while code and not code[-1].strip():
                code.pop()
            md.append("```go" if any("(" in c for c in code) else "```")
            md.extend(code)
            md.append("```")
            md.append("")
            code.clear()

    for l in body:
        if l.startswith("    ") or (code and not l.strip()):
            code.append(l[4:] if l.startswith("    ") else "")
            continue
        flush()
        md.append(l)
    flush()
    return "\n".join(md).strip()


def main():
    for slug, title, order, intro, names in PAGES:
        parts = [f"---\ntitle: {title}\nlayout: default\nnav_order: {order}\n---\n",
                 f"# {title}\n{{: .no_toc }}\n", intro + "\n",
                 "<details open markdown=\"block\">\n<summary>On this page</summary>\n{: .text-delta }\n- TOC\n{:toc}\n</details>\n"]
        for n in names:
            sig, body = godoc(n)
            parts.append(f"## {n}\n\n```go\n{sig}\n```\n\n{render_body(body)}\n")
        with open(f"{ROOT}/docs/{slug}.md", "w") as f:
            f.write("\n".join(parts))
        print("wrote", slug, len(names))


main()
