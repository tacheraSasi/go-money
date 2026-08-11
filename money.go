package money

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

// Money is a precise monetary amount. The zero value is 0.
type Money struct {
	decimal.Decimal
}

// New wraps a decimal.Decimal.
func New(d decimal.Decimal) Money { return Money{Decimal: d} }

// FromFloat creates a Money from a float64. Use for inputs only; prefer
// FromString or arithmetic on existing Money values when precision matters.
func FromFloat(f float64) Money { return Money{Decimal: decimal.NewFromFloat(f)} }

// FromInt creates a Money from an int64.
func FromInt(i int64) Money { return Money{Decimal: decimal.NewFromInt(i)} }

// FromString creates a Money from a decimal string (e.g. "9.99").
func FromString(s string) (Money, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, err
	}
	return Money{Decimal: d}, nil
}

// Zero returns a zero-valued Money.
func Zero() Money { return Money{Decimal: decimal.Zero} }

// Scan implements sql.Scanner so GORM can read DECIMAL/NUMERIC/TEXT columns.
func (m *Money) Scan(src any) error {
	if src == nil {
		m.Decimal = decimal.Zero
		return nil
	}
	return m.Decimal.Scan(src)
}

// Value implements driver.Valuer so GORM writes the precise decimal value.
func (m Money) Value() (driver.Value, error) {
	return m.Decimal.Value()
}

// GormDataType hints the underlying DB type for migrations.
func (m Money) GormDataType() string { return "decimal" }

// MarshalJSON emits Money as a JSON number to preserve the existing
// float64-based API contract. Precision is preserved in storage and backend
// arithmetic; the JSON boundary intentionally rounds to float64 since client
// magnitudes (wallet balances, plan prices) fit cleanly in a float64.
func (m Money) MarshalJSON() ([]byte, error) {
	f, _ := m.Decimal.Float64()
	return json.Marshal(f)
}

// UnmarshalJSON accepts both JSON numbers and decimal strings, so clients may
// send either form.
func (m *Money) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	if s == "null" || s == "" {
		m.Decimal = decimal.Zero
		return nil
	}
	if strings.HasPrefix(s, "\"") {
		var str string
		if err := json.Unmarshal(b, &str); err != nil {
			return err
		}
		d, err := decimal.NewFromString(str)
		if err != nil {
			return err
		}
		m.Decimal = d
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return errors.New("money: invalid amount: " + s)
	}
	m.Decimal = d
	return nil
}

// String returns the canonical decimal string.
func (m Money) String() string { return m.Decimal.String() }

// Float returns the value as a float64 (handy for printf with %f).
func (m Money) Float() float64 {
	f, _ := m.Decimal.Float64()
	return f
}

// ── Arithmetic ───────────────────────────────────────────────────────────

// Add returns the sum of m and o. Neither value is mutated.
func (m Money) Add(o Money) Money { return Money{Decimal: m.Decimal.Add(o.Decimal)} }

// Sub returns the difference m minus o.
func (m Money) Sub(o Money) Money { return Money{Decimal: m.Decimal.Sub(o.Decimal)} }

// Mul returns the product of m and o. Prefer MulFloat when scaling by a
// non-monetary factor such as a quantity or a tax rate.
func (m Money) Mul(o Money) Money { return Money{Decimal: m.Decimal.Mul(o.Decimal)} }

// Div returns the quotient m divided by o. Division by zero panics, the
// same as shopspring/decimal.
func (m Money) Div(o Money) Money { return Money{Decimal: m.Decimal.Div(o.Decimal)} }

// MulFloat scales a Money by a float64 factor (e.g. quantity, tax rate).
func (m Money) MulFloat(f float64) Money {
	return Money{Decimal: m.Decimal.Mul(decimal.NewFromFloat(f))}
}

// DivFloat divides a Money by a float64 (e.g. exchange rate).
func (m Money) DivFloat(f float64) Money {
	return Money{Decimal: m.Decimal.Div(decimal.NewFromFloat(f))}
}

// ── Comparisons ──────────────────────────────────────────────────────────

// LessThan reports whether m is strictly less than o.
func (m Money) LessThan(o Money) bool { return m.Decimal.LessThan(o.Decimal) }

// LessThanOrEqual reports whether m is less than or equal to o.
func (m Money) LessThanOrEqual(o Money) bool { return m.Decimal.LessThanOrEqual(o.Decimal) }

// GreaterThan reports whether m is strictly greater than o.
func (m Money) GreaterThan(o Money) bool { return m.Decimal.GreaterThan(o.Decimal) }

// GreaterThanOrEqual reports whether m is greater than or equal to o.
func (m Money) GreaterThanOrEqual(o Money) bool { return m.Decimal.GreaterThanOrEqual(o.Decimal) }

// Equals reports whether m and o represent the same amount.
func (m Money) Equals(o Money) bool { return m.Decimal.Equal(o.Decimal) }

// IsZero reports whether m is exactly zero.
func (m Money) IsZero() bool { return m.Decimal.IsZero() }

// IsPositive reports whether m is greater than zero.
func (m Money) IsPositive() bool { return m.Decimal.IsPositive() }

// IsNegative reports whether m is less than zero.
func (m Money) IsNegative() bool { return m.Decimal.IsNegative() }
