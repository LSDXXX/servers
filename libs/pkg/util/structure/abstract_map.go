package structure

import (
	"encoding/json"
	"reflect"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

// ToStruct map to struct
// use json tag
func ToStruct(src interface{}, dst interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      dst,
		TagName:     "json",
		ErrorUnused: true,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(src)
}

// AbstractMap map[string]interface{}
type AbstractMap map[string]interface{}

// NewAbstractMap create abstract map
func NewAbstractMap() AbstractMap {
	return make(map[string]interface{})
}

// CopyMap description
// @param src
// @return map
func CopyMap(src interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	mapValue := reflect.ValueOf(src)
	iter := mapValue.MapRange()
	for iter.Next() {
		k := iter.Key().Interface().(string)
		out[k] = iter.Value().Interface()
	}
	return out
}

// ToAbstractMap cast to abstract map
func ToAbstractMap(val interface{}) (m AbstractMap) {
	defer func() {
		if r := recover(); r != nil {
			m = make(AbstractMap)
		}
	}()
	switch val := val.(type) {
	case AbstractMap:
		return val
	case map[string]interface{}:
		return AbstractMap(val)
	case map[string]string:
		return CopyMap(val)
	case map[string][]string:
		return CopyMap(val)
	case map[string][]interface{}:
		return CopyMap(val)
	case map[string]int:
		return CopyMap(val)
	case map[string][]int:
		out := make(map[string]interface{})
		mapValue := reflect.ValueOf(val)
		iter := mapValue.MapRange()
		for iter.Next() {
			k := iter.Key().Interface().(string)
			out[k] = iter.Value()
		}
		return out
	case nil:
		// 避免cast 失败panic
		return make(map[string]interface{})
	default:
		s := structs.New(val)
		s.TagName = "json"
		return s.Map()
	}
}

// SetValue set value
func (m AbstractMap) SetValue(key string, val interface{}) {
	m[key] = val
}

// GetValue get value
func (m AbstractMap) GetValue(key string) interface{} {
	return m[key]
}

// GetMap get map
func (m AbstractMap) GetMap(key string) AbstractMap {
	return ToAbstractMap(m[key])
}

// SetMap set map
func (m AbstractMap) SetMap(key string, val map[string]interface{}) {
	m[key] = val
}

// GetString get string value
func (m AbstractMap) GetString(key string) string {
	return cast.ToString(m[key])
}

// SetString set string value
func (m AbstractMap) SetString(key, val string) {
	m.SetValue(key, val)
}

// GetBool get bool value
func (m AbstractMap) GetBool(key string) bool {
	return cast.ToBool(m[key])
}

// SetBool set bool value
func (m AbstractMap) SetBool(key string, val bool) {
	m.SetValue(key, val)
}

// GetInt get int value
func (m AbstractMap) GetInt(key string) int {
	return cast.ToInt(m[key])
}

// SetInt set int value
func (m AbstractMap) SetInt(key string, val int) {
	m.SetValue(key, val)
}

// GetInt32 get int32 value
func (m AbstractMap) GetInt32(key string) int32 {
	return cast.ToInt32(m[key])
}

// SetInt32 set int32 value
func (m AbstractMap) SetInt32(key string, val int32) {
	m.SetValue(key, val)
}

// GetInt64 get int64 value
func (m AbstractMap) GetInt64(key string) int64 {
	return cast.ToInt64(m[key])
}

// SetInt64 set int64 value
func (m AbstractMap) SetInt64(key string, val int64) {
	m.SetValue(key, val)
}

// GetUint32 get uint32 value
func (m AbstractMap) GetUint32(key string) uint32 {
	return cast.ToUint32(m[key])
}

// SetUint32 set uint32 value
func (m AbstractMap) SetUint32(key string, val uint32) {
	m.SetValue(key, val)
}

// GetUint64 get uint64 value
func (m AbstractMap) GetUint64(key string) uint64 {
	return cast.ToUint64(m[key])
}

// SetUint64 set uint64 value
func (m AbstractMap) SetUint64(key string, val uint64) {
	m.SetValue(key, val)
}

// GetUint get uint value
func (m AbstractMap) GetUint(key string) uint {
	return cast.ToUint(m[key])
}

// SetUint set uint value
func (m AbstractMap) SetUint(key string, val uint) {
	m.SetValue(key, val)
}

// GetFloat64 get float64
func (m AbstractMap) GetFloat64(key string) float64 {
	return cast.ToFloat64(m[key])
}

// SetFloat64 set float64
func (m AbstractMap) SetFloat64(key string, val float64) {
	m[key] = val
}

// GetFloat32 get float32
func (m AbstractMap) GetFloat32(key string) float32 {
	return cast.ToFloat32(m[key])
}

// SetFloat32 set float32
func (m AbstractMap) SetFloat32(key string, val float32) {
	m[key] = val
}

// GetStringSlice get slice
func (m AbstractMap) GetStringSlice(key string) []string {
	return cast.ToStringSlice(m[key])
}

// GetIntSlice get slice
func (m AbstractMap) GetIntSlice(key string) []int {
	return cast.ToIntSlice(m[key])
}

// GetBoolSlice get slice
func (m AbstractMap) GetBoolSlice(key string) []bool {
	return cast.ToBoolSlice(m[key])
}

// GetAbstractSlice get slice
func (m AbstractMap) GetAbstractSlice(key string) AbstractSlice {
	return ToAbstractSlice(m[key])
}

// SetAbstractSlice set slice
func (m AbstractMap) SetAbstractSlice(key string, val AbstractSlice) {
	m[key] = val
}

// GetJSON get object json
func (m AbstractMap) GetJSON(key string) []byte {
	if v := m.GetValue(key); v != nil {
		data, err := json.Marshal(v)
		if err == nil {
			return data
		}
	}
	return nil
}

// TODO slice setter
