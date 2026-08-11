# go-money

Preserve monetary precision in storage and arithmetic in Go.

`go-money` provides a `Money` type built on
[shopspring/decimal](https://github.com/shopspring/decimal) so currency
amounts keep their full precision in the database and in every calculation,
instead of being silently rounded through `float64`. It is small,
dependency-light, and ready to drop into any service that handles money.

## Why

Binary floating point can not represent most decimal fractions, so
arithmetic on `float64` amounts produces values such as
`0.30000000000000004` for `0.1 + 0.2`. That drift is unacceptable for
billing, wallets, invoicing, and reconciliation. `Money` keeps the value
as an exact decimal at every step and only converts to `float64` at
boundaries you control.

## Features

- Exact decimal arithmetic: `0.1 + 0.2` equals `0.3`.
- Constructors for `int64`, `float64`, and decimal strings, plus a zero value.
- Arithmetic: `Add`, `Sub`, `Mul`, `Div`, `MulFloat`, `DivFloat`.
- Comparisons: `LessThan`, `LessThanOrEqual`, `GreaterThan`,
  `GreaterThanOrEqual`, `Equals`, `IsZero`, `IsPositive`, `IsNegative`.
- Database ready: implements `sql.Scanner` and `driver.Valuer` for
  `DECIMAL` / `NUMERIC` / `TEXT` columns, with a Gorm data type hint.
- JSON friendly: marshals as a JSON number and accepts either a JSON
  number or a decimal string on input, so existing `float64` API
  contracts keep working without client changes.
- Reusable and self contained: no business logic, just money.

## Installation

```bash
go get github.com/tacheraSasi/go-money
```

Requires Go 1.26 or later.

## Usage

```go
package main

import (
	"fmt"

	money "github.com/tacheraSasi/go-money"
)

func main() {
	price, _ := money.FromString("9.99")
	quantity, _ := money.FromString("3")

	subtotal := price.Mul(quantity)
	discount := money.FromInt(2)
	total := subtotal.Sub(discount)

	fmt.Println(total.String())    // 27.97
	fmt.Println(total.IsPositive()) // true
	fmt.Println(total.Float())      // 27.97
}
```

Prefer `FromString`, `FromInt`, or `New` over `FromFloat` whenever
precision matters; `FromFloat` is a convenience for inputs that already
arrive as `float64`.

### JSON

`Money` marshals as a JSON number and accepts either a JSON number or a
decimal string on input, so existing `float64` API contracts keep working
without client changes.

```go
type Wallet struct {
	Balance money.Money `json:"balance"`
}

var w Wallet
_ = json.Unmarshal([]byte(`{"balance":"99.95"}`), &w)

out, _ := json.Marshal(w)
// out: {"balance":99.95}
```

### Database / Gorm

`Money` implements `sql.Scanner` and `driver.Valuer`, so it maps
directly to `DECIMAL`, `NUMERIC`, and `TEXT` columns. With Gorm:

```go
type Account struct {
	gorm.Model
	Credit money.Money `gorm:"type:decimal(20,2)"`
}
```

## Precision boundaries

`Money` is exact in storage and in arithmetic. The only lossy boundary is
JSON marshaling, which emits a `float64` number on purpose to preserve an
existing API contract. Values are still stored and computed exactly; the
rounding happens only when crossing JSON. If you need exact decimal
transport, marshal `m.String()` instead.

## API

Full reference at
[pkg.go.dev](https://pkg.go.dev/github.com/tacheraSasi/go-money).

| Category | Symbols |
| --- | --- |
| Constructors | `New`, `FromInt`, `FromFloat`, `FromString`, `Zero` |
| Arithmetic | `Add`, `Sub`, `Mul`, `Div`, `MulFloat`, `DivFloat` |
| Comparisons | `LessThan`, `LessThanOrEqual`, `GreaterThan`, `GreaterThanOrEqual`, `Equals`, `IsZero`, `IsPositive`, `IsNegative` |
| Storage | `Scan`, `Value`, `GormDataType` |
| Transport | `MarshalJSON`, `UnmarshalJSON` |
| Helpers | `String`, `Float` |

## License

MIT