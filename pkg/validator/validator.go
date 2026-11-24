package validator

import (
	"math"
	"net/mail"
	"reflect"
	"strings"
	"time"
)

func IsEmail(value any) bool {
	text, ok := asString(value)
	if !ok || text == "" || text != strings.TrimSpace(text) {
		return false
	}

	address, err := mail.ParseAddress(text)
	if err != nil || address.Address != text {
		return false
	}

	at := strings.LastIndexByte(text, '@')
	if at <= 0 || at == len(text)-1 {
		return false
	}
	domain := text[at+1:]
	return !strings.HasPrefix(domain, ".") &&
		!strings.HasSuffix(domain, ".") &&
		strings.Contains(domain, ".")
}

func IsString(value any) bool {
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Kind() == reflect.String
}

func IsBlankString(value any) bool {
	text, ok := asString(value)
	return ok && text != ""
}

func MinLength(value any, size int) bool {
	text, ok := asString(value)

	return ok && text != "" && len(text) >= size
}

func MaxLength(value any, size int) bool {
	text, ok := asString(value)
	return ok && text != "" && len(text) <= size
}

func IsNumber(value any) bool {
	if value == nil {
		return false
	}

	number := reflect.ValueOf(value)
	switch number.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true
	case reflect.Float32, reflect.Float64:
		f := number.Float()
		return !math.IsNaN(f) && !math.IsInf(f, 0)
	default:
		return false
	}
}

func MinNum(value any, minimum any) bool {
	v, vOK := numberAsFloat64(value)
	minNum, minOK := numberAsFloat64(minimum)
	return vOK && minOK && v >= minNum
}

func MaxNum(value any, maximum any) bool {
	v, vOK := numberAsFloat64(value)
	maxNum, maxOK := numberAsFloat64(maximum)
	return vOK && maxOK && v <= maxNum
}

func IsBoolean(value any) bool {
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Kind() == reflect.Bool
}

func IsTime(value any) bool {
	if value == nil {
		return false
	}

	reflectedValue := reflect.ValueOf(value)
	timeType := reflect.TypeOf(time.Time{})
	for reflectedValue.Kind() == reflect.Pointer {
		if reflectedValue.IsNil() {
			return false
		}
		reflectedValue = reflectedValue.Elem()
	}

	return reflectedValue.Type() == timeType
}

func IsValidDate(s time.Time, layout string) bool {
	if layout == "" || s.IsZero() {
		return false
	}
	return s.Format(layout) != ""
}

func IsArray(value any) bool {
	if value == nil {
		return false
	}
	kind := reflect.ValueOf(value).Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

func asString(value any) (string, bool) {
	if value == nil {
		return "", false
	}

	reflected := reflect.ValueOf(value)
	if reflected.Kind() != reflect.String {
		return "", false
	}
	return reflected.String(), true
}

func numberAsFloat64(value any) (float64, bool) {
	if value == nil {
		return 0, false
	}

	number := reflect.ValueOf(value)
	switch number.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(number.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(number.Uint()), true
	case reflect.Float32, reflect.Float64:
		result := number.Float()
		if math.IsNaN(result) || math.IsInf(result, 0) {
			return 0, false
		}
		return result, true
	default:
		return 0, false
	}
}
