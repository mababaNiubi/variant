package variant

import (
	"errors"
	"math"
	"strings"
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

// isNumChar reports whether c may appear in a numeric string (integer, decimal,
// or scientific notation). Used as a cheap pre-filter so plain strings (the
// common case) are rejected in a single pass instead of running regexes.
func isNumChar(c byte) bool {
	return (c >= '0' && c <= '9') || c == '+' || c == '-' || c == '.' || c == 'e' || c == 'E'
}

func IsNumber(s string) bool {
	// 去除首尾空格
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return false
	}
	// 快速拒绝：含数字不可能出现的字符则直接判非数字
	for i := 0; i < len(s); i++ {
		if !isNumChar(s[i]) {
			return false
		}
	}
	return isInt(s) || isDec(s) || IsSciNum(s)
}

func GetStringValueType(s string) Type {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return TypeEmpty
	}
	// 快速拒绝：纯字符串（最常见）一次扫描即出
	for i := 0; i < len(s); i++ {
		if !isNumChar(s[i]) {

			return TypeString
		}
	}
	if isInt(s) {
		return TypeInt64
	} else if isDec(s) || IsSciNum(s) {
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

// isDec reports whether s matches a decimal literal: optional sign, optional
// digits, a dot, and at least one digit on either side ("1.5", ".5", "5.").
func isDec(s string) bool {
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	before := 0
	for ; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		before++
	}
	if i >= len(s) || s[i] != '.' {
		return false
	}
	i++
	after := 0
	for ; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		after++
	}
	return (before+after > 0) && i == len(s)
}

// isInt reports whether s matches an integer literal: optional sign and 1+
// digits ("123", "-42", "+7").
func isInt(s string) bool {
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	if i >= len(s) {
		return false
	}
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
