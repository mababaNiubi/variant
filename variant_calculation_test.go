package variant

import (
	"testing"
)

func TestReduce_SameType(t *testing.T) {
	// Int64
	r, err := NewInt(10).Reduce(NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 7 {
		t.Errorf("10-3 expected 7, got %d", i)
	}

	// Float64
	r, err = NewFloat64(10.5).Reduce(NewFloat64(3.2))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 7.3) {
		t.Errorf("10.5-3.2 expected 7.3, got %f", f)
	}

	// Bool
	r, err = NewBool(true).Reduce(NewBool(true))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 0 {
		t.Errorf("1-1 expected 0, got %d", i)
	}

	// Empty
	r, err = NewEmpty().Reduce(NewInt(5))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 5 {
		t.Errorf("empty-5 expected 5, got %d", i)
	}

	// Unsupported type
	_, err = NewString("a").Reduce(NewString("b"))
	if err == nil {
		t.Error("expected error for string Reduce")
	}
}

func TestReduce_MixedTypes(t *testing.T) {
	// Float64 - Int64
	r, err := NewFloat64(10.5).Reduce(NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 7.5) {
		t.Errorf("10.5-3 expected 7.5, got %f", f)
	}

	// Int64 - Float64
	r, err = NewInt(10).Reduce(NewFloat64(3.5))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 6.5) {
		t.Errorf("10-3.5 expected 6.5, got %f", f)
	}

	// UInt64 - Int64
	r, err = NewUInt64(10).Reduce(NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	if u, _ := r.AsUInt64(); u != 7 {
		t.Errorf("10-3 expected 7, got %d", u)
	}

	// UInt64 - Float64
	r, err = NewUInt64(10).Reduce(NewFloat64(3.5))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 6.5) {
		t.Errorf("10-3.5 expected 6.5 (via float), got %f", f)
	}
}

func TestIncrease_SameType(t *testing.T) {
	// Int64
	r, err := NewInt(10).Increase(NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 13 {
		t.Errorf("10+3 expected 13, got %d", i)
	}

	// Float64
	r, err = NewFloat64(2.5).Increase(NewFloat64(1.5))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 4.0) {
		t.Errorf("2.5+1.5 expected 4.0, got %f", f)
	}

	// String concat
	r, err = NewString("hello").Increase(NewString(" world"))
	if err != nil {
		t.Fatal(err)
	}
	if r.AsString() != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", r.AsString())
	}

	// List concat
	r, err = NewValueList([]Variant{NewInt(1)}).Increase(NewValueList([]Variant{NewInt(2)}))
	if err != nil {
		t.Fatal(err)
	}
	if r.Len() != 2 {
		t.Errorf("expected len 2, got %d", r.Len())
	}

	// Map merge
	r, err = NewValueMap(map[string]Variant{"a": NewInt(1)}).Increase(
		NewValueMap(map[string]Variant{"b": NewInt(2)}))
	if err != nil {
		t.Fatal(err)
	}
	if r.Len() != 2 {
		t.Errorf("expected len 2, got %d", r.Len())
	}
	// verify overwrite
	r, err = NewValueMap(map[string]Variant{"a": NewInt(1)}).Increase(
		NewValueMap(map[string]Variant{"a": NewInt(99)}))
	if err != nil {
		t.Fatal(err)
	}
	v, _ := r.MapGet("a")
	if i, _ := v.AsInt64(); i != 99 {
		t.Errorf("expected 99 (overwritten), got %d", i)
	}

	// Empty
	r, err = NewEmpty().Increase(NewInt(5))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 5 {
		t.Errorf("empty+5 expected 5, got %d", i)
	}
}

func TestIncrease_MixedTypes(t *testing.T) {
	// String + Int
	r, err := NewString("val:").Increase(NewInt(42))
	if err != nil {
		t.Fatal(err)
	}
	if r.AsString() != "val:42" {
		t.Errorf("expected 'val:42', got '%s'", r.AsString())
	}

	// List + single value
	r, err = NewValueList([]Variant{NewInt(1)}).Increase(NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	if r.Len() != 2 {
		t.Errorf("expected len 2, got %d", r.Len())
	}

	// Float64 + Int64
	r, err = NewFloat64(2.5).Increase(NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 3.5) {
		t.Errorf("2.5+1 expected 3.5, got %f", f)
	}

	// Unsupported
	_, err = NewBool(true).Increase(NewString("a"))
	if err == nil {
		t.Error("expected error")
	}
}

func TestMultiple_SameType(t *testing.T) {
	r, err := NewInt(10).Multiple(NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 30 {
		t.Errorf("10*3 expected 30, got %d", i)
	}

	r, err = NewFloat64(2.5).Multiple(NewFloat64(2.0))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 5.0) {
		t.Errorf("2.5*2 expected 5.0, got %f", f)
	}

	r, err = NewEmpty().Multiple(NewInt(5))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 5 {
		t.Errorf("empty*5 expected 5, got %d", i)
	}
}

func TestMultiple_MixedTypes(t *testing.T) {
	r, err := NewFloat64(2.5).Multiple(NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 5.0) {
		t.Errorf("2.5*2 expected 5.0, got %f", f)
	}
}

func TestDivide_SameType(t *testing.T) {
	r, err := NewInt(10).Divide(NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 5 {
		t.Errorf("10/2 expected 5, got %d", i)
	}

	r, err = NewFloat64(5.0).Divide(NewFloat64(2.0))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 2.5) {
		t.Errorf("5/2 expected 2.5, got %f", f)
	}
}

func TestDivide_MixedTypes(t *testing.T) {
	r, err := NewFloat64(10.0).Divide(NewInt(4))
	if err != nil {
		t.Fatal(err)
	}
	if f, _ := r.AsFloat64(); !IsFloat64Equal(f, 2.5) {
		t.Errorf("10/4 expected 2.5, got %f", f)
	}
}

func TestDivide_Empty(t *testing.T) {
	r, err := NewEmpty().Divide(NewInt(5))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := r.AsInt64(); i != 5 {
		t.Errorf("empty/5 expected 5, got %d", i)
	}
}

func TestDecimal(t *testing.T) {
	v := NewFloat64(3.14159265)
	v = v.Decimal(2)
	if s := v.AsString(); s != "3.14" {
		t.Errorf("expected '3.14', got '%s'", s)
	}

	// List of floats
	list := NewValueList([]Variant{NewFloat64(1.234), NewFloat64(5.678)})
	list = list.Decimal(1)
	v0, _ := list.ListGet(0)
	if s := v0.AsString(); s != "1.2" {
		t.Errorf("expected '1.2', got '%s'", s)
	}

	// Non-float unaffected
	iv := NewInt(42)
	iv = iv.Decimal(2)
	if i, _ := iv.AsInt64(); i != 42 {
		t.Errorf("int should be unchanged, got %d", i)
	}
}

func TestCompareNumberBySymbol_IntInt(t *testing.T) {
	v := NewInt(10)
	r := NewInt(5)

	ok, err := v.CompareNumberBySymbol(r, GreaterThanSymbol)
	if err != nil || !ok {
		t.Error("10 > 5 should be true")
	}
	ok, err = v.CompareNumberBySymbol(r, LessThanSymbol)
	if err != nil || ok {
		t.Error("10 < 5 should be false")
	}
	ok, err = v.CompareNumberBySymbol(r, GreaterEqualSymbol)
	if err != nil || !ok {
		t.Error("10 >= 5 should be true")
	}
	ok, err = v.CompareNumberBySymbol(r, LessEqualSymbol)
	if err != nil || ok {
		t.Error("10 <= 5 should be false")
	}
	ok, err = v.CompareNumberBySymbol(r, EqualSymbol)
	if err != nil || ok {
		t.Error("10 = 5 should be false")
	}
	ok, err = v.CompareNumberBySymbol(r, NotEqualSymbol)
	if err != nil || !ok {
		t.Error("10 != 5 should be true")
	}

	ok, err = v.CompareNumberBySymbol(NewInt(10), EqualSymbol)
	if err != nil || !ok {
		t.Error("10 = 10 should be true")
	}
}

func TestCompareNumberBySymbol_IntFloat(t *testing.T) {
	v := NewInt(10)
	r := NewFloat64(9.5)
	ok, err := v.CompareNumberBySymbol(r, GreaterThanSymbol)
	if err != nil || !ok {
		t.Error("10 > 9.5 should be true")
	}
}

func TestCompareNumberBySymbol_FloatString(t *testing.T) {
	v := NewFloat64(3.14)
	r := NewString("3.14")
	ok, err := v.CompareNumberBySymbol(r, EqualSymbol)
	if err != nil || !ok {
		t.Error("3.14 = '3.14' should be true")
	}
}

func TestCompareNumberBySymbol_InvalidSymbol(t *testing.T) {
	_, err := NewInt(1).CompareNumberBySymbol(NewInt(1), "%%")
	if err == nil {
		t.Error("expected error for invalid symbol")
	}
}

func TestComparable(t *testing.T) {
	if !NewInt(10).Comparable(NewInt(5)) {
		t.Error("10 should be comparable (> 5)")
	}
	if NewInt(5).Comparable(NewInt(10)) {
		t.Error("5 should not be comparable (> 10)")
	}
	if !NewFloat64(3.14).Comparable(NewFloat64(2.71)) {
		t.Error("3.14 should be comparable (> 2.71)")
	}

	// String comparison
	if !NewString("b").Comparable(NewString("a")) {
		t.Error("'b' > 'a' should be true")
	}
}

func TestCompareNumberBySymbol_UInt64Negative(t *testing.T) {
	// UInt64 max value comparing with negative Int64
	v := NewUInt64(100)
	r := NewInt(-1)
	ok, err := v.CompareNumberBySymbol(r, GreaterThanSymbol)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("uint 100 > int -1 should be true")
	}
}
