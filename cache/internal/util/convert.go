package util

import (
	"strconv"
	"reflect"
	"fmt"
)

func Atoi(b []byte) (int, error) {
	return strconv.Atoi(BytesToString(b))
}

func ParseInt(b []byte, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(BytesToString(b), base, bitSize)
}

func ParseUint(b []byte, base int, bitSize int) (uint64, error) {
	return strconv.ParseUint(BytesToString(b), base, bitSize)
}

func ParseFloat(b []byte, bitSize int) (float64, error) {
	return strconv.ParseFloat(BytesToString(b), bitSize)
}

func AsString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}