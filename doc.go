// Package money provides a precise monetary value type for use across
// storage, arithmetic, and transport.
//
// The Money type wraps github.com/shopspring/decimal so that currency
// amounts keep their full precision in the database and in every
// calculation, instead of being rounded through float64. This matters for
// billing: 0.1 plus 0.2 must be 0.3, not 0.30000000000000004.
//
// Money integrates with the rest of a service out of the box:
//
//   - It implements sql.Scanner and driver.Valuer, so it can be stored in
//     DECIMAL, NUMERIC, and TEXT columns and read back losslessly. Gorm is
//     supported via a GormDataType hint.
//   - It marshals as a JSON number and accepts either a JSON number or a
//     decimal string on input, so existing float64-based API contracts
//     keep working without client changes.
//
// Construct values with New, FromInt, FromString (preferred when precision
// matters), or FromFloat when only a float64 is available. Combine them
// with Add, Sub, Mul, Div, MulFloat, and DivFloat, and compare with
// LessThan, GreaterThan, Equals, IsZero, IsPositive, and IsNegative.
//
// Use Money for any currency amount stored in the database or used in
// billing arithmetic. Do NOT use it for non-money floats such as CPU
// percent, memory, or exchange rates; those remain float64.
package money
