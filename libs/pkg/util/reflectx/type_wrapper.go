package reflectx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type TypeWrapper interface {
	GetElem(key interface{}) (TypeWrapper, error)
	GetType() reflect.Type
	StructIndex() []int
	Valid() bool
}

func GenerateTypeWrapper(t reflect.Type, tagName string) TypeWrapper {
	m := make(map[string]TypeWrapper)
	return generateTypeWrapper(t, m, tagName, nil)
}

func generateTypeWrapper(t reflect.Type, used map[string]TypeWrapper, tagName string, index []int) TypeWrapper {
	var out TypeWrapper
	switch t.Kind() {
	case reflect.Int, reflect.Bool, reflect.String,
		reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint16,
		reflect.Uint8, reflect.Uint64:
		out = &commonType{
			t:     t,
			index: index,
		}
	case reflect.Array, reflect.Slice, reflect.Map:
		out = &mapSliceType{
			elemWrapper: generateTypeWrapper(t.Elem(), used, tagName, index),
			t:           t,
			index:       index,
		}
	case reflect.Struct:
		m := make(map[string]TypeWrapper)
		key := t.PkgPath() + ":" + t.Name()
		if o, ok := used[key]; ok {
			return o
		}
		out = &structType{
			t:      t,
			cached: m,
			index:  index,
		}
		used[key] = out
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			m[strings.Split(field.Tag.Get(tagName), ",")[0]] =
				generateTypeWrapper(field.Type, used, tagName, field.Index)
		}
	default:
		out = &invalidType{
			index: index,
		}
	}
	return out
}

type structType struct {
	cached map[string]TypeWrapper
	t      reflect.Type
	index  []int
}

func (s *structType) GetType() reflect.Type {
	return s.t
}

func (s *structType) GetElem(key interface{}) (TypeWrapper, error) {
	k := key.(string)
	v, ok := s.cached[k]
	if !ok {
		return nil, fmt.Errorf("elem not found in struct, key: %s", k)
	}
	return v, nil
}

func (s *structType) StructIndex() []int {
	return s.index
}

func (s *structType) Valid() bool {
	return true
}

type mapSliceType struct {
	t           reflect.Type
	elemWrapper TypeWrapper
	index       []int
}

func (s *mapSliceType) GetType() reflect.Type {
	return s.t
}

func (s *mapSliceType) GetElem(key interface{}) (TypeWrapper, error) {
	return s.elemWrapper, nil
}

func (s *mapSliceType) StructIndex() []int {
	return s.index
}

func (s *mapSliceType) Valid() bool {
	return true
}

type commonType struct {
	t     reflect.Type
	index []int
}

func (s *commonType) GetType() reflect.Type {
	return s.t
}

func (s *commonType) GetElem(key interface{}) (TypeWrapper, error) {
	return nil, errors.New("common type can't get field")
}

func (s *commonType) Valid() bool {
	return true
}

func (s *commonType) StructIndex() []int {
	return s.index
}

type invalidType struct {
	index []int
}

func (s *invalidType) GetType() reflect.Type {
	panic("invalidType GetType")
}

func (s *invalidType) GetElem(key interface{}) (TypeWrapper, error) {
	panic("invalidType GetElem")
}

func (s *invalidType) Valid() bool {
	return false
}

func (s *invalidType) StructIndex() []int {
	return s.index
}
