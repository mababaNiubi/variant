package variant

import (
	"math"
	"testing"
)

func TestIsFloat64Equal(t *testing.T) {
	if !IsFloat64Equal(0.1+0.2, 0.3) {
		t.Error("0.1+0.2 should equal 0.3 with tolerance")
	}
	if IsFloat64Equal(1.0, 2.0) {
		t.Error("1.0 should not equal 2.0")
	}
	if !IsFloat64Equal(1.0, 1.0) {
		t.Error("1.0 should equal 1.0")
	}
	// Test near-boundary tolerance
	if !IsFloat64Equal(1.0, 1.0+1e-15) {
		t.Error("1.0 should equal 1.0+1e-15 with tolerance")
	}
	if IsFloat64Equal(1.0, 1.0+1e-13) {
		t.Error("1.0 should not equal 1.0+1e-13")
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"123", true},
		{"-456", true},
		{"+789", true},
		{"3.14", true},
		{"-0.5", true},
		{"+1.0", true},
		{".5", true},
		{"1.", true},
		{"1e10", true},
		{"-1.5e-3", true},
		{"abc", false},
		{"", false},
		{"  123  ", true},
		{"1.2.3", false},
	}
	for _, tt := range tests {
		if got := IsNumber(tt.s); got != tt.want {
			t.Errorf("IsNumber(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestGetStringValueType(t *testing.T) {
	if GetStringValueType("") != TypeEmpty {
		t.Error("empty string should be TypeEmpty")
	}
	if GetStringValueType("123") != TypeInt64 {
		t.Error("'123' should be TypeInt64")
	}
	if GetStringValueType("3.14") != TypeFloat64 {
		t.Error("'3.14' should be TypeFloat64")
	}
	if GetStringValueType("1e5") != TypeFloat64 {
		t.Error("'1e5' should be TypeFloat64")
	}
	if GetStringValueType("hello") != TypeString {
		t.Error("'hello' should be TypeString")
	}
}

func TestIsSciNum(t *testing.T) {
	if !IsSciNum("1e5") {
		t.Error("'1e5' should be scientific number")
	}
	if !IsSciNum("1.5e-3") {
		t.Error("'1.5e-3' should be scientific number")
	}
	if !IsSciNum("1E10") {
		t.Error("'1E10' should be scientific number")
	}
	if IsSciNum("123") {
		t.Error("'123' should not be scientific number")
	}
	if IsSciNum("e5") {
		t.Error("'e5' should not be scientific number (empty before e)")
	}
	if IsSciNum("5e") {
		t.Error("'5e' should not be scientific number (empty after e)")
	}
}

func TestCompareNumberBySymbol_EdgeCases(t *testing.T) {
	// Int negative < UInt positive
	v := NewInt(-1)
	r := NewUInt64(0)
	ok, err := v.CompareNumberBySymbol(r, LessThanSymbol)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("-1 < 0 should be true")
	}

	// Float NaN-like behavior not directly testable but boundaries are
	ok, err = NewFloat64(math.MaxFloat64).CompareNumberBySymbol(NewFloat64(math.MaxFloat64), EqualSymbol)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("MaxFloat64 should equal itself")
	}
}

func TestIntScientificNotation(t *testing.T) {
	if !isInt("0") {
		t.Error("0 should be int")
	}
	if !isInt("-0") {
		t.Error("-0 should be int")
	}
	if isInt("3.14") {
		t.Error("3.14 should not be int")
	}
	if isInt("") {
		t.Error("empty should not be int")
	}
}

func TestDecimalScientificNotation(t *testing.T) {
	if !isDec("3.14") {
		t.Error("3.14 should be decimal")
	}
	if !isDec(".5") {
		t.Error(".5 should be decimal")
	}
	if !isDec("1.") {
		t.Error("1. should be decimal")
	}
	if isDec("123") {
		t.Error("123 should not be decimal")
	}
}
