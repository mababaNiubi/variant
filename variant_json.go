package variant

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func (p Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.AsInterface())
}

func (m *Variant) UnmarshalJSON(data []byte) error {
	unmarshal, err := UnmarshalJSON(data)
	if err != nil {
		return err
	}
	m.variantType = unmarshal.variantType
	m.complexValue = unmarshal.complexValue
	m.numberValue = unmarshal.numberValue
	return nil
}

func UnmarshalJSON(data []byte) (Variant, error) {
	var value any
	err := json.Unmarshal(data, &value)
	if err != nil {
		return NewEmpty(), err
	}
	return decode(reflect.ValueOf(value))
}

func decode(rValue reflect.Value) (Variant, error) {
	switch rValue.Kind() {
	case reflect.String:
		return NewString(rValue.String()), nil
	case reflect.Bool:
		return NewBool(rValue.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewInt64(rValue.Int()), nil
	case reflect.Float64, reflect.Float32:
		return NewFloat64(rValue.Float()), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return NewUInt64(rValue.Uint()), nil
	case reflect.Array, reflect.Slice:
		list := make([]Variant, rValue.Len())
		for i := 0; i < rValue.Len(); i++ {
			v, err := decode(rValue.Index(i))
			if err != nil {
				return NewEmpty(), err
			}
			list[i] = v
		}
		return Variant{
			variantType:  TypeList,
			complexValue: list,
		}, nil
	case reflect.Map:
		mp := make(map[string]Variant)
		for _, key := range rValue.MapKeys() {
			v, err := decode(rValue.MapIndex(key))
			if err != nil {
				return NewEmpty(), err
			}
			mp[key.String()] = v
		}
		return Variant{
			variantType:  TypeMap,
			complexValue: mp,
		}, nil
	case reflect.Struct:
		mp := make(map[string]Variant)
		number := rValue.NumField()
		typeOf := rValue.Type()
		for i := 0; i < number; i++ {
			u, err := decode(rValue.Field(i))
			if err != nil {
				return NewEmpty(), err
			}
			key := typeOf.Field(i).Tag.Get("json")
			if len(key) == 0 {
				key = typeOf.Field(i).Name
			}
			mp[key] = u
		}
		return Variant{
			variantType:  TypeMap,
			complexValue: mp,
		}, nil
	case reflect.Interface, reflect.Pointer:
		return decode(rValue.Elem())
	case reflect.Invalid:
		return NewEmpty(), nil
	default:
		return NewEmpty(), fmt.Errorf("type not supported: %v", rValue.Kind())
	}
}
