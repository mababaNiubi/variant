package variant

import "testing"

func TestUnmarshalTo_BasicTypes(t *testing.T) {
	var b bool
	if err := NewBool(true).UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if b != true {
		t.Error("expected true")
	}

	var i int
	if err := NewInt(42).UnmarshalTo(&i); err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Errorf("expected 42, got %d", i)
	}

	var i64 int64
	if err := NewInt64(-999).UnmarshalTo(&i64); err != nil {
		t.Fatal(err)
	}
	if i64 != -999 {
		t.Errorf("expected -999, got %d", i64)
	}

	var u64 uint64
	if err := NewUInt64(18446744073709551615).UnmarshalTo(&u64); err != nil {
		t.Fatal(err)
	}
	if u64 != 18446744073709551615 {
		t.Errorf("expected 18446744073709551615, got %d", u64)
	}

	var f64 float64
	if err := NewFloat64(3.14159).UnmarshalTo(&f64); err != nil {
		t.Fatal(err)
	}
	if f64 != 3.14159 {
		t.Errorf("expected 3.14159, got %f", f64)
	}

	var f32 float32
	if err := NewFloat64(2.5).UnmarshalTo(&f32); err != nil {
		t.Fatal(err)
	}
	if f32 != 2.5 {
		t.Errorf("expected 2.5, got %f", f32)
	}

	var s string
	if err := NewString("hello").UnmarshalTo(&s); err != nil {
		t.Fatal(err)
	}
	if s != "hello" {
		t.Errorf("expected hello, got %s", s)
	}

	var emptyInt int
	emptyInt = 5
	if err := NewEmpty().UnmarshalTo(&emptyInt); err != nil {
		t.Fatal(err)
	}
	if emptyInt != 0 {
		t.Errorf("expected 0 after empty, got %d", emptyInt)
	}
}

func TestUnmarshalTo_CrossTypeConversion(t *testing.T) {
	var f float64
	if err := NewInt(42).UnmarshalTo(&f); err != nil {
		t.Fatal(err)
	}
	if f != 42.0 {
		t.Errorf("expected 42.0, got %f", f)
	}

	var i int
	if err := NewFloat64(3.9).UnmarshalTo(&i); err != nil {
		t.Fatal(err)
	}
	if i != 3 {
		t.Errorf("expected 3, got %d", i)
	}

	var u uint64
	if err := NewInt(100).UnmarshalTo(&u); err != nil {
		t.Fatal(err)
	}
	if u != 100 {
		t.Errorf("expected 100, got %d", u)
	}

	var b bool
	if err := NewInt(1).UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("expected true for non-zero int")
	}
	if err := NewInt(0).UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("expected false for zero int")
	}

	var bi int
	if err := NewBool(true).UnmarshalTo(&bi); err != nil {
		t.Fatal(err)
	}
	if bi != 1 {
		t.Errorf("expected 1, got %d", bi)
	}

	var s string
	if err := NewInt(123).UnmarshalTo(&s); err != nil {
		t.Fatal(err)
	}
	if s != "123" {
		t.Errorf("expected '123', got '%s'", s)
	}

	var si int
	if err := NewString("456").UnmarshalTo(&si); err != nil {
		t.Fatal(err)
	}
	if si != 456 {
		t.Errorf("expected 456, got %d", si)
	}

	var sf float64
	if err := NewString("3.14").UnmarshalTo(&sf); err != nil {
		t.Fatal(err)
	}
	if sf != 3.14 {
		t.Errorf("expected 3.14, got %f", sf)
	}

	var sb bool
	if err := NewString("true").UnmarshalTo(&sb); err != nil {
		t.Fatal(err)
	}
	if !sb {
		t.Error("expected true")
	}
}

func TestUnmarshalTo_ListToSlice(t *testing.T) {
	v := NewValueList([]Variant{NewInt(1), NewInt(2), NewInt(3)})

	var s []int
	if err := v.UnmarshalTo(&s); err != nil {
		t.Fatal(err)
	}
	if len(s) != 3 || s[0] != 1 || s[1] != 2 || s[2] != 3 {
		t.Errorf("expected [1,2,3], got %v", s)
	}

	var fs []float64
	if err := NewValueList([]Variant{NewFloat64(1.1), NewFloat64(2.2)}).UnmarshalTo(&fs); err != nil {
		t.Fatal(err)
	}
	if len(fs) != 2 || fs[0] != 1.1 || fs[1] != 2.2 {
		t.Errorf("expected [1.1,2.2], got %v", fs)
	}

	var ss []string
	if err := NewValueList([]Variant{NewString("a"), NewString("b")}).UnmarshalTo(&ss); err != nil {
		t.Fatal(err)
	}
	if len(ss) != 2 || ss[0] != "a" || ss[1] != "b" {
		t.Errorf("expected [a,b], got %v", ss)
	}
}

func TestUnmarshalTo_ListToArray(t *testing.T) {
	v := NewValueList([]Variant{NewInt(10), NewInt(20), NewInt(30)})

	var arr [3]int
	if err := v.UnmarshalTo(&arr); err != nil {
		t.Fatal(err)
	}
	if arr != [3]int{10, 20, 30} {
		t.Errorf("expected [10,20,30], got %v", arr)
	}

	var arr2 [5]int
	for i := range arr2 {
		arr2[i] = 99
	}
	if err := NewValueList([]Variant{NewInt(1)}).UnmarshalTo(&arr2); err != nil {
		t.Fatal(err)
	}
	if arr2[0] != 1 || arr2[1] != 99 {
		t.Errorf("expected [1,99,99,99,99], got %v", arr2)
	}

	var arr3 [2]int
	if err := NewValueList([]Variant{NewInt(1), NewInt(2), NewInt(3)}).UnmarshalTo(&arr3); err != nil {
		t.Fatal(err)
	}
	if arr3 != [2]int{1, 2} {
		t.Errorf("expected [1,2], got %v", arr3)
	}
}

func TestUnmarshalTo_MapToMap(t *testing.T) {
	v := NewValueMap(map[string]Variant{
		"one": NewInt(1),
		"two": NewFloat64(2.2),
	})

	var m map[string]int
	if err := v.UnmarshalTo(&m); err != nil {
		t.Fatal(err)
	}
	if m["one"] != 1 || m["two"] != 2 {
		t.Errorf("expected {one:1, two:2}, got %v", m)
	}

	va := NewValueMap(map[string]Variant{
		"a": NewInt(42),
		"b": NewString("hello"),
		"c": NewFloat64(3.14),
	})
	var ma map[string]any
	if err := va.UnmarshalTo(&ma); err != nil {
		t.Fatal(err)
	}
	if len(ma) != 3 {
		t.Errorf("expected 3 entries, got %d", len(ma))
	}
}

func TestUnmarshalTo_MapToStruct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	v := NewValueMap(map[string]Variant{
		"name": NewString("Alice"),
		"age":  NewInt(30),
	})

	var p Person
	if err := v.UnmarshalTo(&p); err != nil {
		t.Fatal(err)
	}
	if p.Name != "Alice" || p.Age != 30 {
		t.Errorf("expected {Name:Alice, Age:30}, got %+v", p)
	}
}

func TestUnmarshalTo_StructFieldNameFallback(t *testing.T) {
	type Data struct {
		X int
		Y int
	}

	v := NewValueMap(map[string]Variant{
		"X": NewInt(1),
		"Y": NewInt(2),
	})

	var d Data
	if err := v.UnmarshalTo(&d); err != nil {
		t.Fatal(err)
	}
	if d.X != 1 || d.Y != 2 {
		t.Errorf("expected {X:1, Y:2}, got %+v", d)
	}
}

func TestUnmarshalTo_NestedStruct(t *testing.T) {
	type Inner struct {
		Value float64 `json:"value"`
	}
	type Outer struct {
		Title string `json:"title"`
		Nest  Inner  `json:"nest"`
	}

	inner := NewValueMap(map[string]Variant{"value": NewFloat64(3.14)})
	outer := NewValueMap(map[string]Variant{
		"title": NewString("test"),
		"nest":  inner,
	})

	var o Outer
	if err := outer.UnmarshalTo(&o); err != nil {
		t.Fatal(err)
	}
	if o.Title != "test" || o.Nest.Value != 3.14 {
		t.Errorf("expected {Title:test, Nest:{Value:3.14}}, got %+v", o)
	}
}

func TestUnmarshalTo_StructWithSlice(t *testing.T) {
	type Item struct {
		Tags []string `json:"tags"`
	}

	v := NewValueMap(map[string]Variant{
		"tags": NewValueList([]Variant{NewString("a"), NewString("b"), NewString("c")}),
	})

	var item Item
	if err := v.UnmarshalTo(&item); err != nil {
		t.Fatal(err)
	}
	if len(item.Tags) != 3 || item.Tags[0] != "a" || item.Tags[1] != "b" || item.Tags[2] != "c" {
		t.Errorf("expected [a,b,c], got %v", item.Tags)
	}
}

func TestUnmarshalTo_ErrorNonPointer(t *testing.T) {
	var i int
	err := NewInt(1).UnmarshalTo(i)
	if err == nil {
		t.Error("expected error for non-pointer")
	}
}

func TestUnmarshalTo_ErrorNilPointer(t *testing.T) {
	var p *int
	err := NewInt(1).UnmarshalTo(p)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}

func TestUnmarshalTo_ErrorStringToInvalidInt(t *testing.T) {
	var i int
	err := NewString("not-a-number").UnmarshalTo(&i)
	if err == nil {
		t.Error("expected error for invalid string-to-int conversion")
	}
}

func TestUnmarshalTo_InterfaceTarget(t *testing.T) {
	var a any
	if err := NewInt(42).UnmarshalTo(&a); err != nil {
		t.Fatal(err)
	}
	if a != int64(42) {
		t.Errorf("expected int64(42), got %T(%v)", a, a)
	}

	var a2 any
	if err := NewFloat64(3.14).UnmarshalTo(&a2); err != nil {
		t.Fatal(err)
	}
	if a2 != 3.14 {
		t.Errorf("expected 3.14, got %v", a2)
	}

	var a3 any
	if err := NewString("hello").UnmarshalTo(&a3); err != nil {
		t.Fatal(err)
	}
	if a3 != "hello" {
		t.Errorf("expected hello, got %v", a3)
	}
}

func TestUnmarshalTo_RoundTripStruct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	original := Person{Name: "Bob", Age: 25}

	v := New(original)

	var restored Person
	if err := v.UnmarshalTo(&restored); err != nil {
		t.Fatal(err)
	}
	if restored != original {
		t.Errorf("roundtrip mismatch: %+v != %+v", restored, original)
	}
}

func TestUnmarshalTo_SkipUnknownJSONFields(t *testing.T) {
	type Target struct {
		Name string `json:"name"`
	}

	v := NewValueMap(map[string]Variant{
		"name":   NewString("Alice"),
		"extra1": NewInt(123),
		"extra2": NewString("ignored"),
	})

	var tgt Target
	if err := v.UnmarshalTo(&tgt); err != nil {
		t.Fatal(err)
	}
	if tgt.Name != "Alice" {
		t.Errorf("expected Alice, got %s", tgt.Name)
	}
}

func TestUnmarshalTo_UnexportedFieldsIgnored(t *testing.T) {
	type Data struct {
		Exported   string `json:"exported"`
		unexported string `json:"unexported"`
	}

	v := NewValueMap(map[string]Variant{
		"exported":   NewString("visible"),
		"unexported": NewString("hidden"),
	})

	var d Data
	if err := v.UnmarshalTo(&d); err != nil {
		t.Fatal(err)
	}
	if d.Exported != "visible" {
		t.Errorf("expected visible, got %s", d.Exported)
	}
	if d.unexported != "" {
		t.Errorf("expected empty for unexported field, got %s", d.unexported)
	}
}

func TestUnmarshalTo_MapIntKey(t *testing.T) {
	v := NewValueMap(map[string]Variant{"0": NewString("a"), "1": NewString("b")})

	var m map[int]string
	if err := v.UnmarshalTo(&m); err != nil {
		t.Fatal(err)
	}
	if m[0] != "a" || m[1] != "b" {
		t.Errorf("expected {0:a, 1:b}, got %v", m)
	}
}

func TestUnmarshalTo_SetBool_AllTargets(t *testing.T) {
	var s string
	if err := NewBool(true).UnmarshalTo(&s); err != nil {
		t.Fatal(err)
	}
	if s != "true" {
		t.Errorf("expected 'true', got '%s'", s)
	}

	var f float64
	if err := NewBool(false).UnmarshalTo(&f); err != nil {
		t.Fatal(err)
	}
	if f != 0 {
		t.Errorf("expected 0, got %f", f)
	}
}

func TestUnmarshalTo_SetUInt64_AllTargets(t *testing.T) {
	var s string
	if err := NewUInt64(99).UnmarshalTo(&s); err != nil {
		t.Fatal(err)
	}
	if s != "99" {
		t.Errorf("expected '99', got '%s'", s)
	}

	var b bool
	if err := NewUInt64(1).UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("expected true")
	}
}

func TestUnmarshalTo_SetFloat64_AllTargets(t *testing.T) {
	var s string
	if err := NewFloat64(3.14159).UnmarshalTo(&s); err != nil {
		t.Fatal(err)
	}
	if s != "3.14159" {
		t.Errorf("expected '3.14159', got '%s'", s)
	}

	var b bool
	if err := NewFloat64(1.5).UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("expected true for non-zero float")
	}
}

func TestUnmarshalTo_SetString_AllTargets(t *testing.T) {
	var u uint
	if err := NewString("42").UnmarshalTo(&u); err != nil {
		t.Fatal(err)
	}
	if u != 42 {
		t.Errorf("expected 42, got %d", u)
	}

	var b bool
	if err := NewString("true").UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("expected true")
	}
}

func TestUnmarshalTo_ListToInterface(t *testing.T) {
	v := NewValueList([]Variant{NewInt(1), NewInt(2)})
	var a any
	if err := v.UnmarshalTo(&a); err != nil {
		t.Fatal(err)
	}
	list, ok := a.([]Variant)
	if !ok || len(list) != 2 {
		t.Errorf("expected []Variant of len 2, got %T", a)
	}
}

func TestUnmarshalTo_MapToInterface(t *testing.T) {
	v := NewValueMap(map[string]Variant{"x": NewInt(1)})
	var a any
	if err := v.UnmarshalTo(&a); err != nil {
		t.Fatal(err)
	}
	mp, ok := a.(map[string]Variant)
	if !ok || len(mp) != 1 {
		t.Errorf("expected map[string]Variant of len 1, got %T", a)
	}
}

func TestUnmarshalTo_BoolToString(t *testing.T) {
	var s string
	if err := NewBool(true).UnmarshalTo(&s); err != nil {
		t.Fatal(err)
	}
	if s != "true" {
		t.Errorf("expected 'true', got '%s'", s)
	}
}

func TestUnmarshalTo_Uint64ToFloat(t *testing.T) {
	var f float64
	if err := NewUInt64(42).UnmarshalTo(&f); err != nil {
		t.Fatal(err)
	}
	if f != 42.0 {
		t.Errorf("expected 42.0, got %f", f)
	}
}

func TestUnmarshalTo_Uint64ToBool(t *testing.T) {
	var b bool
	if err := NewUInt64(0).UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("expected false for zero uint")
	}
}

func TestUnmarshalTo_StringToUint(t *testing.T) {
	var u uint64
	if err := NewString("100").UnmarshalTo(&u); err != nil {
		t.Fatal(err)
	}
	if u != 100 {
		t.Errorf("expected 100, got %d", u)
	}
}

func TestUnmarshalTo_StringToBool(t *testing.T) {
	var b bool
	if err := NewString("1").UnmarshalTo(&b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("expected true for string '1'")
	}
}

func TestUnmarshalTo_MapToStructWithIntConversion(t *testing.T) {
	type Data struct {
		Score float64 `json:"score"`
	}
	v := NewValueMap(map[string]Variant{"score": NewInt(95)})
	var d Data
	if err := v.UnmarshalTo(&d); err != nil {
		t.Fatal(err)
	}
	if d.Score != 95.0 {
		t.Errorf("expected 95.0, got %f", d.Score)
	}
}

func TestUnmarshalTo_SetBool_ToUint(t *testing.T) {
	var u uint64
	if err := NewBool(true).UnmarshalTo(&u); err != nil {
		t.Fatal(err)
	}
	if u != 1 {
		t.Errorf("expected 1, got %d", u)
	}
}

func TestUnmarshalTo_UintViaInterface(t *testing.T) {
	var a any
	if err := NewUInt64(42).UnmarshalTo(&a); err != nil {
		t.Fatal(err)
	}
	if a != uint64(42) {
		t.Errorf("expected uint64(42), got %T(%v)", a, a)
	}
}

func TestUnmarshalTo_FloatViaInterface(t *testing.T) {
	var a any
	if err := NewFloat64(2.5).UnmarshalTo(&a); err != nil {
		t.Fatal(err)
	}
	if a != 2.5 {
		t.Errorf("expected 2.5, got %v", a)
	}
}

func TestUnmarshalTo_StringViaInterface(t *testing.T) {
	var a any
	if err := NewString("hello").UnmarshalTo(&a); err != nil {
		t.Fatal(err)
	}
	if a != "hello" {
		t.Errorf("expected hello, got %v", a)
	}
}

func TestUnmarshalTo_BoolViaInterface(t *testing.T) {
	var a any
	if err := NewBool(true).UnmarshalTo(&a); err != nil {
		t.Fatal(err)
	}
	if a != true {
		t.Errorf("expected true, got %v", a)
	}
}

func TestUnmarshalTo_DifferentIntSizes(t *testing.T) {
	v := NewInt(127)

	var i8 int8
	if err := v.UnmarshalTo(&i8); err != nil {
		t.Fatal(err)
	}
	if i8 != 127 {
		t.Errorf("expected 127, got %d", i8)
	}

	var i16 int16
	if err := v.UnmarshalTo(&i16); err != nil {
		t.Fatal(err)
	}
	if i16 != 127 {
		t.Errorf("expected 127, got %d", i16)
	}

	var i32 int32
	if err := v.UnmarshalTo(&i32); err != nil {
		t.Fatal(err)
	}
	if i32 != 127 {
		t.Errorf("expected 127, got %d", i32)
	}
}
