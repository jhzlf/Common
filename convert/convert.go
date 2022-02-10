package convert

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// StringToInt : string -> int
func StringToInt(s string) int {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return TmpInt
}

// StringToInt16 : string -> int16
func StringToInt16(s string) int16 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return int16(TmpInt)
}

// StringToInt32 : string -> int32
func StringToInt32(s string) int32 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return int32(TmpInt)
}

// StringToInt64 : string -> int64
func StringToInt64(s string) int64 {
	if s == "" {
		return 0
	}

	tmp, _ := strconv.ParseInt(s, 10, 64)

	return tmp
}

// StringToUint : string -> uint
func StringToUint(s string) uint {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint(TmpInt)
}

// StringToUint16 : string -> uint16
func StringToUint16(s string) uint16 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint16(TmpInt)
}

// StringToUint32 : string -> uint32
func StringToUint32(s string) uint32 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint32(TmpInt)
}

// StringToUint64 : string -> uint64
func StringToUint64(s string) uint64 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint64(TmpInt)
}

// StringToFloat64 : string -> float64
func StringToFloat64(s string) float64 {
	if s == "" {
		return 0
	}

	tmp, _ := strconv.ParseFloat(s, 64)

	return tmp
}

// Int64toString : int64 -> string
func Int64toString(a int64) string {
	return strconv.FormatInt(a, 10)
}

// Int32toString : int32 -> string
func Int32toString(a int32) string {
	return strconv.FormatInt(int64(a), 10)
}

// Int16toString : int16 -> string
func Int16toString(a int16) string {
	return strconv.FormatInt(int64(a), 10)
}

// InttoString : int -> string
func InttoString(a int) string {
	return strconv.Itoa(a)
}

// Uint64toString : uint64 -> string
func Uint64toString(a uint64) string {
	return strconv.FormatUint(a, 10)
}

// Uint32toString : uint32 -> string
func Uint32toString(a uint32) string {
	return strconv.FormatUint(uint64(a), 10)
}

// Uint16toString : uint16 -> string
func Uint16toString(a uint16) string {
	return strconv.FormatUint(uint64(a), 10)
}

// Float64toString : float64 -> string
func Float64toString(a float64) string {
	return strconv.FormatFloat(a, 'f', -1, 64)
}

// StringToInt32Slice :
func StringToInt32Slice(s string, sep string) (ret []int32) {
	tokens := strings.Split(s, sep)
	for _, k := range tokens {
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil
		}
		ret = append(ret, int32(i))
	}
	return
}

// BytesToString :
func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

// ToString
func ToString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case int:
		return InttoString(x)
	case int64:
		return Int64toString(x)
	default:
		return fmt.Sprint(v)
	}
}

func StringToSliceByte(s string) []byte {
	l := len(s)
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
		Len:  l,
		Cap:  l,
	}))
}
