package variant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

func (p Variant) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 256)
	return p.appendJSON(buf)
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
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	v, err := decodeToken(dec)
	if err != nil {
		return NewEmpty(), err
	}
	return v, nil
}

const hexDigits = "0123456789abcdef"

func (v Variant) appendJSON(dst []byte) ([]byte, error) {
	switch v.variantType {
	case TypeEmpty:
		return append(dst, "null"...), nil
	case TypeBool:
		if v.numberValue != 0 {
			return append(dst, "true"...), nil
		}
		return append(dst, "false"...), nil
	case TypeInt64:
		return strconv.AppendInt(dst, v.numberValue, 10), nil
	case TypeUInt64:
		return strconv.AppendUint(dst, uint64(v.numberValue), 10), nil
	case TypeFloat64:
		f := *(*float64)(unsafe.Pointer(&v.numberValue))
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return dst, fmt.Errorf("json: unsupported value: %v", f)
		}
		return strconv.AppendFloat(dst, f, 'g', -1, 64), nil
	case TypeString:
		return appendJSONString(dst, v.complexValue.(string)), nil
	case TypeList:
		list := v.complexValue.([]Variant)
		dst = append(dst, '[')
		for i := range list {
			if i > 0 {
				dst = append(dst, ',')
			}
			var err error
			dst, err = list[i].appendJSON(dst)
			if err != nil {
				return dst, err
			}
		}
		return append(dst, ']'), nil
	case TypeMap:
		mp, _ := v.mapVariant()
		dst = append(dst, '{')
		first := true
		for k, val := range mp {
			if !first {
				dst = append(dst, ',')
			}
			first = false
			dst = appendJSONString(dst, k)
			dst = append(dst, ':')
			var err error
			dst, err = val.appendJSON(dst)
			if err != nil {
				return dst, err
			}
		}
		return append(dst, '}'), nil
	default:
		return append(dst, "null"...), nil
	}
}

func appendJSONString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			dst = append(dst, '\\', '"')
		case c == '\\':
			dst = append(dst, '\\', '\\')
		case c == '\n':
			dst = append(dst, '\\', 'n')
		case c == '\r':
			dst = append(dst, '\\', 'r')
		case c == '\t':
			dst = append(dst, '\\', 't')
		case c < 0x20:
			dst = append(dst, '\\', 'u', '0', '0', hexDigits[c>>4], hexDigits[c&0xF])
		case c == '<' || c == '>' || c == '&':
			dst = append(dst, '\\', 'u', '0', '0', hexDigits[c>>4], hexDigits[c&0xF])
		default:
			dst = append(dst, c)
		}
	}
	return append(dst, '"')
}

func decodeToken(dec *json.Decoder) (Variant, error) {
	t, err := dec.Token()
	if err != nil {
		return NewEmpty(), err
	}
	switch v := t.(type) {
	case nil:
		return NewEmpty(), nil
	case bool:
		return NewBool(v), nil
	case json.Number:
		return decodeJSONNumber(string(v))
	case string:
		return NewString(v), nil
	case json.Delim:
		if v == '[' {
			return decodeJSONArray(dec)
		}
		if v == '{' {
			return decodeJSONObject(dec)
		}
	}
	return NewEmpty(), fmt.Errorf("unexpected JSON token: %v", t)
}

func decodeJSONNumber(s string) (Variant, error) {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return NewInt64(i), nil
	}
	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		return NewUInt64(u), nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return NewEmpty(), fmt.Errorf("invalid JSON number: %s", s)
	}
	return NewFloat64(f), nil
}

func decodeJSONArray(dec *json.Decoder) (Variant, error) {
	var list []Variant
	for dec.More() {
		val, err := decodeToken(dec)
		if err != nil {
			return NewEmpty(), err
		}
		list = append(list, val)
	}
	if _, err := dec.Token(); err != nil {
		return NewEmpty(), err
	}
	return NewValueList(list), nil
}

func decodeJSONObject(dec *json.Decoder) (Variant, error) {
	mp := make(map[string]Variant)
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return NewEmpty(), err
		}
		key, ok := keyToken.(string)
		if !ok {
			return NewEmpty(), fmt.Errorf("expected string key, got %T", keyToken)
		}
		val, err := decodeToken(dec)
		if err != nil {
			return NewEmpty(), err
		}
		mp[key] = val
	}
	if _, err := dec.Token(); err != nil {
		return NewEmpty(), err
	}
	return NewValueMap(mp), nil
}

func looksLikeJSON(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return s[i] == '{' || s[i] == '['
		}
	}
	return false
}

func decodeSlice(v []any) Variant {
	list := make([]Variant, len(v))
	for i := range v {
		list[i] = New(v[i])
	}
	return Variant{variantType: TypeList, complexValue: list}
}

func decodeMap(v map[string]any) Variant {
	mp := make(map[string]Variant, len(v))
	for k, val := range v {
		mp[k] = New(val)
	}
	return Variant{variantType: TypeMap, complexValue: mp}
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
