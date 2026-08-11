// Package money provides a Money type that wraps shopspring/decimal to
// preserve monetary precision in storage and arithmetic, while serializing
// as a JSON number so the existing float64-based API contract stays intact.
//
// Use this for any currency amount stored in the database or used in
// billing arithmetic. Do NOT use it for non-money floats (CPU percent,
// memory, exchange rates, etc.) those stay float64.
package money