package variant

import (
	"math"
	"testing"
)

func TestFloat64Boundaries(t *testing.T) {
	v := NewFloat64(math.MaxFloat64)
	f, err := v.AsFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f != math.MaxFloat64 {
		t.Errorf("expected MaxFloat64")
	}

	v2 := NewFloat64(math.SmallestNonzeroFloat64)
	f2, _ := v2.AsFloat64()
	if f2 != math.SmallestNonzeroFloat64 {
		t.Errorf("expected SmallestNonzeroFloat64")
	}
}

func TestMarshalBinary(t *testing.T) {
	tests := []Variant{
		NewEmpty(),
		NewBool(true),
		NewBool(false),
		NewInt64(42),
		NewInt64(-42),
		NewUInt64(18446744073709551615),
		NewFloat64(3.14159),
		NewFloat64(-1.5),
		NewString("hello"),
		NewString(""),
	}
	for _, original := range tests {
		data, err := original.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary(%s): %v", original.AsString(), err)
		}
		if !IsBinaryFormat(data) {
			t.Fatalf("expected binary format for %s", original.AsString())
		}
		decoded, n, err := UnmarshalBinary(data)
		if err != nil {
			t.Fatalf("UnmarshalBinary: %v", err)
		}
		if n != len(data) {
			t.Fatalf("consumed %d bytes, expected %d", n, len(data))
		}
		if !original.IsEqual(decoded) {
			t.Fatalf("roundtrip mismatch: type=%d vs type=%d, str=%s != %s",
				original.Type(), decoded.Type(), original.AsString(), decoded.AsString())
		}
	}
}

func TestMarshalBinaryNested(t *testing.T) {
	tests := []Variant{
		NewValueList([]Variant{NewInt64(1), NewFloat64(2.5), NewString("x")}),
		NewValueMap(map[string]Variant{"a": NewInt64(1), "b": NewString("hello")}),
		NewValueMap(map[string]Variant{
			"nested": NewValueList([]Variant{NewBool(true), NewEmpty()}),
		}),
	}
	for _, original := range tests {
		data, err := original.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary(%s): %v", original.AsString(), err)
		}
		decoded, _, err := UnmarshalBinary(data)
		if err != nil {
			t.Fatalf("UnmarshalBinary: %v", err)
		}
		if !original.IsEqual(decoded) {
			t.Fatalf("nested roundtrip mismatch: type=%d vs type=%d",
				original.Type(), decoded.Type())
		}
	}
}

// ─── Benchmarks ────────────────────────────────────────────────────────────────

func BenchmarkMarshalBinary_Float64(b *testing.B) {
	v := NewFloat64(3.141592653589793)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary_Float64(b *testing.B) {
	v := NewFloat64(3.141592653589793)
	data, _ := v.MarshalBinary()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = UnmarshalBinary(data)
	}
}

func BenchmarkMarshalBinary_Int64(b *testing.B) {
	v := NewInt64(9223372036854775807)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary_Int64(b *testing.B) {
	v := NewInt64(9223372036854775807)
	data, _ := v.MarshalBinary()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = UnmarshalBinary(data)
	}
}

func BenchmarkMarshalBinary_String(b *testing.B) {
	v := NewString("hello world")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary_String(b *testing.B) {
	v := NewString("hello world")
	data, _ := v.MarshalBinary()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = UnmarshalBinary(data)
	}
}

func BenchmarkAppendBinary_Batch(b *testing.B) {
	variants := make([]Variant, 1000)
	for i := range variants {
		variants[i] = NewFloat64(float64(i) * 1.5)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf []byte
		for j := range variants {
			buf = variants[j].AppendBinary(buf)
		}
	}
}

func BenchmarkMarshalBinary_NestedMap(b *testing.B) {
	v := NewValueMap(map[string]Variant{
		"name": NewString("test"),
		"val":  NewFloat64(3.14),
		"nested": NewValueList([]Variant{
			NewInt64(1),
			NewBool(true),
			NewString("x"),
		}),
	})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = v.MarshalBinary()
	}
}

func BenchmarkUnmarshalBinary_NestedMap(b *testing.B) {
	v := NewValueMap(map[string]Variant{
		"name": NewString("test"),
		"val":  NewFloat64(3.14),
		"nested": NewValueList([]Variant{
			NewInt64(1),
			NewBool(true),
			NewString("x"),
		}),
	})
	data, _ := v.MarshalBinary()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = UnmarshalBinary(data)
	}
}
