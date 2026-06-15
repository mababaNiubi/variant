package variant

import (
	"encoding/json"
	"math"
	"testing"
)

// ============================================================================
// looksLikeJSON tests
// ============================================================================

func TestLooksLikeJSON(t *testing.T) {
	if !looksLikeJSON(`{"a":1}`) {
		t.Error("JSON object should be detected")
	}
	if !looksLikeJSON(`[1,2,3]`) {
		t.Error("JSON array should be detected")
	}
	if !looksLikeJSON(`  {"a":1}`) {
		t.Error("JSON with leading whitespace should be detected")
	}
	if !looksLikeJSON("\t\n [1]") {
		t.Error("JSON with whitespace chars should be detected")
	}
	if looksLikeJSON(`hello`) {
		t.Error("plain string should not be detected as JSON")
	}
	if looksLikeJSON("") {
		t.Error("empty string should not be detected as JSON")
	}
	if looksLikeJSON("   ") {
		t.Error("whitespace-only string should not be detected as JSON")
	}
}

// ============================================================================
// decode / decodeSlice / decodeMap tests (used by New)
// ============================================================================

func TestDecodeSlice(t *testing.T) {
	v := New([]any{int64(1), float64(2.5), "hello"})
	if v.Type() != TypeList || v.Len() != 3 {
		t.Fatalf("expected list of 3, got type=%v len=%d", v.Type(), v.Len())
	}
	if i, _ := v.ListGet(0); i.AsString() != "1" {
		t.Errorf("index 0: expected 1")
	}
	if i, _ := v.ListGet(1); i.AsString() != "2.5" {
		t.Errorf("index 1: expected 2.5")
	}
	if i, _ := v.ListGet(2); i.AsString() != "hello" {
		t.Errorf("index 2: expected hello")
	}
}

func TestDecode_InvalidAndPointer(t *testing.T) {
	v := New(nil)
	if !v.IsEmpty() {
		t.Errorf("expected empty for nil input, got %v", v.Type())
	}

	type Item struct{ Name string }
	item := &Item{Name: "test"}
	v = New(item)
	if v.Type() != TypeMap {
		t.Errorf("expected TypeMap for pointer to struct, got %v", v.Type())
	}
	name, _ := v.MapGet("Name")
	if name.AsString() != "test" {
		t.Errorf("expected test, got %s", name.AsString())
	}
}

func TestDecode_Uint(t *testing.T) {
	v := New(uint(42))
	if u, _ := v.AsUInt64(); u != 42 {
		t.Errorf("expected uint64(42), got %d", u)
	}
}

func TestDecodeMap(t *testing.T) {
	v := New(map[string]any{"name": "Alice", "age": int64(30)})
	if v.Type() != TypeMap || v.Len() != 2 {
		t.Fatalf("expected map of 2, got type=%v len=%d", v.Type(), v.Len())
	}
	name, _ := v.MapGet("name")
	if name.AsString() != "Alice" {
		t.Errorf("expected Alice, got %s", name.AsString())
	}
	age, _ := v.MapGet("age")
	if a, _ := age.AsInt64(); a != 30 {
		t.Errorf("expected 30, got %d", a)
	}
}

func TestNew_DefaultDecode(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	s := TestStruct{Name: "Alice", Age: 30}
	v := New(s)
	if v.Type() != TypeMap {
		t.Fatalf("expected TypeMap, got %v", v.Type())
	}
	name, ok := v.MapGet("name")
	if !ok || name.AsString() != "Alice" {
		t.Errorf("expected name=Alice, got %v", name.AsInterface())
	}
	age, ok := v.MapGet("age")
	if !ok {
		t.Error("expected age key")
	} else if a, _ := age.AsInt64(); a != 30 {
		t.Errorf("expected age=30, got %d", a)
	}
}

// ============================================================================
// MarshalJSON — direct JSON writing tests (new code paths)
// ============================================================================

func TestMarshalJSON_Empty(t *testing.T) {
	data, err := NewEmpty().MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "null" {
		t.Errorf("expected null, got %s", data)
	}
}

func TestMarshalJSON_Bool(t *testing.T) {
	data, err := NewBool(true).MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "true" {
		t.Errorf("expected true, got %s", data)
	}

	data, err = NewBool(false).MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "false" {
		t.Errorf("expected false, got %s", data)
	}
}

func TestMarshalJSON_Int64(t *testing.T) {
	tests := []struct {
		val  int64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{math.MaxInt64, "9223372036854775807"},
		{math.MinInt64, "-9223372036854775808"},
	}
	for _, tt := range tests {
		data, err := NewInt64(tt.val).MarshalJSON()
		if err != nil {
			t.Fatalf("NewInt64(%d): %v", tt.val, err)
		}
		if string(data) != tt.want {
			t.Errorf("NewInt64(%d): want %s, got %s", tt.val, tt.want, data)
		}
	}
}

func TestMarshalJSON_UInt64(t *testing.T) {
	tests := []struct {
		val  uint64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{math.MaxInt64, "9223372036854775807"},
		{math.MaxInt64 + 1, "9223372036854775808"},
		{math.MaxUint64, "18446744073709551615"},
	}
	for _, tt := range tests {
		data, err := NewUInt64(tt.val).MarshalJSON()
		if err != nil {
			t.Fatalf("NewUInt64(%d): %v", tt.val, err)
		}
		if string(data) != tt.want {
			t.Errorf("NewUInt64(%d): want %s, got %s", tt.val, tt.want, data)
		}
	}
}

func TestMarshalJSON_Float64(t *testing.T) {
	// Normal floats
	tests := []struct {
		val  float64
		want string
	}{
		{0.0, "0"},
		{1.5, "1.5"},
		{-3.14, "-3.14"},
		{1e10, "1e+10"},
		{1.5e-3, "0.0015"},
	}
	for _, tt := range tests {
		data, err := NewFloat64(tt.val).MarshalJSON()
		if err != nil {
			t.Fatalf("NewFloat64(%v): %v", tt.val, err)
		}
		if string(data) != tt.want {
			t.Errorf("NewFloat64(%v): want %s, got %s", tt.val, tt.want, data)
		}
	}

	// NaN and Inf should error
	_, err := NewFloat64(math.NaN()).MarshalJSON()
	if err == nil {
		t.Error("NaN should produce error")
	}
	_, err = NewFloat64(math.Inf(1)).MarshalJSON()
	if err == nil {
		t.Error("+Inf should produce error")
	}
	_, err = NewFloat64(math.Inf(-1)).MarshalJSON()
	if err == nil {
		t.Error("-Inf should produce error")
	}
}

func TestMarshalJSON_String(t *testing.T) {
	tests := []struct {
		val  string
		want string
	}{
		{"", `""`},
		{"hello", `"hello"`},
		{`he"llo`, `"he\"llo"`},
		{`back\slash`, `"back\\slash"`},
		{"line\nbreak", `"line\nbreak"`},
		{"tab\there", `"tab\there"`},
		{"carriage\rreturn", `"carriage\rreturn"`},
		{"\x00\x01\x1f", "\"\\u0000\\u0001\\u001f\""},
		{"中文", `"中文"`},
		{"<html>", "\"\\u003chtml\\u003e\""},
		{"a & b", "\"a \\u0026 b\""},
	}
	for _, tt := range tests {
		data, err := NewString(tt.val).MarshalJSON()
		if err != nil {
			t.Fatalf("NewString(%q): %v", tt.val, err)
		}
		if string(data) != tt.want {
			t.Errorf("NewString(%q): want %s, got %s", tt.val, tt.want, data)
		}
	}

	// Round-trip for HTML-escaped and Unicode strings
	roundTripStrings := []string{
		"<script>alert('xss')</script>",
		"a & b && c",
		"line1\nline2\r\nline3\tindented",
		"中文テスト🎉",
		"back\\slash and \"quote\"",
	}
	for _, s := range roundTripStrings {
		data, err := NewString(s).MarshalJSON()
		if err != nil {
			t.Fatalf("marshal %q: %v", s, err)
		}
		// Parse back with stdlib to verify correctness
		var decoded string
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("stdlib cannot parse our JSON for %q: %v (got: %s)", s, err, data)
		}
		if decoded != s {
			t.Errorf("round-trip mismatch: want %q, got %q (JSON: %s)", s, decoded, data)
		}
	}
}

func TestMarshalJSON_List(t *testing.T) {
	// Empty list
	data, err := NewValueList(nil).MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "[]" {
		t.Errorf("empty list: want [], got %s", data)
	}

	// Single element
	v := NewValueList([]Variant{NewInt(1)})
	data, err = v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "[1]" {
		t.Errorf("single: want [1], got %s", data)
	}

	// Mixed types
	v = NewValueList([]Variant{NewInt(1), NewString("two"), NewBool(true), NewEmpty()})
	data, err = v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `[1,"two",true,null]` {
		t.Errorf("mixed: got %s", data)
	}
}

func TestMarshalJSON_Map(t *testing.T) {
	// Empty map
	data, err := NewValueMap(nil).MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "{}" {
		t.Errorf("empty map: want {}, got %s", data)
	}

	// Simple map
	v := NewValueMap(map[string]Variant{
		"name": NewString("test"),
		"val":  NewInt(42),
	})
	data, err = v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}
	// Verify it's valid JSON
	if !json.Valid(data) {
		t.Errorf("invalid JSON: %s", data)
	}
}

func TestMarshalJSON_Nested(t *testing.T) {
	v := NewValueMap(map[string]Variant{
		"items": NewValueList([]Variant{
			NewValueMap(map[string]Variant{
				"id":   NewInt(1),
				"name": NewString("a"),
			}),
			NewValueMap(map[string]Variant{
				"id":   NewInt(2),
				"name": NewString("b"),
			}),
		}),
		"total": NewInt(2),
	})
	data, err := v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Errorf("invalid nested JSON: %s", data)
	}
}

func TestMarshalJSON_StdlibCompatible(t *testing.T) {
	// Verify our Marshal output can be parsed by stdlib json.Unmarshal
	v := NewValueMap(map[string]Variant{
		"name":    NewString("test"),
		"age":     NewInt(42),
		"score":   NewFloat64(99.5),
		"enabled": NewBool(true),
		"items":   NewValueList([]Variant{NewInt(1), NewInt(2), NewInt(3)}),
		"extra":   NewEmpty(),
	})
	data, err := v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("stdlib can't parse our JSON: %v", err)
	}

	if m["name"] != "test" {
		t.Errorf("name: want test, got %v", m["name"])
	}
	if m["age"] != float64(42) {
		t.Errorf("age: want 42, got %v (%T)", m["age"], m["age"])
	}
	if m["score"] != 99.5 {
		t.Errorf("score: want 99.5, got %v", m["score"])
	}
	if m["enabled"] != true {
		t.Errorf("enabled: want true, got %v", m["enabled"])
	}
	if m["extra"] != nil {
		t.Errorf("extra: want nil, got %v", m["extra"])
	}
}

// ============================================================================
// UnmarshalJSON — token-based decoding tests (new code paths)
// ============================================================================

func TestUnmarshalJSON_Null(t *testing.T) {
	v, err := UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeEmpty {
		t.Errorf("expected TypeEmpty, got %v", v.Type())
	}
}

func TestUnmarshalJSON_Bool(t *testing.T) {
	v, err := UnmarshalJSON([]byte("true"))
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := v.AsBool(); !b {
		t.Error("expected true")
	}

	v, err = UnmarshalJSON([]byte("false"))
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := v.AsBool(); b {
		t.Error("expected false")
	}
}

func TestUnmarshalJSON_Int64(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"9223372036854775807", math.MaxInt64},
		{"-9223372036854775808", math.MinInt64},
	}
	for _, tt := range tests {
		v, err := UnmarshalJSON([]byte(tt.input))
		if err != nil {
			t.Fatalf("UnmarshalJSON(%q): %v", tt.input, err)
		}
		if v.Type() != TypeInt64 {
			t.Errorf("UnmarshalJSON(%q): want TypeInt64, got %v", tt.input, v.Type())
		}
		got, _ := v.AsInt64()
		if got != tt.want {
			t.Errorf("UnmarshalJSON(%q): want %d, got %d", tt.input, tt.want, got)
		}
	}
}

func TestUnmarshalJSON_UInt64(t *testing.T) {
	// Values that overflow int64 should become UInt64
	v, err := UnmarshalJSON([]byte("9223372036854775808")) // MaxInt64 + 1
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeUInt64 {
		t.Errorf("want TypeUInt64, got %v", v.Type())
	}
	if u, _ := v.AsUInt64(); u != 9223372036854775808 {
		t.Errorf("want 9223372036854775808, got %d", u)
	}

	v, err = UnmarshalJSON([]byte("18446744073709551615")) // MaxUint64
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeUInt64 {
		t.Errorf("want TypeUInt64, got %v", v.Type())
	}
	if u, _ := v.AsUInt64(); u != math.MaxUint64 {
		t.Errorf("want MaxUint64, got %d", u)
	}
}

func TestUnmarshalJSON_Float64(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"3.14", 3.14},
		{"1.5", 1.5},
		{"-0.5", -0.5},
		{"1e10", 1e10},
		{"1e-3", 0.001},
		{"1.5e2", 150},
	}
	for _, tt := range tests {
		v, err := UnmarshalJSON([]byte(tt.input))
		if err != nil {
			t.Fatalf("UnmarshalJSON(%q): %v", tt.input, err)
		}
		if v.Type() != TypeFloat64 {
			t.Errorf("UnmarshalJSON(%q): want TypeFloat64, got %v", tt.input, v.Type())
		}
		got, _ := v.AsFloat64()
		if got != tt.want {
			t.Errorf("UnmarshalJSON(%q): want %v, got %v", tt.input, tt.want, got)
		}
	}
}

func TestUnmarshalJSON_OverflowToFloat64(t *testing.T) {
	// Number too large for uint64 → falls back to float64
	v, err := UnmarshalJSON([]byte("99999999999999999999"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeFloat64 {
		t.Errorf("want TypeFloat64, got %v", v.Type())
	}
	f, _ := v.AsFloat64()
	if f != 1e20 {
		t.Errorf("want 1e20, got %v", f)
	}
}

func TestUnmarshalJSON_String(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`""`, ""},
		{`"hello"`, "hello"},
		{`"line\nbreak"`, "line\nbreak"},
		{`"tab\there"`, "tab\there"},
		{`"escaped\"quote"`, `escaped"quote`},
		{`"back\\slash"`, `back\slash`},
		{`"<div>"`, "<div>"},
		{`"中文"`, "中文"},
	}
	for _, tt := range tests {
		v, err := UnmarshalJSON([]byte(tt.input))
		if err != nil {
			t.Fatalf("UnmarshalJSON(%q): %v", tt.input, err)
		}
		if v.Type() != TypeString {
			t.Errorf("UnmarshalJSON(%q): want TypeString, got %v", tt.input, v.Type())
		}
		got := v.AsString()
		if got != tt.want {
			t.Errorf("UnmarshalJSON(%q): want %q, got %q", tt.input, tt.want, got)
		}
	}
}

func TestUnmarshalJSON_EmptyArray(t *testing.T) {
	v, err := UnmarshalJSON([]byte("[]"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeList {
		t.Errorf("want TypeList, got %v", v.Type())
	}
	if v.Len() != 0 {
		t.Errorf("want len 0, got %d", v.Len())
	}
}

func TestUnmarshalJSON_EmptyObject(t *testing.T) {
	v, err := UnmarshalJSON([]byte("{}"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeMap {
		t.Errorf("want TypeMap, got %v", v.Type())
	}
	if v.Len() != 0 {
		t.Errorf("want len 0, got %d", v.Len())
	}
}

func TestUnmarshalJSON_Array(t *testing.T) {
	v, err := UnmarshalJSON([]byte(`[1,"two",true,null,3.14]`))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeList || v.Len() != 5 {
		t.Fatalf("want list of 5, got type=%v len=%d", v.Type(), v.Len())
	}

	// [0] = 1 (int64)
	e0, _ := v.ListGet(0)
	if e0.Type() != TypeInt64 {
		t.Errorf("[0]: want TypeInt64, got %v", e0.Type())
	}
	if i, _ := e0.AsInt64(); i != 1 {
		t.Errorf("[0]: want 1, got %d", i)
	}

	// [1] = "two"
	e1, _ := v.ListGet(1)
	if e1.Type() != TypeString || e1.AsString() != "two" {
		t.Errorf("[1]: want two, got %s", e1.AsString())
	}

	// [2] = true
	e2, _ := v.ListGet(2)
	if b, _ := e2.AsBool(); !b {
		t.Error("[2]: want true")
	}

	// [3] = null
	e3, _ := v.ListGet(3)
	if !e3.IsEmpty() {
		t.Error("[3]: want empty")
	}

	// [4] = 3.14 (float64)
	e4, _ := v.ListGet(4)
	if e4.Type() != TypeFloat64 {
		t.Errorf("[4]: want TypeFloat64, got %v", e4.Type())
	}
	if f, _ := e4.AsFloat64(); f != 3.14 {
		t.Errorf("[4]: want 3.14, got %v", f)
	}
}

func TestUnmarshalJSON_Object(t *testing.T) {
	data := []byte(`{"a":1,"b":"hello","c":true,"d":null}`)
	v, err := UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeMap || v.Len() != 4 {
		t.Fatalf("want map of 4, got type=%v len=%d", v.Type(), v.Len())
	}

	av, ok := v.MapGet("a")
	if !ok {
		t.Error("expected key 'a'")
	}
	ai, err := av.AsInt64()
	if err != nil || ai != 1 {
		t.Errorf("expected a=1, got %d (err=%v)", ai, err)
	}

	bv, ok := v.MapGet("b")
	if !ok || bv.AsString() != "hello" {
		t.Error("expected b=hello")
	}

	cv, ok := v.MapGet("c")
	if !ok {
		t.Error("expected key 'c'")
	}
	cb, _ := cv.AsBool()
	if !cb {
		t.Error("expected c=true")
	}

	dv, ok := v.MapGet("d")
	if !ok || !dv.IsEmpty() {
		t.Error("expected d=null")
	}
}

func TestUnmarshalJSON_Nested(t *testing.T) {
	data := []byte(`{"items":[{"id":1,"name":"a"},{"id":2,"name":"b"}],"total":2}`)
	v, err := UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeMap {
		t.Fatalf("want TypeMap, got %v", v.Type())
	}

	items, ok := v.MapGet("items")
	if !ok || items.Type() != TypeList || items.Len() != 2 {
		t.Fatalf("expected items list of 2")
	}

	i0, _ := items.ListGet(0)
	if i0.Type() != TypeMap {
		t.Errorf("items[0]: want TypeMap, got %v", i0.Type())
	}
	id0, _ := i0.MapGet("id")
	if id, _ := id0.AsInt64(); id != 1 {
		t.Errorf("items[0].id: want 1, got %d", id)
	}

	total, ok := v.MapGet("total")
	if !ok {
		t.Error("expected total key")
	}
	if tval, _ := total.AsInt64(); tval != 2 {
		t.Errorf("total: want 2, got %d", tval)
	}
}

func TestUnmarshalJSON_NestedArray(t *testing.T) {
	v, err := UnmarshalJSON([]byte("[[1,2],[3,4]]"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeList || v.Len() != 2 {
		t.Fatalf("want list of 2, got type=%v len=%d", v.Type(), v.Len())
	}

	inner, _ := v.ListGet(0)
	if inner.Type() != TypeList || inner.Len() != 2 {
		t.Errorf("inner[0]: want list of 2")
	}
}

func TestUnmarshalJSON_NumberEdgeCases(t *testing.T) {
	// Negative zero → parsed as int 0
	v, err := UnmarshalJSON([]byte("-0"))
	if err != nil {
		t.Fatal(err)
	}
	if i, _ := v.AsInt64(); i != 0 {
		t.Errorf("-0: want 0, got %d", i)
	}

	// Zero
	v, err = UnmarshalJSON([]byte("0"))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeInt64 {
		t.Errorf("0: want TypeInt64, got %v", v.Type())
	}
	if i, _ := v.AsInt64(); i != 0 {
		t.Errorf("0: want 0, got %d", i)
	}
}

// ============================================================================
// UnmarshalJSON error cases
// ============================================================================

func TestUnmarshalJSON_Errors(t *testing.T) {
	errCases := []struct {
		name  string
		input string
	}{
		{"invalid", "invalid"},
		{"truncated array", "[1,"},
		{"truncated object", `{"a":`},
		{"unclosed array", "[1,2"},
		{"unclosed object", `{"a":1`},
		{"unclosed string", `"hello`},
		{"bad token", `[}`},
	}
	for _, tc := range errCases {
		_, err := UnmarshalJSON([]byte(tc.input))
		if err == nil {
			t.Errorf("%s: expected error, got nil", tc.name)
		}
	}
}

func TestUnmarshalJSON_MethodError(t *testing.T) {
	var v Variant
	err := v.UnmarshalJSON([]byte(`invalid json`))
	if err == nil {
		t.Error("expected error for invalid JSON in method")
	}
}

func TestUnmarshalJSON_Standalone(t *testing.T) {
	v, err := UnmarshalJSON([]byte(`[1,2,3]`))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeList || v.Len() != 3 {
		t.Errorf("expected list of 3, got type=%v len=%d", v.Type(), v.Len())
	}

	_, err = UnmarshalJSON([]byte(`invalid`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// ============================================================================
// Round-trip tests (marshal → unmarshal)
// ============================================================================

func TestRoundTrip_Scalars(t *testing.T) {
	originals := []Variant{
		NewEmpty(),
		NewBool(true),
		NewBool(false),
		NewInt64(0),
		NewInt64(42),
		NewInt64(-1),
		NewInt64(math.MaxInt64),
		NewUInt64(math.MaxUint64),
		NewFloat64(3.14),
		NewFloat64(-0.5),
		NewString("hello"),
		NewString(""),
	}

	for _, orig := range originals {
		data, err := orig.MarshalJSON()
		if err != nil {
			t.Fatalf("marshal %v: %v", orig.Type(), err)
		}

		var restored Variant
		err = restored.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("unmarshal %s: %v", data, err)
		}

		if !orig.IsEqual(restored) {
			t.Errorf("round-trip mismatch for type %v: orig=%v, restored=%v (JSON: %s)",
				orig.Type(), orig.AsInterface(), restored.AsInterface(), data)
		}
	}
}

func TestRoundTrip_List(t *testing.T) {
	orig := NewValueList([]Variant{
		NewInt(1), NewString("two"), NewBool(true), NewEmpty(),
	})
	data, err := orig.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var restored Variant
	err = restored.UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	if restored.Len() != 4 {
		t.Errorf("want len 4, got %d", restored.Len())
	}
	e0, _ := restored.ListGet(0)
	if i, _ := e0.AsInt64(); i != 1 {
		t.Errorf("[0]: want 1, got %d", i)
	}
}

func TestRoundTrip_Map(t *testing.T) {
	orig := NewValueMap(map[string]Variant{
		"a": NewInt(1),
		"b": NewString("hello"),
	})
	data, err := orig.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	v, err := UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeMap || v.Len() != 2 {
		t.Fatalf("want map of 2, got type=%v len=%d", v.Type(), v.Len())
	}
	if a, _ := v.MapGet("a"); a.AsString() != "1" {
		t.Errorf("a: want 1, got %s", a.AsString())
	}
}

func TestRoundTrip_Nested(t *testing.T) {
	orig := NewValueMap(map[string]Variant{
		"data": NewValueList([]Variant{
			NewValueMap(map[string]Variant{
				"id":   NewInt(1),
				"name": NewString("first"),
			}),
		}),
	})

	data, err := orig.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var result Variant
	err = result.UnmarshalJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	wrapper, _ := result.MapGet("data")
	items, _ := wrapper.ListGet(0)
	name, _ := items.MapGet("name")
	if name.AsString() != "first" {
		t.Errorf("expected first, got %s", name.AsString())
	}
}

func TestRoundTrip_StdlibJSON(t *testing.T) {
	// Variant → our JSON → stdlib unmarshal → stdlib marshal → our unmarshal
	orig := NewValueMap(map[string]Variant{
		"name":  NewString("test"),
		"count": NewInt(42),
	})

	// Step 1: Variant → our JSON
	ourJSON, err := orig.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	// Step 2: our JSON → stdlib unmarshal
	var stdObj map[string]any
	if err := json.Unmarshal(ourJSON, &stdObj); err != nil {
		t.Fatalf("stdlib unmarshal: %v", err)
	}

	// Step 3: stdlib marshal
	stdJSON, err := json.Marshal(stdObj)
	if err != nil {
		t.Fatal(err)
	}

	// Step 4: stdlib JSON → our unmarshal
	var result Variant
	err = result.UnmarshalJSON(stdJSON)
	if err != nil {
		t.Fatal(err)
	}

	name, _ := result.MapGet("name")
	if name.AsString() != "test" {
		t.Errorf("name: want test, got %s", name.AsString())
	}
	count, _ := result.MapGet("count")
	// stdlib produces float64 for all numbers, so count becomes TypeFloat64
	if f, _ := count.AsFloat64(); f != 42 {
		t.Errorf("count: want 42, got %v", count.AsInterface())
	}
}

// ============================================================================
// Method-style unmarshal (on *Variant)
// ============================================================================

func TestVariant_UnmarshalJSON(t *testing.T) {
	var v Variant
	err := v.UnmarshalJSON([]byte(`{"a":1,"b":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeMap {
		t.Fatalf("expected TypeMap, got %v", v.Type())
	}
	av, ok := v.MapGet("a")
	if !ok {
		t.Error("expected key 'a'")
	}
	ai, err := av.AsInt64()
	if err != nil || ai != 1 {
		t.Errorf("expected a=1, got %d (err=%v)", ai, err)
	}
}

// ============================================================================
// Existing basic tests kept for backward compatibility
// ============================================================================

func TestVariant_MarshalJSON(t *testing.T) {
	v := NewValueMap(map[string]Variant{
		"name": NewString("test"),
		"val":  NewInt(42),
	})
	data, err := v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}
}
