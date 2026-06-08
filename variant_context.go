package variant

import (
	"errors"
	"math"
	"regexp"
	"strings"
)

var (
	isIntRe  = regexp.MustCompile(`^[\+-]?\d+$`)
	isDecRe1 = regexp.MustCompile(`^[\+-]?\d*\.\d+$`)
	isDecRe2 = regexp.MustCompile(`^[\+-]?\d+\.\d*$`)
)

type Type int8

const (
	TypeEmpty Type = iota
	TypeBool
	TypeInt64
	TypeUInt64
	TypeFloat64
	TypeString
	TypeList
	TypeMap
)

const (
	errorVariantEmpty         = "variant is nil"
	errorVariantTypeNotNumber = "is not number type"
	errorVariantValueOverFlow = "value overflow"
	errUnsupportedType        = "unsupported type"
)

const (
	EqualSymbol        string = "="
	NotEqualSymbol     string = "!="
	GreaterThanSymbol  string = ">"
	LessThanSymbol     string = "<"
	GreaterEqualSymbol string = ">="
	LessEqualSymbol    string = "<="
)

func compareNumberBySymbol[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float64 | float32](left T, right T, symbol string) (bool, error) {
	switch symbol {
	case EqualSymbol:
		return left == right, nil
	case NotEqualSymbol:
		return left != right, nil
	case GreaterThanSymbol:
		return left > right, nil
	case LessThanSymbol:
		return left < right, nil
	case GreaterEqualSymbol:
		return left >= right, nil
	case LessEqualSymbol:
		return left <= right, nil
	default:
		return false, errors.New("invalid symbol")
	}
}
func IsFloat64Equal(a, b float64) bool {
	return math.Abs(a-b) < 1e-14
}
func IsNumber(s string) bool {
	// 去除首尾空格
	s = strings.TrimSpace(s)
	// 否则判断是否为整数或小数
	return isInt(s) || isDec(s) || IsSciNum(s)
}

func GetStringValueType(s string) Type {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return TypeEmpty
	}
	if isInt(s) {
		return TypeInt64
	} else if IsNumber(s) {
		return TypeFloat64
	}
	return TypeString
}

// 是否为科学计数法
func checkSciParts(num1, num2 string) bool {
	// e 前后字符串长度为0 是错误的
	if len(num1) == 0 || len(num2) == 0 {
		return false
	}
	// e 后面必须是整数，前面可以是整数或小数  4  +
	return (isInt(num1) || isDec(num1)) && isInt(num2)
}

func IsSciNum(s string) bool {
	for i := 0; i < len(s); i++ {
		// 存在 e 或 E, 判断是否为科学计数法
		if s[i] == 'e' || s[i] == 'E' {
			return checkSciParts(s[:i], s[i+1:])
		}
	}
	return false
}

func isDec(s string) bool {
	return isDecRe1.MatchString(s) || isDecRe2.MatchString(s)
}

func isInt(s string) bool {
	return isIntRe.MatchString(s)
}
