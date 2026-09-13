---
title: Getting started
layout: default
nav_order: 2
---

# Getting started

```sh
go get github.com/bakhod1r/readable
```

Requires **Go 1.24+**.

```go
package main

import (
	"fmt"
	"time"

	"github.com/bakhod1r/readable"
)

func main() {
	fmt.Println(readable.Bytes(1536))                  // 1.5 KB
	fmt.Println(readable.Duration(90 * time.Minute))   // 1h 30m
	fmt.Println(readable.MoneySymbol(150050, "USD"))   // $1,500.50
	fmt.Println(readable.MaskCard("8600123412345678")) // **** **** **** 5678
}
```

## Conventions

- The base function takes the natural Go type and uses defaults: `Number(int64)`, `Bytes(uint64)`, `Duration(time.Duration)`.
- Variants use suffixes: `...WithPrecision`, `...WithOptions`, `...WithLabel`, `...From` (explicit `now`), unit families `...IEC` / `...SI`.
- Options structs have a zero value that means the defaults.
- Functions return `string`; `(string, error)` only when the result can be undefined (`ErrUndefined`).

Full API reference: [pkg.go.dev](https://pkg.go.dev/github.com/bakhod1r/readable).
