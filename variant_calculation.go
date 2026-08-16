package variant

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

func (v Variant) Reduce(va Variant) (Variant, error) {
	if v.variantType == va.variantType {
		switch v.variantType {
		case TypeInt64, TypeUInt64, TypeBool:
			return NewInt64(v.numberValue - va.numberValue), nil
		case TypeFloat64:
			value := *(*float64)(unsafe.Pointer(&v.numberValue)) - *(*float64)(unsafe.Pointer(&va.numberValue))
			return NewFloat64(value), nil
		case TypeEmpty:
			return va, nil
		default:
			return v, errors.New(errorVariantTypeNotNumber)
		}
	}
	switch v.variantType {
	case TypeFloat64:
		fl, err := va.AsFloat64()
		if err != nil {
			return v, err
		}
		numberValue := *(*float64)(unsafe.Pointer(&v.numberValue))
		return NewFloat64(numberValue - fl), nil
	case TypeUInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 - *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsUInt64()
		if err != nil {
			return v, err
		}
		return NewUInt64(uint64(v.numberValue) - numberValue), nil
	case TypeInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 - *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsInt64()
		if err != nil {
			return v, err
		}
		return NewInt64(v.numberValue - numberValue), nil
	case TypeEmpty:
		return va, nil
	default:
		return v, errors.New(errUnsupportedType)
	}
}

func (v Variant) Increase(va Variant) (Variant, error) {
	if v.variantType == va.variantType {
		switch v.variantType {
		case TypeInt64, TypeUInt64, TypeBool:
			return NewInt64(v.numberValue + va.numberValue), nil
		case TypeFloat64:
			value := *(*float64)(unsafe.Pointer(&v.numberValue)) + *(*float64)(unsafe.Pointer(&va.numberValue))
			return NewFloat64(value), nil
		case TypeList:
			return NewValueList(append(v.complexValue.([]Variant), va.complexValue.([]Variant)...)), nil
		case TypeMap:
			mp, _ := v.mapVariant()
			if vaMap, ok := va.mapVariant(); ok {
				for key, value := range vaMap {
					mp[key] = value
				}
			}
			return NewValueMap(mp), nil
		case TypeString:
			return NewString(v.complexValue.(string) + va.complexValue.(string)), nil
		case TypeEmpty:
			return va, nil
		default:
			return v, errors.New(errUnsupportedType)
		}
	}
	switch v.variantType {
	case TypeString:
		return NewString(v.complexValue.(string) + va.AsString()), nil
	case TypeList:
		valList := v.complexValue.([]Variant)
		return NewValueList(append(valList, va)), nil
	case TypeFloat64:
		fl, err := va.AsFloat64()
		if err != nil {
			return v, err
		}
		numberValue := *(*float64)(unsafe.Pointer(&v.numberValue))
		return NewFloat64(numberValue + fl), nil
	case TypeUInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 + *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsUInt64()
		if err != nil {
			return v, err
		}
		return NewUInt64(uint64(v.numberValue) + numberValue), nil
	case TypeInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 + *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsInt64()
		if err != nil {
			return v, err
		}
		return NewInt64(v.numberValue + numberValue), nil
	case TypeEmpty:
		return va, nil
	default:
		return v, errors.New(errUnsupportedType)
	}
}

func (v Variant) Multiple(va Variant) (Variant, error) {
	if v.variantType == va.variantType {
		switch v.variantType {
		case TypeInt64, TypeUInt64, TypeBool:
			return NewInt64(v.numberValue * va.numberValue), nil
		case TypeFloat64:
			value := *(*float64)(unsafe.Pointer(&v.numberValue)) * *(*float64)(unsafe.Pointer(&va.numberValue))
			return NewFloat64(value), nil
		case TypeEmpty:
			return va, nil
		default:
			return v, errors.New(errorVariantTypeNotNumber)
		}
	}
	switch v.variantType {
	case TypeFloat64:
		fl, err := va.AsFloat64()
		if err != nil {
			return v, err
		}
		numberValue := *(*float64)(unsafe.Pointer(&v.numberValue))
		return NewFloat64(numberValue * fl), nil
	case TypeUInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 * *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsUInt64()
		if err != nil {
			return v, err
		}
		return NewUInt64(uint64(v.numberValue) * numberValue), nil
	case TypeInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 * *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsInt64()
		if err != nil {
			return v, err
		}
		return NewInt64(v.numberValue * numberValue), nil
	case TypeEmpty:
		return va, nil
	default:
		return v, errors.New(errUnsupportedType)
	}
}

func (v Variant) Divide(va Variant) (Variant, error) {
	if v.variantType == va.variantType {
		switch v.variantType {
		case TypeInt64, TypeUInt64, TypeBool:
			return NewInt64(v.numberValue / va.numberValue), nil
		case TypeFloat64:
			value := *(*float64)(unsafe.Pointer(&v.numberValue)) / *(*float64)(unsafe.Pointer(&va.numberValue))
			return NewFloat64(value), nil
		case TypeEmpty:
			return va, nil
		default:
			return v, errors.New(errorVariantTypeNotNumber)
		}
	}
	switch v.variantType {
	case TypeFloat64:
		fl, err := va.AsFloat64()
		if err != nil {
			return v, err
		}
		numberValue := *(*float64)(unsafe.Pointer(&v.numberValue))
		return NewFloat64(numberValue / fl), nil
	case TypeUInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 / *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsUInt64()
		if err != nil {
			return v, err
		}
		return NewUInt64(uint64(v.numberValue) / numberValue), nil
	case TypeInt64:
		if va.Type() == TypeFloat64 {
			asFloat64, err := v.AsFloat64()
			if err != nil {
				return va, err
			}
			return NewFloat64(asFloat64 / *(*float64)(unsafe.Pointer(&va.numberValue))), nil
		}
		numberValue, err := va.AsInt64()
		if err != nil {
			return v, err
		}
		return NewInt64(v.numberValue / numberValue), nil
	case TypeEmpty:
		return va, nil
	default:
		return v, errors.New(errUnsupportedType)
	}
}

func (v Variant) Decimal(accuracy int) Variant {
	switch v.variantType {
	case TypeList:
		valList := v.complexValue.([]Variant)
		for i := range valList {
			valList[i] = valList[i].Decimal(accuracy)
		}
		v.complexValue = valList
	case TypeFloat64:
		*(*float64)(unsafe.Pointer(&v.numberValue)), _ = strconv.ParseFloat(fmt.Sprintf("%.*f", accuracy, *(*float64)(unsafe.Pointer(&v.numberValue))), 64)
	default:
		//不需要操作
	}
	return v
}

func (v Variant) CompareNumberBySymbol(r Variant, symbol string) (bool, error) {
	switch v.variantType {
	case TypeInt64, TypeEmpty, TypeBool:
		switch r.variantType {
		case TypeInt64, TypeBool, TypeEmpty:
			return compareNumberBySymbol(v.numberValue, r.numberValue, symbol)
		case TypeUInt64:
			if v.numberValue < 0 && r.numberValue >= 0 {
				return compareNumberBySymbol(-1, 0, symbol)
			}
			return compareNumberBySymbol(v.numberValue, r.numberValue, symbol)
		case TypeFloat64:
			return compareNumberBySymbol(float64(v.numberValue), *(*float64)(unsafe.Pointer(&r.numberValue)), symbol)
		case TypeString:
			right, err := r.AsInt64()
			if err != nil {
				return false, err
			}
			return compareNumberBySymbol(v.numberValue, right, symbol)
		default:
			return false, errors.New(errorVariantTypeNotNumber)
		}
	case TypeFloat64:
		switch r.variantType {
		case TypeInt64, TypeBool, TypeEmpty:
			return compareNumberBySymbol(*(*float64)(unsafe.Pointer(&v.numberValue)), float64(r.numberValue), symbol)
		case TypeUInt64:
			return compareNumberBySymbol(*(*float64)(unsafe.Pointer(&v.numberValue)), float64(uint64(r.numberValue)), symbol)
		case TypeFloat64:
			return compareNumberBySymbol(*(*float64)(unsafe.Pointer(&v.numberValue)), *(*float64)(unsafe.Pointer(&r.numberValue)), symbol)
		case TypeString:
			right, err := r.AsFloat64()
			if err != nil {
				return false, err
			}
			return compareNumberBySymbol(*(*float64)(unsafe.Pointer(&v.numberValue)), right, symbol)
		default:
			return false, errors.New(errorVariantTypeNotNumber)
		}
	case TypeUInt64:
		switch r.variantType {
		case TypeInt64:
			if r.numberValue < 0 {
				return compareNumberBySymbol(1, 0, symbol)
			}
			return compareNumberBySymbol(uint64(v.numberValue), uint64(r.numberValue), symbol)
		case TypeUInt64, TypeBool, TypeEmpty:
			return compareNumberBySymbol(uint64(v.numberValue), uint64(r.numberValue), symbol)
		case TypeFloat64:
			if r.numberValue < 0 {
				return compareNumberBySymbol(1, 0, symbol)
			}
			return compareNumberBySymbol(float64(uint64(v.numberValue)), *(*float64)(unsafe.Pointer(&r.numberValue)), symbol)
		case TypeString:
			right, err := r.AsFloat64()
			if err != nil {
				return false, err
			}
			return compareNumberBySymbol(float64(v.numberValue), right, symbol)
		default:
			return false, errors.New(errorVariantTypeNotNumber)
		}
	case TypeString:
		right, err := r.AsFloat64()
		if err != nil {
			return false, err
		}
		left, err := v.AsFloat64()
		if err != nil {
			return false, err
		}
		return compareNumberBySymbol(left, right, symbol)
	default:
		return false, errors.New(errorVariantTypeNotNumber)
	}
}

func (v Variant) Comparable(r Variant) bool {
	result, err := v.CompareNumberBySymbol(r, ">")
	if err == nil {
		return result
	}
	return v.AsString() > r.AsString()
}
