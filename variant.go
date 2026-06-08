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
		return decodeMap(val)
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
	return Variant{
		variantType:  TypeMap,
		complexValue: v,
	}
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
		vs := v.complexValue.(map[string]Variant)
		return len(vs) == 0
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
		mp, ok := v.complexValue.(map[string]Variant)
		result := make(map[string]any)
		if !ok {
			return result
		}
		for k, value := range mp {
			result[k] = value.AsInterface()
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
		valMap := v.complexValue.(map[string]Variant)
		var strKeyValues []string
		for k, value := range valMap {
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
