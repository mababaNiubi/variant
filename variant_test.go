package variant

import (
	"math"
	"testing"
)

func TestNewEmpty(t *testing.T) {
	v := NewEmpty()
	if v.Type() != TypeEmpty {
		t.Errorf("expected TypeEmpty, got %v", v.Type())
	}
	if !v.IsEmpty() {
		t.Error("expected IsEmpty to be true")
	}
}

func TestNewBool(t *testing.T) {
	tv := NewBool(true)
	fv := NewBool(false)
	if b, _ := tv.AsBool(); !b {
		t.Error("expected true")
	}
	if b, _ := fv.AsBool(); b {
		t.Error("expected false")
	}
	if tv.Type() != TypeBool {
		t.Errorf("expected TypeBool, got %v", tv.Type())
	}
}

func TestNewInt(t *testing.T) {
	v := NewInt(42)
	if i, _ := v.AsInt64(); i != 42 {
		t.Errorf("expected 42, got %d", i)
	}
}

func TestNewInt64(t *testing.T) {
	v := NewInt64(9223372036854775807)
	if i, _ := v.AsInt64(); i != 9223372036854775807 {
		t.Errorf("expected 9223372036854775807, got %d", i)
	}
}

func TestNewUInt64(t *testing.T) {
	v := NewUInt64(18446744073709551615)
	if u, _ := v.AsUInt64(); u != 18446744073709551615 {
		t.Errorf("expected 18446744073709551615, got %d", u)
	}
}

func TestNewFloat64(t *testing.T) {
	v := NewFloat64(3.141592653589793)
	if f, _ := v.AsFloat64(); !IsFloat64Equal(f, 3.141592653589793) {
		t.Errorf("expected 3.141592653589793, got %f", f)
	}
}

func TestNewValue(t *testing.T) {
	v1 := NewValue(10.0)
	if v1.Type() != TypeInt64 {
		t.Errorf("expected TypeInt64 for whole float, got %v", v1.Type())
	}

	v2 := NewValue(3.14)
	if v2.Type() != TypeFloat64 {
		t.Errorf("expected TypeFloat64 for decimal float, got %v", v2.Type())
	}
}

func TestNewString(t *testing.T) {
	v := NewString("hello")
	if s := v.AsString(); s != "hello" {
		t.Errorf("expected 'hello', got '%s'", s)
	}
}

func TestNewValueList(t *testing.T) {
	v := NewValueList([]Variant{NewInt(1), NewInt(2)})
	if v.Type() != TypeList {
		t.Errorf("expected TypeList, got %v", v.Type())
	}
	if v.Len() != 2 {
		t.Errorf("expected len 2, got %d", v.Len())
	}
}

func TestNewValueMap(t *testing.T) {
	v := NewValueMap(map[string]Variant{"a": NewInt(1)})
	if v.Type() != TypeMap {
		t.Errorf("expected TypeMap, got %v", v.Type())
	}
	if v.Len() != 1 {
		t.Errorf("expected len 1, got %d", v.Len())
	}
}

func TestNew_StringTypes(t *testing.T) {
	v1 := New("123")
	if i, _ := v1.AsInt64(); i != 123 {
		t.Errorf("expected 123, got %d", i)
	}

	v2 := New("3.14")
	if f, _ := v2.AsFloat64(); !IsFloat64Equal(f, 3.14) {
		t.Errorf("expected 3.14, got %f", f)
	}

	v3 := New("")
	if !v3.IsEmpty() {
		t.Error("expected IsEmpty for empty string")
	}

	v4 := New("hello world")
	if s := v4.AsString(); s != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", s)
	}

	v5 := New(`{"key":"value"}`)
	if v5.Type() != TypeMap {
		t.Errorf("expected TypeMap for JSON object, got %v", v5.Type())
	}

	v6 := New(`[1,2,3]`)
	if v6.Type() != TypeList {
		t.Errorf("expected TypeList for JSON array, got %v", v6.Type())
	}
}

func TestNew_NumericTypes(t *testing.T) {
	v := New(int8(8))
	if i, _ := v.AsInt64(); i != 8 {
		t.Errorf("int8: expected 8, got %d", i)
	}
	v = New(int16(16))
	if i, _ := v.AsInt64(); i != 16 {
		t.Errorf("int16: expected 16, got %d", i)
	}
	v = New(int32(32))
	if i, _ := v.AsInt64(); i != 32 {
		t.Errorf("int32: expected 32, got %d", i)
	}
	v = New(int64(64))
	if i, _ := v.AsInt64(); i != 64 {
		t.Errorf("int64: expected 64, got %d", i)
	}
	v = New(uint(1))
	if u, _ := v.AsUInt64(); u != 1 {
		t.Errorf("uint: expected 1, got %d", u)
	}
	v = New(uint8(8))
	if u, _ := v.AsUInt64(); u != 8 {
		t.Errorf("uint8: expected 8, got %d", u)
	}
	v = New(uint16(16))
	if u, _ := v.AsUInt64(); u != 16 {
		t.Errorf("uint16: expected 16, got %d", u)
	}
	v = New(uint32(32))
	if u, _ := v.AsUInt64(); u != 32 {
		t.Errorf("uint32: expected 32, got %d", u)
	}
	v = New(uint64(64))
	if u, _ := v.AsUInt64(); u != 64 {
		t.Errorf("uint64: expected 64, got %d", u)
	}
}

func TestNew_Bool(t *testing.T) {
	v := New(true)
	if v.Type() != TypeBool {
		t.Errorf("expected TypeBool, got %v", v.Type())
	}
	if b, _ := v.AsBool(); !b {
		t.Error("expected true")
	}
}

func TestNew_VariantAndPointer(t *testing.T) {
	orig := NewInt(42)
	v := New(orig)
	if i, _ := v.AsInt64(); i != 42 {
		t.Errorf("expected 42, got %d", i)
	}

	vp := New(&orig)
	if i, _ := vp.AsInt64(); i != 42 {
		t.Errorf("expected 42, got %d", i)
	}
}

func TestNew_VariantSliceAndMap(t *testing.T) {
	vs := []Variant{NewInt(1), NewInt(2)}
	v1 := New(vs)
	if v1.Type() != TypeList || v1.Len() != 2 {
		t.Error("expected TypeList with 2 elements")
	}

	vm := map[string]Variant{"a": NewInt(1)}
	v2 := New(vm)
	if v2.Type() != TypeMap || v2.Len() != 1 {
		t.Error("expected TypeMap with 1 element")
	}
}

func TestNew_Bytes(t *testing.T) {
	v1 := New([]byte(`{"a":1}`))
	if v1.Type() != TypeMap {
		t.Errorf("expected TypeMap for JSON bytes, got %v", v1.Type())
	}

	// Non-JSON bytes are decoded as list of byte values via reflect
	v2 := New([]byte("plain text"))
	if v2.Type() != TypeList {
		t.Errorf("expected TypeList for plain bytes (decoded as []uint8), got %v", v2.Type())
	}
}

func TestIsEmpty(t *testing.T) {
	if !NewEmpty().IsEmpty() {
		t.Error("Empty should be empty")
	}
	if !NewString("").IsEmpty() {
		t.Error("empty string should be empty")
	}
	if NewString("a").IsEmpty() {
		t.Error("non-empty string should not be empty")
	}
	if !NewValueList([]Variant{NewEmpty(), NewEmpty()}).IsEmpty() {
		t.Error("list of empties should be empty")
	}
	if NewValueList([]Variant{NewEmpty(), NewInt(1)}).IsEmpty() {
		t.Error("list with non-empty element should not be empty")
	}
	if !NewValueMap(map[string]Variant{}).IsEmpty() {
		t.Error("empty map should be empty")
	}
	if NewInt(0).IsEmpty() {
		t.Error("int 0 should not be empty")
	}
}

func TestIsZero(t *testing.T) {
	if !NewInt(0).IsZero() {
		t.Error("int 0 should be zero")
	}
	if !NewFloat64(0.0).IsZero() {
		t.Error("float 0 should be zero")
	}
	if !NewString("0").IsZero() {
		t.Error(`string "0" should be zero`)
	}
	if NewInt(1).IsZero() {
		t.Error("int 1 should not be zero")
	}
	if NewString("hello").IsZero() {
		t.Error(`string "hello" should not be zero`)
	}
}

func TestVariantIsNumber(t *testing.T) {
	if !NewInt(1).IsNumber() {
		t.Error("int should be number")
	}
	if !NewFloat64(1.5).IsNumber() {
		t.Error("float should be number")
	}
	if !NewString("123").IsNumber() {
		t.Error(`string "123" should be number`)
	}
	if NewString("abc").IsNumber() {
		t.Error(`string "abc" should not be number`)
	}
}

func TestIsFloat(t *testing.T) {
	if !NewFloat64(1.5).IsFloat() {
		t.Error("float64 should be float")
	}
	if NewInt(1).IsFloat() {
		t.Error("int should not be float")
	}
	if !NewString("1.5").IsFloat() {
		t.Error(`string "1.5" should be float`)
	}
	if NewString("123").IsFloat() {
		t.Error(`string "123" should not be float`)
	}
}

func TestIsTrue(t *testing.T) {
	if !NewBool(true).IsTrue() {
		t.Error("true should be true")
	}
	if NewBool(false).IsTrue() {
		t.Error("false should not be true")
	}
	if !NewInt(1).IsTrue() {
		t.Error("int 1 should be true")
	}
	if NewInt(0).IsTrue() {
		t.Error("int 0 should not be true")
	}
}

func TestAsBool(t *testing.T) {
	tests := []struct {
		v    Variant
		want bool
		err  bool
	}{
		{NewEmpty(), false, false},
		{NewBool(true), true, false},
		{NewBool(false), false, false},
		{NewInt(1), true, false},
		{NewInt(0), false, false},
		{NewFloat64(1.5), true, false},
		{NewFloat64(0.0), false, false},
		{NewString("true"), true, false},
		{NewString("false"), false, false},
		{NewValueList([]Variant{}), false, true},
	}
	for _, tt := range tests {
		got, err := tt.v.AsBool()
		if (err != nil) != tt.err {
			t.Errorf("AsBool(%v) error = %v, wantErr %v", tt.v.Type(), err, tt.err)
		}
		if got != tt.want {
			t.Errorf("AsBool(%v) = %v, want %v", tt.v.Type(), got, tt.want)
		}
	}
}

func TestAsInt64(t *testing.T) {
	v := NewInt(42)
	if i, _ := v.AsInt64(); i != 42 {
		t.Errorf("expected 42, got %d", i)
	}

	fv := NewFloat64(3.14)
	if i, _ := fv.AsInt64(); i != 3 {
		t.Errorf("expected 3, got %d", i)
	}

	sv := NewString("100")
	if i, _ := sv.AsInt64(); i != 100 {
		t.Errorf("expected 100, got %d", i)
	}

	bigF := NewFloat64(1e20)
	if _, err := bigF.AsInt64(); err == nil {
		t.Error("expected overflow error")
	}

	lv := NewValueList([]Variant{})
	if _, err := lv.AsInt64(); err == nil {
		t.Error("expected type error")
	}

	bv := NewBool(true)
	if i, _ := bv.AsInt64(); i != 1 {
		t.Errorf("expected 1, got %d", i)
	}
}

func TestAsUInt64(t *testing.T) {
	v := NewUInt64(42)
	if u, _ := v.AsUInt64(); u != 42 {
		t.Errorf("expected 42, got %d", u)
	}

	nv := NewInt(-1)
	if _, err := nv.AsUInt64(); err == nil {
		t.Error("expected overflow for negative int")
	}

	bv := NewBool(true)
	if u, _ := bv.AsUInt64(); u != 1 {
		t.Errorf("expected 1, got %d", u)
	}
}

func TestAsFloat64(t *testing.T) {
	v := NewFloat64(3.14)
	if f, _ := v.AsFloat64(); !IsFloat64Equal(f, 3.14) {
		t.Errorf("expected 3.14, got %f", f)
	}

	iv := NewInt(42)
	if f, _ := iv.AsFloat64(); !IsFloat64Equal(f, 42.0) {
		t.Errorf("expected 42.0, got %f", f)
	}

	sv := NewString("3.14")
	if f, _ := sv.AsFloat64(); !IsFloat64Equal(f, 3.14) {
		t.Errorf("expected 3.14, got %f", f)
	}

	ev := NewEmpty()
	if f, _ := ev.AsFloat64(); f != 0 {
		t.Errorf("expected 0 for empty, got %f", f)
	}
}

func TestAsFloat32(t *testing.T) {
	v := NewFloat64(3.14)
	if f, _ := v.AsFloat32(); f < 3.139 || f > 3.141 {
		t.Errorf("expected ~3.14, got %f", f)
	}
}

func TestAsString(t *testing.T) {
	if s := NewString("hello").AsString(); s != "hello" {
		t.Errorf("expected 'hello', got '%s'", s)
	}
	if s := NewInt(42).AsString(); s != "42" {
		t.Errorf("expected '42', got '%s'", s)
	}
	if s := NewFloat64(3.14).AsString(); s != "3.14" {
		t.Errorf("expected '3.14', got '%s'", s)
	}

	list := NewValueList([]Variant{NewInt(1), NewString("two")})
	sl := list.AsString()
	if sl != `[1,"two"]` {
		t.Errorf("expected '[1,\"two\"]', got '%s'", sl)
	}

	mp := NewValueMap(map[string]Variant{"key": NewInt(1)})
	sm := mp.AsString()
	if sm != `{"key":1}` {
		t.Errorf("expected '{\"key\":1}', got '%s'", sm)
	}
}

func TestAsInterface(t *testing.T) {
	v1 := NewInt(42)
	if a := v1.AsInterface(); a.(int64) != 42 {
		t.Errorf("expected 42, got %v", a)
	}
	v2 := NewFloat64(3.14)
	if a := v2.AsInterface(); !IsFloat64Equal(a.(float64), 3.14) {
		t.Errorf("expected 3.14, got %v", a)
	}
	v3 := NewString("hello")
	if a := v3.AsInterface(); a.(string) != "hello" {
		t.Errorf("expected 'hello', got %v", a)
	}
	v4 := NewBool(true)
	if a := v4.AsInterface(); !a.(bool) {
		t.Errorf("expected true, got %v", a)
	}

	vlist := NewValueList([]Variant{NewInt(1), NewInt(2)})
	listIf := vlist.AsInterface().([]any)
	if len(listIf) != 2 {
		t.Errorf("expected len 2, got %d", len(listIf))
	}

	vmap := NewValueMap(map[string]Variant{"a": NewInt(1)})
	mapIf := vmap.AsInterface().(map[string]any)
	if mapIf["a"].(int64) != 1 {
		t.Errorf("expected 1, got %v", mapIf["a"])
	}
}

func TestAddString(t *testing.T) {
	v := NewEmpty()
	v.AddString("hello")
	if v.AsString() != "hello" {
		t.Errorf("expected 'hello', got '%s'", v.AsString())
	}
	v.AddString(" world")
	if v.AsString() != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", v.AsString())
	}
}

func TestIsEqual(t *testing.T) {
	if !NewInt(42).IsEqual(NewInt(42)) {
		t.Error("42 should equal 42")
	}
	if NewInt(42).IsEqual(NewInt(43)) {
		t.Error("42 should not equal 43")
	}
	if !NewFloat64(3.14).IsEqual(NewFloat64(3.14)) {
		t.Error("3.14 should equal 3.14")
	}
	if NewFloat64(3.14).IsEqual(NewFloat64(3.15)) {
		t.Error("3.14 should not equal 3.15")
	}
	if NewInt(42).IsEqual(NewFloat64(42.0)) {
		t.Error("different types should not be equal")
	}
	if !NewString("hello").IsEqual(NewString("hello")) {
		t.Error("strings should be equal")
	}

	list1 := NewValueList([]Variant{NewInt(1), NewInt(2)})
	list2 := NewValueList([]Variant{NewInt(1), NewInt(2)})
	list3 := NewValueList([]Variant{NewInt(1), NewInt(3)})
	if !list1.IsEqual(list2) {
		t.Error("identical lists should be equal")
	}
	if list1.IsEqual(list3) {
		t.Error("different lists should not be equal")
	}

	map1 := NewValueMap(map[string]Variant{"a": NewInt(1)})
	map2 := NewValueMap(map[string]Variant{"a": NewInt(1)})
	map3 := NewValueMap(map[string]Variant{"a": NewInt(2)})
	if !map1.IsEqual(map2) {
		t.Error("identical maps should be equal")
	}
	if map1.IsEqual(map3) {
		t.Error("different maps should not be equal")
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

func TestIsContainer(t *testing.T) {
	v1 := NewValueList([]Variant{})
	if v1.IsContainer() != true {
		t.Error("list should be container")
	}
	v2 := NewValueMap(map[string]Variant{})
	if v2.IsContainer() != true {
		t.Error("map should be container")
	}
	v3 := NewInt(1)
	if v3.IsContainer() != false {
		t.Error("int should not be container")
	}
}

func TestLen(t *testing.T) {
	s := NewString("hello")
	if s.Len() != 5 {
		t.Errorf("expected len 5, got %d", s.Len())
	}
	l := NewValueList([]Variant{NewInt(1), NewInt(2)})
	if l.Len() != 2 {
		t.Errorf("expected len 2, got %d", l.Len())
	}
	i := NewInt(42)
	if i.Len() != 0 {
		t.Errorf("expected len 0 for int, got %d", i.Len())
	}
}

func TestAdd(t *testing.T) {
	v := NewEmpty()
	_, err := v.Add(NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	if v.Type() != TypeList || v.Len() != 1 {
		t.Errorf("expected list of 1, got type=%v len=%d", v.Type(), v.Len())
	}

	_, err = v.Add(NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	if v.Len() != 2 {
		t.Errorf("expected len 2, got %d", v.Len())
	}

	sv := NewString("hello")
	_, err = sv.Add(" world")
	if err != nil {
		t.Fatal(err)
	}
	if sv.AsString() != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", sv.AsString())
	}

	iv := NewInt(1)
	_, err = iv.Add(2)
	if err == nil {
		t.Error("expected error adding to int")
	}
}

func TestRange(t *testing.T) {
	mp := NewValueMap(map[string]Variant{"a": NewInt(1), "b": NewInt(2)})
	count := 0
	mp.Range(func(key string, value Variant) bool {
		count++
		return true
	})
	if count != 2 {
		t.Errorf("expected 2 iterations, got %d", count)
	}

	count = 0
	mp.Range(func(key string, value Variant) bool {
		count++
		return false
	})
	if count != 1 {
		t.Errorf("expected 1 iteration (early stop), got %d", count)
	}

	list := NewValueList([]Variant{NewInt(1), NewInt(2), NewInt(3)})
	count = 0
	list.Range(func(key string, value Variant) bool {
		count++
		return true
	})
	if count != 3 {
		t.Errorf("expected 3 iterations for list, got %d", count)
	}
}

func TestRangeByIndex(t *testing.T) {
	list := NewValueList([]Variant{NewInt(10), NewInt(20), NewInt(30)})
	sum := 0
	list.RangeByIndex(func(index int, value Variant) bool {
		v, _ := value.AsInt()
		sum += v
		return true
	})
	if sum != 60 {
		t.Errorf("expected sum 60, got %d", sum)
	}
}

func TestGet(t *testing.T) {
	list := NewValueList([]Variant{NewInt(10), NewInt(20)})
	v, ok := list.Get(0)
	if !ok {
		t.Error("expected ok for index 0")
	}
	if i, _ := v.AsInt64(); i != 10 {
		t.Errorf("expected 10, got %d", i)
	}

	_, ok = list.Get(5)
	if ok {
		t.Error("expected not ok for out-of-bounds")
	}

	mp := NewValueMap(map[string]Variant{"key": NewString("val")})
	v2, ok := mp.Get("key")
	if !ok {
		t.Error("expected ok for key 'key'")
	}
	if v2.AsString() != "val" {
		t.Errorf("expected 'val', got '%s'", v2.AsString())
	}
}

func TestSet(t *testing.T) {
	list := NewValueList([]Variant{NewInt(10), NewInt(20)})
	list.Set(0, NewInt(99))
	v, _ := list.ListGet(0)
	if i, _ := v.AsInt64(); i != 99 {
		t.Errorf("expected 99, got %d", i)
	}

	// Set on map
	mp := NewValueMap(map[string]Variant{"key": NewInt(1)})
	mp.Set("key", NewInt(42))
	mapVal, ok := mp.MapGet("key")
	iv, _ := mapVal.AsInt64()
	if !ok || iv != 42 {
		t.Errorf("expected key='key'=42 after Set on map")
	}
}

func TestRemove(t *testing.T) {
	list := NewValueList([]Variant{NewInt(1), NewInt(2), NewInt(3)})
	list.Remove(1)
	if list.Len() != 2 {
		t.Errorf("expected len 2 after remove, got %d", list.Len())
	}
	v, _ := list.ListGet(1)
	if i, _ := v.AsInt64(); i != 3 {
		t.Errorf("expected 3 at index 1, got %d", i)
	}

	mp := NewValueMap(map[string]Variant{"a": NewInt(1), "b": NewInt(2)})
	mp.Remove("a")
	if mp.Len() != 1 {
		t.Errorf("expected len 1 after remove, got %d", mp.Len())
	}
	_, ok := mp.MapGet("a")
	if ok {
		t.Error("expected key 'a' to be removed")
	}

	iv := NewInt(42)
	err := iv.Remove(0)
	if err == nil {
		t.Error("expected error removing from non-container")
	}
}

func TestMapSet_OnEmpty(t *testing.T) {
	v := NewEmpty()
	v.MapSet("key", NewString("value"))
	if v.Type() != TypeMap {
		t.Fatalf("expected TypeMap after MapSet on empty, got %v", v.Type())
	}
	val, ok := v.MapGet("key")
	if !ok || val.AsString() != "value" {
		t.Error("MapSet on empty failed")
	}

	v.MapSet("key2", NewInt(100))
	if v.Len() != 2 {
		t.Errorf("expected len 2, got %d", v.Len())
	}
}

func TestListSet_Invalid(t *testing.T) {
	v := NewValueList([]Variant{NewInt(1)})
	v.ListSet(-1, NewInt(0))
	v.ListSet(10, NewInt(0))
	if v.Len() != 1 {
		t.Error("list should not change with invalid indices")
	}
}

func TestListGet_Invalid(t *testing.T) {
	v := NewValueList([]Variant{NewInt(1)})
	_, ok := v.ListGet(-1)
	if ok {
		t.Error("negative index should not be ok")
	}
	_, ok = v.ListGet(10)
	if ok {
		t.Error("out of bounds should not be ok")
	}

	nv := NewInt(42)
	got, ok := nv.ListGet(0)
	if ok {
		t.Error("non-list should not return ok")
	}
	if got.AsString() != "42" {
		t.Errorf("expected '42', got '%s'", got.AsString())
	}
}

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
