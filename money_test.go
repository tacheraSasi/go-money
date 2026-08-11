package money

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

// ── Constructors ────────────────────────────────────────────────────────────

func TestNew(t *testing.T) {
	m := New(decimal.NewFromInt(42))
	if m.String() != "42" {
		t.Fatalf("New: got %q want %q", m.String(), "42")
	}
}

func TestFromInt(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{-7, "-7"},
		{1000000, "1000000"},
	}
	for _, c := range cases {
		got := FromInt(c.in).String()
		if got != c.want {
			t.Errorf("FromInt(%d): got %q want %q", c.in, got, c.want)
		}
	}
}

func TestFromFloat(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{9.99, "9.99"},
		{-1.5, "-1.5"},
		{100, "100"},
	}
	for _, c := range cases {
		got := FromFloat(c.in).String()
		if got != c.want {
			t.Errorf("FromFloat(%v): got %q want %q", c.in, got, c.want)
		}
	}
}

func TestFromString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"0", "0"},
		{"9.99", "9.99"},
		{"-12.345", "-12.345"},
		{"100", "100"},
	}
	for _, c := range cases {
		m, err := FromString(c.in)
		if err != nil {
			t.Fatalf("FromString(%q): unexpected error: %v", c.in, err)
		}
		if m.String() != c.want {
			t.Errorf("FromString(%q): got %q want %q", c.in, m.String(), c.want)
		}
	}

	if _, err := FromString("abc"); err == nil {
		t.Error("FromString(invalid): expected error, got nil")
	}
}

func TestZero(t *testing.T) {
	if !Zero().IsZero() {
		t.Fatalf("Zero: got %s, want 0", Zero().String())
	}
}

// TestZeroValue checks that an uninitialized Money behaves as 0.
func TestZeroValue(t *testing.T) {
	var m Money
	if !m.IsZero() {
		t.Fatalf("zero-value Money: got %s, want 0", m.String())
	}
	if m.String() != "0" {
		t.Errorf("zero-value Money String: got %q want %q", m.String(), "0")
	}
}

// ── Arithmetic ───────────────────────────────────────────────────────────────

func TestAdd(t *testing.T) {
	cases := []struct {
		a, b, want string
	}{
		{"0", "0", "0"},
		{"1", "2", "3"},
		{"0.1", "0.2", "0.3"}, // the float64 footgun this package exists to avoid
		{"9.99", "0.01", "10"},
		{"-5", "5", "0"},
	}
	for _, c := range cases {
		a, _ := FromString(c.a)
		b, _ := FromString(c.b)
		got := a.Add(b).String()
		if got != c.want {
			t.Errorf("Add(%s, %s): got %s want %s", c.a, c.b, got, c.want)
		}
	}
}

func TestSub(t *testing.T) {
	cases := []struct {
		a, b, want string
	}{
		{"10", "3", "7"},
		{"5", "5", "0"},
		{"5", "10", "-5"},
		{"1.00", "0.25", "0.75"},
	}
	for _, c := range cases {
		a, _ := FromString(c.a)
		b, _ := FromString(c.b)
		got := a.Sub(b).String()
		if got != c.want {
			t.Errorf("Sub(%s, %s): got %s want %s", c.a, c.b, got, c.want)
		}
	}
}

func TestMul(t *testing.T) {
	cases := []struct {
		a, b, want string
	}{
		{"2", "3", "6"},
		{"9.99", "100", "999"},
		{"0.1", "0.1", "0.01"},
		{"-2", "3", "-6"},
	}
	for _, c := range cases {
		a, _ := FromString(c.a)
		b, _ := FromString(c.b)
		got := a.Mul(b).String()
		if got != c.want {
			t.Errorf("Mul(%s, %s): got %s want %s", c.a, c.b, got, c.want)
		}
	}
}

func TestDiv(t *testing.T) {
	cases := []struct {
		a, b, want string
	}{
		{"6", "3", "2"},
		{"10", "4", "2.5"},
		{"1", "100", "0.01"},
		{"-9", "3", "-3"},
	}
	for _, c := range cases {
		a, _ := FromString(c.a)
		b, _ := FromString(c.b)
		got := a.Div(b).String()
		if got != c.want {
			t.Errorf("Div(%s, %s): got %s want %s", c.a, c.b, got, c.want)
		}
	}
}

// TestDiv_ByZeroPanics documents that Div follows shopspring/decimal and
// panics on a zero divisor rather than returning an error.
func TestDiv_ByZeroPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Div by zero: expected panic, got none")
		}
	}()
	FromInt(1).Div(Zero())
}

func TestMulFloat(t *testing.T) {
	m, _ := FromString("9.99")
	if got := m.MulFloat(3).String(); got != "29.97" {
		t.Errorf("MulFloat(3): got %s want 29.97", got)
	}
	if got := m.MulFloat(-1).String(); got != "-9.99" {
		t.Errorf("MulFloat(-1): got %s want -9.99", got)
	}
	if got := m.MulFloat(0).String(); got != "0" {
		t.Errorf("MulFloat(0): got %s want 0", got)
	}
}

func TestDivFloat(t *testing.T) {
	m, _ := FromString("100")
	if got := m.DivFloat(4).String(); got != "25" {
		t.Errorf("DivFloat(4): got %s want 25", got)
	}
}

// ── Comparisons ─────────────────────────────────────────────────────────────

func TestComparisons(t *testing.T) {
	one, _ := FromString("1")
	two, _ := FromString("2")
	oneEq, _ := FromString("1")

	cases := []struct {
		name string
		got  bool
	}{
		{"1 < 2", one.LessThan(two)},
		{"1 <= 2", one.LessThanOrEqual(two)},
		{"1 <= 1", one.LessThanOrEqual(oneEq)},
		{"!(2 < 1)", !two.LessThan(one)},
		{"2 > 1", two.GreaterThan(one)},
		{"2 >= 1", two.GreaterThanOrEqual(one)},
		{"2 >= 2", two.GreaterThanOrEqual(two)},
		{"1 == 1", one.Equals(oneEq)},
		{"!(1 == 2)", !one.Equals(two)},
	}
	for _, c := range cases {
		if !c.got {
			t.Errorf("comparison %s: expected true", c.name)
		}
	}
}

func TestSignChecks(t *testing.T) {
	pos, _ := FromString("5")
	neg, _ := FromString("-5")
	z := Zero()

	if !pos.IsPositive() {
		t.Error("5: IsPositive expected true")
	}
	if pos.IsNegative() {
		t.Error("5: IsNegative expected false")
	}
	if !neg.IsNegative() {
		t.Error("-5: IsNegative expected true")
	}
	if neg.IsPositive() {
		t.Error("-5: IsPositive expected false")
	}
	if !z.IsZero() {
		t.Error("zero: IsZero expected true")
	}
	if z.IsPositive() {
		t.Error("zero: IsPositive expected false")
	}
	if z.IsNegative() {
		t.Error("zero: IsNegative expected false")
	}
}

// ── JSON ───────────────────────────────────────────────────────────────────

func TestMarshalJSON(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"0", "0"},
		{"10", "10"},
		{"9.99", "9.99"},
		{"-25.5", "-25.5"},
	}
	for _, c := range cases {
		m, _ := FromString(c.in)
		out, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("MarshalJSON(%s): %v", c.in, err)
		}
		if string(out) != c.want {
			t.Errorf("MarshalJSON(%s): got %s want %s", c.in, out, c.want)
		}
		// MarshalJSON must emit a JSON number, not a quoted string.
		if strings.HasPrefix(string(out), "\"") {
			t.Errorf("MarshalJSON(%s): output %s is a string, want a number", c.in, out)
		}
	}
}

// TestMarshalJSON_LossyBoundary documents the intentional precision boundary:
// MarshalJSON emits a float64 number, so high-precision decimals lose precision
// when crossing JSON. Storage and arithmetic remain exact; this only affects
// the JSON transport.
func TestMarshalJSON_LossyBoundary(t *testing.T) {
	m, _ := FromString("1.123456789123456789")
	out, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "1.1234567891234568" {
		t.Fatalf("lossy marshal: got %s want 1.1234567891234568", out)
	}

	// Round-tripping through JSON does not recover the original precision.
	var round Money
	if err := json.Unmarshal(out, &round); err != nil {
		t.Fatal(err)
	}
	if round.Equals(m) {
		t.Fatalf("lossy boundary: round-trip %s equals original %s, expected precision loss", round, m)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"number", "9.99", "9.99"},
		{"string", "\"9.99\"", "9.99"},
		{"integer", "100", "100"},
		{"negative", "-1.5", "-1.5"},
		{"null", "null", "0"},
	}
	for _, c := range cases {
		var m Money
		if err := json.Unmarshal([]byte(c.in), &m); err != nil {
			t.Fatalf("%s: UnmarshalJSON(%s): %v", c.name, c.in, err)
		}
		if m.String() != c.want {
			t.Errorf("%s: UnmarshalJSON(%s): got %s want %s", c.name, c.in, m.String(), c.want)
		}
	}

	// Invalid decimal string.
	if err := json.Unmarshal([]byte(`"abc"`), new(Money)); err == nil {
		t.Error("UnmarshalJSON(invalid string): expected error, got nil")
	}
	// Invalid JSON token: encoding/json rejects this before the hook runs.
	if err := json.Unmarshal([]byte(`abc`), new(Money)); err == nil {
		t.Error("UnmarshalJSON(invalid token): expected error, got nil")
	}
	// Empty quoted string is not a valid decimal and must error.
	if err := json.Unmarshal([]byte(`""`), new(Money)); err == nil {
		t.Error("UnmarshalJSON(empty quoted string): expected error, got nil")
	}
}

// TestUnmarshalJSON_DirectErrors pins the internal error branches that are
// unreachable via json.Unmarshal, because encoding/json validates the token
// boundary before delegating to the hook. Calling UnmarshalJSON directly
// exercises the hook's own error returns for malformed quoted strings and
// non-numeric unquoted tokens.
func TestUnmarshalJSON_DirectErrors(t *testing.T) {
	cases := [][]byte{
		[]byte(`"9.99`), // unterminated quoted string: inner json.Unmarshal fails
		[]byte(`abc`),   // non-numeric unquoted token: decimal.NewFromString fails
		[]byte(`""`),    // empty quoted string is not a valid decimal
	}
	for _, in := range cases {
		var m Money
		if err := m.UnmarshalJSON(in); err == nil {
			t.Errorf("UnmarshalJSON(%s): expected error, got nil", in)
		}
		if !m.IsZero() {
			t.Errorf("UnmarshalJSON(%s): on error got %s, want unchanged 0", in, m.String())
		}
	}
}

// TestUnmarshalJSON_EmptyRaw exercises the defensive empty-bytes branch by
// invoking the method directly, since encoding/json never delivers empty
// input to an UnmarshalJSON hook.
func TestUnmarshalJSON_EmptyRaw(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte("")); err != nil {
		t.Fatalf("UnmarshalJSON(empty): unexpected error %v", err)
	}
	if !m.IsZero() {
		t.Fatalf("UnmarshalJSON(empty): got %s, want 0", m.String())
	}
}

func TestJSON_RoundTrip(t *testing.T) {
	cases := []string{"0", "9.99", "10.5", "-25.5", "100"}
	for _, in := range cases {
		orig, _ := FromString(in)
		out, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("marshal %s: %v", in, err)
		}
		var round Money
		if err := json.Unmarshal(out, &round); err != nil {
			t.Fatalf("unmarshal %s: %v", in, err)
		}
		if !round.Equals(orig) {
			t.Errorf("round-trip %s: got %s, want %s", in, round.String(), orig.String())
		}
	}
}

// TestJSON_StructField covers the realistic API path: a struct field of
// type Money unmarshaling from a decimal string and marshaling back as a
// JSON number.
func TestJSON_StructField(t *testing.T) {
	type Wallet struct {
		Balance Money `json:"balance"`
	}

	var w Wallet
	if err := json.Unmarshal([]byte(`{"balance":"99.95"}`), &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if w.Balance.String() != "99.95" {
		t.Fatalf("balance: got %s want 99.95", w.Balance.String())
	}

	out, err := json.Marshal(w)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != `{"balance":99.95}` {
		t.Errorf("marshal: got %s want {\"balance\":99.95}", out)
	}
}

// ── SQL ─────────────────────────────────────────────────────────────────────

func TestScan(t *testing.T) {
	cases := []struct {
		name string
		src  any
		want string
	}{
		{"nil", nil, "0"},
		{"string", "9.99", "9.99"},
		{"bytes", []byte("9.99"), "9.99"},
		{"float64", float64(12.5), "12.5"},
		{"int64", int64(7), "7"},
		{"int64 zero", int64(0), "0"},
		{"quoted string", "\"9.99\"", "9.99"},
	}
	for _, c := range cases {
		var m Money
		if err := m.Scan(c.src); err != nil {
			t.Fatalf("%s: Scan(%v): %v", c.name, c.src, err)
		}
		if m.String() != c.want {
			t.Errorf("%s: Scan(%v): got %s want %s", c.name, c.src, m.String(), c.want)
		}
	}
}

func TestScan_Invalid(t *testing.T) {
	var m Money
	if err := m.Scan(true); err == nil {
		t.Fatal("Scan(bool): expected error, got nil")
	}
}

func TestValue(t *testing.T) {
	m, _ := FromString("9.99")
	v, err := m.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	// driver.Valuer returns the canonical decimal string for SQL drivers.
	s, ok := v.(string)
	if !ok {
		t.Fatalf("Value: expected string, got %T", v)
	}
	if s != "9.99" {
		t.Errorf("Value: got %q want %q", s, "9.99")
	}
}

func TestValue_Zero(t *testing.T) {
	v, err := Zero().Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("Value: expected string, got %T", v)
	}
	if s != "0" {
		t.Errorf("Value: got %q want %q", s, "0")
	}
}

func TestGormDataType(t *testing.T) {
	var m Money
	if got := m.GormDataType(); got != "decimal" {
		t.Errorf("GormDataType: got %q want %q", got, "decimal")
	}
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func TestString(t *testing.T) {
	m, _ := FromString("12.345")
	if got := m.String(); got != "12.345" {
		t.Errorf("String: got %q want %q", got, "12.345")
	}
}

func TestFloat(t *testing.T) {
	m, _ := FromString("12.5")
	if got := m.Float(); got != 12.5 {
		t.Errorf("Float: got %v want %v", got, 12.5)
	}
}
