package reflectx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/LSDXXX/libs/pkg/mapstructure"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/spf13/cast"
)

// ValueOf
//  @param obj
//  @return *reflect.Value
func ValueOf(obj interface{}) reflect.Value {
	var out reflect.Value
	switch v := obj.(type) {
	case *reflect.Value:
		out = *v
	case reflect.Value:
		out = v
	default:
		out = reflect.ValueOf(obj)
	}
	return out
}

// FieldAttr field
type FieldAttr struct {
	isIndex   bool
	fieldName string
	index     string
}

// ParseExpr parse
//  @param attr
//  @param split
//  @return []FieldAttr
//  @return error
func ParseExpr(attr string, split bool) ([]FieldAttr, error) {
	fields := strings.Split(attr, ".")
	var out []FieldAttr
	for _, field := range fields {
		fieldAndIndex := strings.Split(field, "[")
		for _, fieldOrIndex := range fieldAndIndex {
			right := strings.Index(fieldOrIndex, "]")
			if right != -1 {
				index := fieldOrIndex[:right]
				out = append(out, FieldAttr{
					isIndex: true,
					index:   strings.Trim(index, " "),
				})
			} else {
				out = append(out, FieldAttr{
					fieldName: strings.Trim(fieldOrIndex, " "),
				})
			}
		}
	}
	return out, nil
}

// IsValid
//  @param val
//  @return bool
func IsValid(val reflect.Value) bool {
	return val.Kind() != reflect.Invalid
}

func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return v.Elem()
	}
	return v
}

func clone(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}
	switch v.Kind() {
	case reflect.Pointer:
		// 暂时不支持除struct以外的类型
		out := reflect.New(v.Type().Elem())
		v = clone(v.Elem())
		if !v.IsValid() {
			return out
		}
		out.Elem().Set(v)
		return out
	case reflect.Slice:
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		reflect.Copy(out, v)
		return out
	case reflect.Map:
		out := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			out.SetMapIndex(key, v.MapIndex(key))
		}
		return out
	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			newValField := out.Field(i)
			if newValField.CanSet() {
				newValField.Set(v.Field(i))
			}
		}
		return out
	default:
		out := reflect.New(v.Type()).Elem()
		out.Set(v)
		return out
	}
}

func createFromType(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0)
	case reflect.Map:
		return reflect.MakeMap(t)
	case reflect.Pointer:
		v := reflect.New(t.Elem())
		v.Elem().Set(createFromType(t.Elem()))
		return v
	default:
		return reflect.New(t).Elem()
	}
}

func fieldSetter(o reflect.Value) func(interface{}) reflect.Value {
	return func(val interface{}) reflect.Value {
		v := ValueOf(val)
		o.Set(v)
		return o
	}
}

func getAssignableToSimple(target reflect.Value,
	output reflect.Value) (reflect.Value, bool) {
	if output.Kind() == reflect.Interface {
		return target, true
	}
	if target.Type().AssignableTo(output.Type()) {
		return target, true
	}
	if indir := reflect.Indirect(target); indir.Type().AssignableTo(output.Type()) {
		return indir, true
	}
	return reflect.Value{}, false
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

func doSet(fields []FieldAttr, wrapper TypeWrapper, cur int,
	obj reflect.Value, target reflect.Value, tagNames ...string) (reflect.Value, error) {
	if obj.Kind() == reflect.Pointer &&
		reflect.Indirect(obj).Kind() == reflect.Interface {
		obj = reflect.Indirect(obj)
	}
	if !IsValid(obj) || isNil(obj) {
		obj = createFromType(obj.Type())
	}
	if len(fields) == cur {
		// obj.Type().AssignableTo()
		// TODO: do set
		if out, ok := getAssignableToSimple(target, obj); ok {
			return out, nil
		}
		if obj.Kind() == reflect.Slice {
			obj = reflect.Indirect(reflect.New(obj.Type()))
		}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           obj,
			TagName:          tagNames[0],
			WeaklyTypedInput: true,
		})
		if err != nil {
			return reflect.Value{}, err
		}
		err = decoder.Decode(target)
		return obj, err
	}
	field := fields[cur]
	// obj = reflect.Indirect(obj)
	switch obj.Kind() {
	case reflect.Pointer:
		out, err := doSet(fields, nil, cur, obj.Elem(), target, tagNames...)
		if err != nil {
			return reflect.Value{}, err
		}
		obj.Elem().Set(out)
		return obj, nil
	case reflect.Interface:
		if IsValid(obj.Elem()) {
			return doSet(fields, nil, cur, obj.Elem(),
				target, tagNames...)
		}
		if field.isIndex {
			n := cast.ToInt(field.index)
			obj = reflect.MakeSlice(reflect.TypeOf([]interface{}{}),
				n+1, n+1)
		} else {
			obj = reflect.MakeMap(reflect.TypeOf(map[string]interface{}{}))
		}
		return doSet(fields, nil, cur, obj, target, tagNames...)
	case reflect.Map:
		var elem reflect.Value
		var key interface{}
		if field.isIndex {
			key = field.index
			elem, _, _ = getIndex(obj, ValueOf(field.index))
		} else {
			key = field.fieldName
			elem, _, _ = getAttr(obj, wrapper, field.fieldName, tagNames...)
		}
		if !elem.IsValid() {
			elem = createFromType(obj.Type().Elem())
		} else {
			elem = clone(elem)
		}
		var nextWrapper TypeWrapper
		if wrapper != nil && wrapper.Valid() {
			tmp, err := wrapper.GetElem(nil)
			if err == nil {
				nextWrapper = tmp
			}
		}
		next, err := doSet(fields, nextWrapper, cur+1, elem, target, tagNames...)
		if err != nil {
			return reflect.Value{}, err
		}
		obj.SetMapIndex(reflect.ValueOf(key), next)
		return obj, nil
	case reflect.Slice:
		index, err := strconv.Atoi(field.index)
		if err != nil || !field.isIndex {
			return reflect.Value{}, fmt.Errorf("error at: %s", field.fieldName)
		}
		if obj.Len() <= index {
			if obj.Cap() > index {
				obj.SetLen(index + 1)
			} else {
				tmp := reflect.MakeSlice(obj.Type(), index+1, index*2+1)
				reflect.Copy(tmp, obj)
				obj = tmp
			}
		}
		v, setter, err := getIndex(obj, reflect.ValueOf(index))
		if err != nil {
			return v, err
		}
		var nextWrapper TypeWrapper
		if wrapper != nil && wrapper.Valid() {
			tmp, err := wrapper.GetElem(nil)
			if err == nil {
				nextWrapper = tmp
			}
		}
		v, err = doSet(fields, nextWrapper, cur+1, v, target, tagNames...)
		if err != nil {
			return v, err
		}
		setter(v)
		return obj, nil
	case reflect.Struct:
		v, setter, err := getAttr(obj, wrapper, field.fieldName, tagNames...)
		if err != nil {
			return reflect.Value{}, err
		}
		v = clone(v)
		var nextWrapper TypeWrapper
		if wrapper != nil && wrapper.Valid() {
			tmp, err := wrapper.GetElem(field.fieldName)
			if err == nil {
				nextWrapper = tmp
			}
		}
		v, err = doSet(fields, nextWrapper, cur+1, v, target, tagNames...)
		if err != nil {
			return v, err
		}
		setter(v)
		return obj, nil
	default:
		return reflect.Value{}, fmt.Errorf("error at: %s", field.fieldName)
	}
}

/*
func doSetAttr(expr ast.Expr, root, target reflect.Value,
	tagNames ...string) (reflect.Value, error) {

	var walk func(expr ast.Expr, obj, target reflect.Value) (reflect.Value, error)
	walk = func(expr ast.Expr, obj, target reflect.Value) (reflect.Value, error) {
		switch expr := expr.(type) {
		case *ast.IndexExpr:

		}
	}
}
*/

func SetAttrWithWrapper(obj interface{}, wrapper TypeWrapper, attr string,
	val interface{}, tagNames ...string) error {
	fields, err := ParseExpr(attr, false)
	if err != nil {
		return err
	}
	return SetAttrParsed(obj, wrapper, fields, val, tagNames...)
}

func SetAttr(obj interface{}, attr string,
	val interface{}, tagNames ...string) error {
	return SetAttrWithWrapper(obj, nil, attr, val, tagNames...)
}

func SetAttrParsed(obj interface{}, wrapper TypeWrapper, fields []FieldAttr, val interface{}, tagNames ...string) error {
	if len(tagNames) == 0 {
		tagNames = append(tagNames, "json")
	}
	o := deref(ValueOf(obj))
	v := ValueOf(val)
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	_, err := doSet(fields, wrapper, 0, o, v, tagNames...)
	return err
}

func getAttr(o reflect.Value, wrapper TypeWrapper, name string, tagNames ...string) (reflect.Value, func(interface{}) reflect.Value, error) {
	names := []string{name}
	if name == "Case" {
		names = append(names, "case")
	}
	if o.Kind() == reflect.Map {
		for _, name := range names {
			key := reflect.ValueOf(name)
			val := o.MapIndex(key)
			if val.Kind() != reflect.Invalid {
				return val, nil, nil
			}
		}
		return reflect.Value{}, nil, fmt.Errorf("can't find key: %v in map", name)
	}
	if o.Kind() == reflect.Interface || o.Kind() == reflect.Ptr {
		tmp := o.Elem()
		return getAttr(tmp, wrapper, name)
	}
	if o.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("object kind must be map or struct, not %v", o.Kind().String())
	}
	//val := o.FieldByName(name)

	//if IsValid(val) {
	//return val, fieldSetter(val), nil
	//}

	var val reflect.Value
	if wrapper != nil && wrapper.Valid() {
		for _, name := range names {
			w, err := wrapper.GetElem(name)
			if err == nil {
				val = o.FieldByIndex(w.StructIndex())
			}
		}
	} else {
		valueType := o.Type()
		for i := 0; i < valueType.NumField(); i++ {
			field := valueType.Field(i)
			found := false
			for _, name := range names {
				for _, tag := range tagNames {
					t := field.Tag.Get(tag)
					if t == name || strings.HasPrefix(t, name+",") {
						val = o.FieldByIndex(field.Index)
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
	}
	if IsValid(val) {
		return val, fieldSetter(val), nil
	}
	return reflect.Value{}, nil, fmt.Errorf("can't find attr %s", name)
}

func getIndex(o reflect.Value, index reflect.Value) (reflect.Value, func(interface{}) reflect.Value, error) {
	if o.Kind() == reflect.Interface || o.Kind() == reflect.Ptr {
		tmp := o.Elem()
		return getIndex(tmp, index)
	}
	switch o.Kind() {
	case reflect.Array, reflect.Slice:
		i, err := cast.ToIntE(index.Interface())
		if err != nil {
			return reflect.Value{}, nil, fmt.Errorf("cast slice or array index to int err: %w", err)
		}
		if i < 0 || i >= o.Len() {
			return reflect.Value{}, nil, fmt.Errorf("slice or array index is out of range: %d", i)
		}
		tmp := o.Index(i)
		return tmp, fieldSetter(tmp), nil
	case reflect.Map:
		if index.Kind() == reflect.Interface {
			tmp := index.Elem()
			return getIndex(o, tmp)
		}
		if index.Kind() == reflect.String {
			index = reflect.ValueOf(strings.Trim(index.String(), " \""))
		}
		tmp := o.MapIndex(index)
		return tmp, nil, nil
	}
	return reflect.Value{}, nil, errors.New("index object kind is not in (map, array, slice)")
}

func GetAttr(obj interface{}, attr string, tagNames ...string) (reflect.Value, func(interface{}) reflect.Value, error) {
	attr = strings.Replace(attr, "case", "Case", -1)
	// expr, err := parser.ParseExpr(attr)
	fields, err := ParseExpr(attr, false)
	if err != nil {
		return reflect.Value{}, nil, err
	}
	return GetAttrParsed(obj, fields, tagNames...)
}

// ReflectGetValue description
// @param obj
// @param attr
// @return *reflect.Value
// @return func(interface{})
// @return error
func GetAttrParsed(obj interface{}, fields []FieldAttr, tagNames ...string) (reflect.Value, func(interface{}) reflect.Value, error) {

	if len(tagNames) == 0 || util.SliceIndex(tagNames, "json") == -1 {
		tagNames = append(tagNames, "json")
	}
	root := ValueOf(obj)
	// TODO 拆出去

	v := root
	var err error
	var setter func(interface{}) reflect.Value
	for _, field := range fields {
		if field.isIndex {
			tmp := reflect.ValueOf(field.index)
			v, setter, err = getIndex(v, tmp)
		} else {
			v, setter, err = getAttr(v, nil, field.fieldName, tagNames...)
		}
		if err != nil {
			break
		}
	}
	return v, setter, err

	//invalidErr := errors.New("invalid stmt")
	/*
		var walk func(node ast.Expr) (reflect.Value, func(interface{}) reflect.Value, error)
		walk = func(node ast.Expr) (reflect.Value, func(interface{}) reflect.Value, error) {
			switch node.(type) {
			case *ast.IndexExpr:
				indexExpr := node.(*ast.IndexExpr)
				index, _, err := walk(indexExpr.Index)
				if err != nil {
					return reflect.Value{}, nil, err
				}
				val, _, err := walk(indexExpr.X)
				if err != nil {
					return reflect.Value{}, nil, err
				}
				val, setter, err := getIndex(val, index)
				if err != nil {
					return reflect.Value{}, nil, fmt.Errorf("%w at offset %v", err, indexExpr.Pos())
				}
				return val, setter, nil
			case *ast.SelectorExpr:
				selectorExpr := node.(*ast.SelectorExpr)
				n := selectorExpr.Sel.Name
				val, _, err := walk(selectorExpr.X)
				if err != nil {
					return reflect.Value{}, nil, err
				}
				val, setter, err := getAttr(val, n)
				if err != nil {
					return reflect.Value{}, nil, fmt.Errorf("%w at offset %d", err, selectorExpr.Sel.Pos())
				}
				return val, setter, nil
			case *ast.BasicLit:
				bascLit := node.(*ast.BasicLit)
				switch bascLit.Kind {
				case token.INT:
					tmp := reflect.ValueOf(cast.ToInt(bascLit.Value))
					return tmp, nil, nil
				case token.STRING:
					tmp := reflect.ValueOf(bascLit.Value)
					return tmp, nil, nil
				}
				return reflect.Value{}, nil, fmt.Errorf("unknown basicLit: %s", bascLit.Kind.String())
			case *ast.Ident:
				identExpr := node.(*ast.Ident)
				val, setter, err := getAttr(root, identExpr.String())
				if err != nil {
					return reflect.Value{}, nil, fmt.Errorf("%w at offset %d", err, identExpr.Pos())
				}
				return val, setter, nil
			}
			return reflect.Value{}, nil, fmt.Errorf("unknown node type: %s at offset %d", reflect.TypeOf(node).String(), node.Pos())
		}
		return walk(expr)
	*/
}
