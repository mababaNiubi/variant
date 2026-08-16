package variant

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

type Variant struct {
	variantType  Type
	numberValue  int64 // 数字类型用该变量存储 （采集场景下，大部分是数字类型，故不用interface{}，避免高频操作的性能影响）
	complexValue any   // 字符串、list、map 用该变量存储
}

// structPairs is the compact representation of a TypeMap structure: parallel
// key/value slices. Building and iterating a small structure costs a fraction
// of a map (no hashing, no growth rehash), which is what makes high-volume
// structure writes and reads viable. A map never aliases the caller's input,
// so callers may freely reuse their map after New.
type structPairs struct {
	keys []string
	vals []Variant
}

func (p *structPairs) len() int {
	if p == nil {
		return 0
	}
	return len(p.keys)
}

func (p *structPairs) get(key string) (Variant, bool) {
	for i, k := range p.keys {
		if k == key {
			return p.vals[i], true
		}
	}
	return NewEmpty(), false
}

func NewEmpty() Variant {
	return Variant{
		variantType: TypeEmpty,
	}
}

func New(v any) Variant {
	switch val := v.(type) {
	case bool:
		return NewBool(val)
	case string:
		switch GetStringValueType(val) {
		case TypeInt64:
			ato, err := strconv.Atoi(val)
			if err != nil {
				return NewString(val)
			}
			return NewInt(ato)
		case TypeEmpty:
			return NewEmpty()
		case TypeFloat64:
			ato, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return NewString(val)
			}
			return NewFloat64(ato)
		default:
			if looksLikeJSON(val) && json.Valid([]byte(val)) {
				va, _ := UnmarshalJSON([]byte(val))
				return va
			}
		}
		return NewString(val)
	case float64:
		return NewFloat64(val)
	case float32:
		return NewFloat64(float64(val))
	case int:
		return NewInt(val)
	case int8:
		return NewInt(int(val))
	case int16:
		return NewInt(int(val))
	case int32:
		return NewInt(int(val))
	case int64:
		return NewInt64(val)
	case uint:
		return NewUInt64(uint64(val))
	case uint64:
		return NewUInt64(val)
	case uint32:
		return NewUInt64(uint64(val))
	case uint16:
		return NewUInt64(uint64(val))
	case uint8:
		return NewUInt64(uint64(val))
	case []byte:
		if json.Valid(val) {
			va, _ := UnmarshalJSON(val)
			return va
		}
		va, _ := decode(reflect.ValueOf(v))
		return va
	case Variant:
		return val
	case *Variant:
		return *val
	case []Variant:
		return NewValueList(val)
	case map[string]Variant:
		return NewValueMap(val)
	case []any:
		return decodeSlice(val)
	case map[string]any:
		// Compact: convert each value to a Variant (NewRawValue wraps by Go
		// type, no string reclassification) and store as key/value pair slices.
		// Values are copied into the slices, so callers may freely reuse or
		// mutate their map afterwards.
		sp := &structPairs{keys: make([]string, 0, len(val)), vals: make([]Variant, 0, len(val))}
		for k, v := range val {
			sp.keys = append(sp.keys, k)
			sp.vals = append(sp.vals, NewRawValue(v))
		}
		return Variant{variantType: TypeMap, complexValue: sp}
	default:
		va, _ := decode(reflect.ValueOf(v))
		return va
	}
}

func NewBool(v bool) Variant {
	if v {
		return Variant{
			variantType: TypeBool,
			numberValue: 1,
		}
	}

	return Variant{
		variantType: TypeBool,
		numberValue: 0,
	}
}

func NewInt(v int) (r Variant) {
	return Variant{
		variantType: TypeInt64,
		numberValue: int64(v),
	}
}

func NewInt64(v int64) (r Variant) {
	return Variant{
		variantType: TypeInt64,
		numberValue: v,
	}
}

func NewUInt64(v uint64) Variant {
	return Variant{
		variantType: TypeUInt64,
		numberValue: int64(v),
	}
}

func NewValue(v float64) (r Variant) {
	_, sp := math.Modf(v)
	if sp == 0 && (v <= float64(int64(math.MaxInt64)) && v >= float64(int64(math.MinInt64))) {
		return NewInt64(int64(v))
	}
	return NewFloat64(v)
}

func NewFloat64(v float64) (r Variant) {
	r.variantType = TypeFloat64
	*(*float64)(unsafe.Pointer(&r.numberValue)) = v
	return r
}

func NewString(value string) Variant {
	return Variant{
		variantType:  TypeString,
		complexValue: value,
	}
}

func NewValueList(v []Variant) Variant {
	return Variant{
		variantType:  TypeList,
		complexValue: v,
	}
}

func NewValueMap(v map[string]Variant) Variant {
	sp := &structPairs{keys: make([]string, 0, len(v)), vals: make([]Variant, 0, len(v))}
	for k, val := range v {
		sp.keys = append(sp.keys, k)
		sp.vals = append(sp.vals, val)
	}
	return Variant{
		variantType:  TypeMap,
		complexValue: sp,
	}
}

// NewStruct builds a compact structure from parallel key/value slices without
// allocating a map. Keys and vals must be the same length. Used by the columnar
// read path to reconstruct structures cheaply.
func NewStruct(keys []string, vals []Variant) Variant {
	return Variant{
		variantType:  TypeMap,
		complexValue: &structPairs{keys: keys, vals: vals},
	}
}

// NewRawValue wraps a raw structure field value as a Variant using its Go type
// as authoritative. Unlike New, a Go string is never reclassified as a number:
// in a raw map[string]any structure the value's concrete type already tells us
// what it is. Used by the WAL encoder and the columnar write path so the live
// and replayed representations agree.
func NewRawValue(val any) Variant {
	switch t := val.(type) {
	case bool:
		return NewBool(t)
	case float64:
		return NewFloat64(t)
	case float32:
		return NewFloat64(float64(t))
	case int:
		return NewInt(t)
	case int8:
		return NewInt(int(t))
	case int16:
		return NewInt(int(t))
	case int32:
		return NewInt(int(t))
	case int64:
		return NewInt64(t)
	case uint:
		return NewUInt64(uint64(t))
	case uint8:
		return NewUInt64(uint64(t))
	case uint16:
		return NewUInt64(uint64(t))
	case uint32:
		return NewUInt64(uint64(t))
	case uint64:
		return NewUInt64(t)
	case string:
		return NewString(t)
	case []byte:
		return NewString(string(t))
	default:
		// Nested structures, lists, and unknown types fall through to New.
		return New(val)
	}
}

// mapVariant converts a structure into a map[string]Variant. Consumers that
// only need to iterate use StructPairs or Range directly; this exists for the
// few callers that require map semantics (mutation by key in container ops).
func (v Variant) mapVariant() (map[string]Variant, bool) {
	sp, ok := v.complexValue.(*structPairs)
	if !ok || sp == nil {
		return nil, false
	}
	mp := make(map[string]Variant, len(sp.keys))
	for i := range sp.keys {
		mp[sp.keys[i]] = sp.vals[i]
	}
	return mp, true
}

// StructPairs exposes the compact key/value slices backing a TypeMap.
func (v Variant) StructPairs() (*structPairs, bool) {
	sp, ok := v.complexValue.(*structPairs)
	return sp, ok
}

func (v Variant) Type() Type {
	return v.variantType
}

func (v Variant) IsEmpty() bool {
	switch v.variantType {
	case TypeEmpty:
		return true
	case TypeString:
		return len(v.complexValue.(string)) == 0
	case TypeList:
		vs := v.complexValue.([]Variant)
		num := 0
		for i := range vs {
			if vs[i].IsEmpty() {
				num++
			}
		}
		return num == len(vs)
	case TypeMap:
		// Length check only — IsEmpty runs on every written structure value.
		sp, _ := v.StructPairs()
		return sp.len() == 0
	default:
		return false
	}
}

func (v Variant) IsZero() bool {
	switch v.variantType {
	case TypeString:
		asFloat64, err := v.AsFloat64()
		if err != nil {
			return false
		}
		return asFloat64 == 0
	case TypeInt64, TypeUInt64, TypeFloat64:
		return v.numberValue == 0
	default:
		return false
	}
}

func (v Variant) IsNumber() bool {
	switch v.variantType {
	case TypeInt64, TypeUInt64, TypeFloat64:
		return true
	case TypeString:
		return IsNumber(v.complexValue.(string))
	default:
		return false
	}
}

func (v Variant) IsFloat() bool {
	switch v.variantType {
	case TypeFloat64:
		return true
	case TypeString:
		return IsNumber(v.complexValue.(string)) && strings.IndexByte(v.complexValue.(string), '.') != -1
	default:
		return false
	}
}

func (v Variant) IsTrue() bool {
	asBool, err := v.AsBool()
	if err != nil {
		return false
	}
	return asBool
}
func (v Variant) AsBool() (bool, error) {
	switch v.variantType {
	case TypeEmpty:
		return false, nil
	case TypeBool, TypeInt64, TypeUInt64:
		if v.numberValue == 0 {
			return false, nil
		}
		return true, nil
	case TypeFloat64:
		if *(*float64)(unsafe.Pointer(&v.numberValue)) == 0 {
			return false, nil
		}
		return true, nil
	case TypeString:
		return strconv.ParseBool(v.complexValue.(string))
	default:
		return false, errors.New(errorVariantTypeNotNumber)
	}
}

func (v Variant) AsInt64() (int64, error) {
	switch v.variantType {
	case TypeBool, TypeInt64, TypeUInt64, TypeEmpty:
		return v.numberValue, nil
	case TypeFloat64:
		f := *(*float64)(unsafe.Pointer(&v.numberValue))
		if f < math.MinInt64 || f > math.MaxInt64 {
			return 0, errors.New(errorVariantValueOverFlow)
		}
		return int64(f), nil
	case TypeString:
		i, err := strconv.ParseFloat(v.complexValue.(string), 64)
		return int64(i), err
	default:
		return 0, errors.New(errorVariantTypeNotNumber)
	}
}

func (v Variant) AsInt() (int, error) {
	asInt64, err := v.AsInt64()
	if err != nil {
		return 0, err
	}
	return int(asInt64), nil
}

func (v Variant) AsUInt64() (uint64, error) {
	switch v.variantType {
	case TypeBool, TypeUInt64, TypeEmpty:
		return uint64(v.numberValue), nil
	case TypeInt64:
		if v.numberValue < 0 {
			return 0, errors.New(errorVariantValueOverFlow)
		}
		return uint64(v.numberValue), nil
	case TypeFloat64:
		f := *(*float64)(unsafe.Pointer(&v.numberValue))
		if f < 0 || f > float64(^uint64(0)) {
			return 0, errors.New(errorVariantValueOverFlow)
		}
		return uint64(f), nil
	case TypeString:
		d, err := strconv.ParseInt(v.complexValue.(string), 10, 64)
		return uint64(d), err
	default:
		return 0, errors.New(errorVariantTypeNotNumber)
	}
}

func (v Variant) AsFloat32() (float32, error) {
	switch v.variantType {
	case TypeEmpty:
		return 0, nil
	case TypeBool, TypeInt64, TypeUInt64:
		return float32(v.numberValue), nil
	case TypeFloat64:
		return float32(*(*float64)(unsafe.Pointer(&v.numberValue))), nil
	case TypeString:
		f, err := strconv.ParseFloat(v.complexValue.(string), 32)
		return float32(f), err
	default:
		return 0, errors.New(errorVariantTypeNotNumber)
	}
}

func (v Variant) AsFloat64() (float64, error) {
	switch v.variantType {
	case TypeEmpty:
		return 0, nil
	case TypeBool, TypeInt64:
		return float64(v.numberValue), nil
	case TypeUInt64:
		return float64(uint64(v.numberValue)), nil
	case TypeFloat64:
		return *(*float64)(unsafe.Pointer(&v.numberValue)), nil
	case TypeString:
		return strconv.ParseFloat(v.complexValue.(string), 64)
	default:
		return 0, errors.New(errorVariantTypeNotNumber)
	}
}

func (v *Variant) AsInterface() any {
	switch v.variantType {
	case TypeString:
		return v.complexValue
	case TypeFloat64:
		return *(*float64)(unsafe.Pointer(&v.numberValue))
	case TypeBool:
		return v.numberValue != 0
	case TypeInt64, TypeUInt64:
		return v.numberValue
	case TypeList:
		list, ok := v.complexValue.([]Variant)
		if !ok {
			return make([]any, 0)
		}
		result := make([]any, len(list))
		for i := range list {
			result[i] = list[i].AsInterface()
		}
		return result
	case TypeMap:
		sp, _ := v.StructPairs()
		result := make(map[string]any, sp.len())
		for i := range sp.keys {
			result[sp.keys[i]] = sp.vals[i].AsInterface()
		}
		return result
	default:
		return v.complexValue
	}
}

func (v Variant) String() string {
	return v.AsString()
}

func (v Variant) AsString() string {
	switch v.variantType {
	case TypeEmpty:
		return ""
	case TypeBool, TypeInt64, TypeUInt64:
		return strconv.FormatInt(v.numberValue, 10)
	case TypeFloat64:
		return strconv.FormatFloat(*(*float64)(unsafe.Pointer(&v.numberValue)), 'f', -1, 64)
	case TypeString:
		return v.complexValue.(string)
	case TypeList:
		valList := v.complexValue.([]Variant)
		var strList []string
		for i := range valList {
			if valList[i].IsNumber() {
				strList = append(strList, valList[i].AsString())
			} else {
				strList = append(strList, "\""+valList[i].AsString()+"\"")
			}
		}
		return "[" + strings.Join(strList, ",") + "]"
	case TypeMap:
		sp, _ := v.StructPairs()
		var strKeyValues []string
		for i := range sp.keys {
			k := sp.keys[i]
			value := sp.vals[i]
			if value.variantType == TypeString || value.variantType == TypeEmpty {
				strKeyValues = append(strKeyValues, fmt.Sprintf("\"%s\":\"%s\"", k, value.AsString()))
			} else {
				strKeyValues = append(strKeyValues, fmt.Sprintf("\"%s\":%s", k, value.AsString()))
			}
		}
		return "{" + strings.Join(strKeyValues, ",") + "}"
	default:
		return ""
	}
}

func (c *Variant) AddString(value string) *Variant {
	switch c.variantType {
	case TypeEmpty:
		c.complexValue = value
		c.variantType = TypeString
		return c
	case TypeString:
		if c.complexValue == nil {
			c.complexValue = ""
		}
		str, ok := c.complexValue.(string)
		if !ok {
			return c
		}
		c.complexValue = str + value
	}
	return c
}

func (v Variant) IsEqual(r Variant) bool {
	if v.variantType != r.variantType {
		return false
	}
	switch v.variantType {
	case TypeEmpty:
		return true
	case TypeBool, TypeInt64, TypeUInt64:
		if v.numberValue == r.numberValue {
			return true
		}
	case TypeFloat64:
		if IsFloat64Equal(*(*float64)(unsafe.Pointer(&v.numberValue)), *(*float64)(unsafe.Pointer(&r.numberValue))) {
			return true
		}
	case TypeMap:
		bl := true
		v.Range(func(key string, value Variant) bool {
			rv, ok := r.MapGet(key)
			if !ok || !value.IsEqual(rv) {
				bl = false
				return false
			}
			return true
		})
		return bl
	case TypeList:
		if v.Len() != r.Len() {
			return false
		}
		bl := true
		v.RangeByIndex(func(index int, value Variant) bool {
			rv, ok := r.ListGet(index)
			if !ok || !value.IsEqual(rv) {
				bl = false
				return false
			}
			return true
		})
		return bl
	default:
		if v.AsString() == r.AsString() {
			return true
		}
	}
	return false
}
