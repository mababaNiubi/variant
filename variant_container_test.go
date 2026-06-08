package variant

import "testing"

func TestRangeByIndex_Map(t *testing.T) {
	mp := NewValueMap(map[string]Variant{"a": NewInt(10), "b": NewInt(20)})
	sum := 0
	mp.RangeByIndex(func(index int, value Variant) bool {
		v, _ := value.AsInt()
		sum += v
		return true
	})
	if sum != 30 {
		t.Errorf("expected sum 30, got %d", sum)
	}
}

func TestRangeByIndex_NonContainer(t *testing.T) {
	v := NewInt(42)
	called := false
	v.RangeByIndex(func(index int, value Variant) bool {
		called = true
		return true
	})
	if called {
		t.Error("RangeByIndex should not iterate on non-container")
	}
}

func TestGet_NonContainer(t *testing.T) {
	v := NewInt(42)
	_, ok := v.Get("anything")
	if ok {
		t.Error("Get on non-container should not return ok")
	}
}

func TestGet_ListNonIntKey(t *testing.T) {
	v := NewValueList([]Variant{NewInt(1), NewInt(2)})
	_, ok := v.Get("not-an-int")
	if ok {
		t.Error("Get on list with non-int key should not return ok")
	}
}

func TestGet_MapMissingKey(t *testing.T) {
	v := NewValueMap(map[string]Variant{"a": NewInt(1)})
	_, ok := v.Get("nonexistent")
	if ok {
		t.Error("Get with missing key should not return ok")
	}
}

func TestRemove_ListOutOfBounds(t *testing.T) {
	v := NewValueList([]Variant{NewInt(1), NewInt(2)})
	v.Remove(5) // out of bounds
	if v.Len() != 2 {
		t.Errorf("expected len 2, got %d", v.Len())
	}
}

func TestRemove_MapMissingKey(t *testing.T) {
	v := NewValueMap(map[string]Variant{"a": NewInt(1)})
	v.Remove("nonexistent")
	if v.Len() != 1 {
		t.Errorf("expected len 1, got %d", v.Len())
	}
}

func TestMapGet_NonMap(t *testing.T) {
	v := NewInt(42)
	_, ok := v.MapGet("key")
	if ok {
		t.Error("MapGet on non-map should not return ok")
	}
}

func TestMapSet_Default(t *testing.T) {
	v := NewInt(42)
	v.MapSet("key", NewInt(1))
	if v.Type() != TypeInt64 {
		t.Error("MapSet on int should be no-op")
	}
}

func TestRange_NonContainer(t *testing.T) {
	v := NewInt(42)
	called := false
	v.Range(func(key string, value Variant) bool {
		called = true
		return true
	})
	if called {
		t.Error("Range should not iterate on non-container")
	}
}

func TestListGet_NonListEmpty(t *testing.T) {
	v := NewEmpty()
	got, ok := v.ListGet(0)
	if ok {
		t.Error("ListGet index 0 on non-list should return ok=false")
	}
	if !got.IsEmpty() {
		t.Errorf("expected empty, got %v", got.Type())
	}
}

func TestListSet_NonList(t *testing.T) {
	v := NewString("hello")
	v.ListSet(0, NewInt(42))
	if v.AsString() != "hello" {
		t.Error("ListSet on non-list should be no-op")
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
