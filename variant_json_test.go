package variant

import "testing"

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
	// nil → reflect.Invalid
	v := New(nil)
	if !v.IsEmpty() {
		t.Errorf("expected empty for nil input, got %v", v.Type())
	}

	// pointer → reflect.Pointer (through default)
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

func TestUnmarshalJSON_MethodError(t *testing.T) {
	var v Variant
	err := v.UnmarshalJSON([]byte(`invalid json`))
	if err == nil {
		t.Error("expected error for invalid JSON in method")
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
