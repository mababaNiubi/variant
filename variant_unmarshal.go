package variant

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// UnmarshalTo populates target with the Variant value, similar to json.Unmarshal.
// target must be a non-nil pointer. Struct fields are matched by json tag (preferred)
// or field name. Variant type is converted to match the target field type automatically.
func (v Variant) UnmarshalTo(target any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("variant: UnmarshalTo requires a non-nil pointer")
	}
	return v.setValue(rv.Elem())
}

func (v Variant) setValue(rv reflect.Value) error {
	if !rv.CanSet() {
		return fmt.Errorf("variant: cannot set value")
	}

	switch v.variantType {
	case TypeEmpty:
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	case TypeBool:
		return v.setBool(rv)
	case TypeInt64:
		return v.setInt64(rv)
	case TypeUInt64:
		return v.setUInt64(rv)
	case TypeFloat64:
		return v.setFloat64(rv)
	case TypeString:
		return v.setString(rv)
	case TypeList:
		return v.setList(rv)
	case TypeMap:
		return v.setMap(rv)
	default:
		return fmt.Errorf("variant: unknown type %v", v.variantType)
	}
}

func (v Variant) setBool(rv reflect.Value) error {
	b := v.numberValue != 0
	switch rv.Kind() {
	case reflect.Bool:
		rv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if b {
			rv.SetInt(1)
		} else {
			rv.SetInt(0)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if b {
			rv.SetUint(1)
		} else {
			rv.SetUint(0)
		}
	case reflect.Float32, reflect.Float64:
		if b {
			rv.SetFloat(1)
		} else {
			rv.SetFloat(0)
		}
	case reflect.String:
		rv.SetString(strconv.FormatBool(b))
	default:
		return v.setViaInterface(rv, b)
	}
	return nil
}

func (v Variant) setInt64(rv reflect.Value) error {
	n := v.numberValue
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(uint64(n))
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(float64(n))
	case reflect.Bool:
		rv.SetBool(n != 0)
	case reflect.String:
		rv.SetString(strconv.FormatInt(n, 10))
	default:
		return v.setViaInterface(rv, n)
	}
	return nil
}

func (v Variant) setUInt64(rv reflect.Value) error {
	n := uint64(v.numberValue)
	switch rv.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(int64(n))
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(float64(n))
	case reflect.Bool:
		rv.SetBool(n != 0)
	case reflect.String:
		rv.SetString(strconv.FormatUint(n, 10))
	default:
		return v.setViaInterface(rv, n)
	}
	return nil
}

func (v Variant) setFloat64(rv reflect.Value) error {
	f := *(*float64)(unsafe.Pointer(&v.numberValue))
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(int64(f))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(uint64(f))
	case reflect.Bool:
		rv.SetBool(f != 0)
	case reflect.String:
		rv.SetString(strconv.FormatFloat(f, 'f', -1, 64))
	default:
		return v.setViaInterface(rv, f)
	}
	return nil
}

func (v Variant) setString(rv reflect.Value) error {
	s := v.complexValue.(string)
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("variant: cannot convert string %q to int: %w", s, err)
		}
		rv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("variant: cannot convert string %q to uint: %w", s, err)
		}
		rv.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("variant: cannot convert string %q to float: %w", s, err)
		}
		rv.SetFloat(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("variant: cannot convert string %q to bool: %w", s, err)
		}
		rv.SetBool(b)
	default:
		return v.setViaInterface(rv, s)
	}
	return nil
}

func (v Variant) setList(rv reflect.Value) error {
	list, _ := v.complexValue.([]Variant)
	switch rv.Kind() {
	case reflect.Slice:
		nv := reflect.MakeSlice(rv.Type(), len(list), len(list))
		for i := range list {
			if err := list[i].setValue(nv.Index(i)); err != nil {
				return fmt.Errorf("variant: list[%d]: %w", i, err)
			}
		}
		rv.Set(nv)
	case reflect.Array:
		length := rv.Len()
		for i := 0; i < len(list) && i < length; i++ {
			if err := list[i].setValue(rv.Index(i)); err != nil {
				return fmt.Errorf("variant: list[%d]: %w", i, err)
			}
		}
	case reflect.Interface:
		rv.Set(reflect.ValueOf(list))
	default:
		return fmt.Errorf("variant: cannot unmarshal list into %s", rv.Kind())
	}
	return nil
}

func (v Variant) setMap(rv reflect.Value) error {
	mp, _ := v.mapVariant()
	switch rv.Kind() {
	case reflect.Map:
		nv := reflect.MakeMapWithSize(rv.Type(), len(mp))
		kt := rv.Type().Key()
		vt := rv.Type().Elem()
		for key, val := range mp {
			kv := reflect.New(kt).Elem()
			if err := setReflectValueFromString(key, kv); err != nil {
				return fmt.Errorf("variant: map key %q: %w", key, err)
			}
			vv := reflect.New(vt).Elem()
			if err := val.setValue(vv); err != nil {
				return fmt.Errorf("variant: map[%q]: %w", key, err)
			}
			nv.SetMapIndex(kv, vv)
		}
		rv.Set(nv)
	case reflect.Struct:
		t := rv.Type()
		for key, val := range mp {
			field := findFieldByJSONTag(t, key)
			if field.Name == "" {
				continue
			}
			if err := val.setValue(rv.FieldByIndex(field.Index)); err != nil {
				return fmt.Errorf("variant: field %s: %w", key, err)
			}
		}
	case reflect.Interface:
		rv.Set(reflect.ValueOf(mp))
	default:
		return fmt.Errorf("variant: cannot unmarshal map into %s", rv.Kind())
	}
	return nil
}

func findFieldByJSONTag(t reflect.Type, tag string) reflect.StructField {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		if f.Tag.Get("json") == tag {
			return f
		}
	}
	// fallback: match by field name
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		if f.Name == tag {
			return f
		}
	}
	return reflect.StructField{}
}

func setReflectValueFromString(s string, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		rv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		rv.SetUint(n)
	default:
		rv.SetString(s)
	}
	return nil
}

// setViaInterface uses AsInterface() to get a Go-native value and sets it.
// This handles edge cases like any/interface{} targets.
func (v Variant) setViaInterface(rv reflect.Value, val any) error {
	if rv.Kind() == reflect.Interface && rv.NumMethod() == 0 {
		rv.Set(reflect.ValueOf(val))
		return nil
	}
	return fmt.Errorf("variant: cannot unmarshal into %s", rv.Kind())
}
