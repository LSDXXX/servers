package structure

import (
	"reflect"

	"github.com/spf13/cast"
)

// AbstractSlice []interface{} 方便使用
type AbstractSlice []interface{}

// NewAbstractSlice create abstract slice
func NewAbstractSlice() AbstractSlice {
	return make([]interface{}, 0)
}

func copySlice(src interface{}) []interface{} {
	s := reflect.ValueOf(src)
	var out []interface{}
	for i := 0; i < s.Len(); i++ {
		out = append(out, s.Index(i).Interface())
	}
	return out
}

// ToAbstractSlice cast to abstractslice
func ToAbstractSlice(val interface{}) AbstractSlice {
	switch val := val.(type) {
	case []interface{}:
		return AbstractSlice(val)
	case *AbstractSlice:
		return *val
	case AbstractSlice:
		return val
	case []int:
		return copySlice(val)
	case []string:
		return copySlice(val)
	case []AbstractMap:
		return copySlice(val)
	case []map[string]interface{}:
		return copySlice(val)
	}
	return AbstractSlice{}
}

// Len length
func (s AbstractSlice) Len() int {
	return len(s)
}

// Value raw value
func (s AbstractSlice) Value(i int) interface{} {
	return s[i]
}

// Int int
func (s AbstractSlice) Int(i int) int {
	return cast.ToInt(s[i])
}

// Int32 int32
func (s AbstractSlice) Int32(i int) int32 {
	return cast.ToInt32(s[i])
}

// Int64 int64
func (s AbstractSlice) Int64(i int) int64 {
	return cast.ToInt64(s[i])
}

// Uint uint
func (s AbstractSlice) Uint(i int) uint {
	return cast.ToUint(s[i])
}

// Uint32 uint32
func (s AbstractSlice) Uint32(i int) uint32 {
	return cast.ToUint32(s[i])
}

// Uint64 uint64
func (s AbstractSlice) Uint64(i int) uint64 {
	return cast.ToUint64(s[i])
}

// String string
func (s AbstractSlice) String(i int) string {
	return cast.ToString(s[i])
}

// Slice slice
func (s AbstractSlice) Slice(i int) AbstractSlice {
	return cast.ToSlice(s[i])
}

// Bool bool
func (s AbstractSlice) Bool(i int) bool {
	return cast.ToBool(s[i])
}

// Map map
func (s AbstractSlice) Map(i int) AbstractMap {
	return ToAbstractMap(s[i])
}

// StringSlice []string
func (s AbstractSlice) StringSlice(i int) []string {
	return cast.ToStringSlice(s[i])
}

// IntSlice []int
func (s AbstractSlice) IntSlice(i int) []int {
	return cast.ToIntSlice(s[i])
}
