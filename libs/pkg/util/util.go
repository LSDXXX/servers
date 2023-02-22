package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/LSDXXX/libs/pkg/log"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unsafe"

	"github.com/dengsgo/math-engine/engine"
	"github.com/go-zookeeper/zk"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

var (
	localIP string

	// ErrStandardVarNotMatched .
	ErrStandardVarNotMatched = errors.New("standard var not matched")

	ParamRegex = regexp.MustCompile(`\${([a-z|A-Z|0-9|\.|_|\[|\]]*)}`)
)

// GetLocalIP 获取本机ip
func GetLocalIP() string {
	if len(localIP) > 0 {
		return localIP
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				return localIP
			}
		}
	}

	return "127.0.0.1"
}

// EncodeURL add query
func EncodeURL(rawURL string, querys map[string]string) (*url.URL, error) {
	surl, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := surl.Query()
	for k, v := range querys {
		query.Set(k, v)
	}
	surl.RawQuery = query.Encode()
	return surl, nil
}

// UnsafeString return string without allocation
func UnsafeString(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

// UnsafeBytes return bytes without allocation
func UnsafeBytes(s string) (b []byte) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Cap = pstring.Len
	pbytes.Len = pstring.Len
	return
}

// PrintStack print stack
func PrintStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

// ParseAndExec description
// @param root
// @param statement
// @return out
// @return err
func ParseAndExec(root interface{}, statement string) (out string, err error) {
	statement, err = FormatObject(statement, root)
	if err != nil {
		return
	}
	r, err := engine.ParseAndExec(statement)
	if err != nil {
		return
	}
	return cast.ToString(r), nil
}

func reflectValue(obj interface{}) *reflect.Value {
	var out *reflect.Value
	switch v := obj.(type) {
	case *reflect.Value:
		out = v
	case reflect.Value:
		out = &v
	default:
		tmp := reflect.ValueOf(obj)
		out = &tmp
	}
	return out
}

type fieldAttr struct {
	isIndex   bool
	fieldName string
	index     string
}

func parseExpr(attr string) ([]fieldAttr, error) {
	fields := strings.Split(attr, ".")
	var out []fieldAttr
	for _, field := range fields {
		left := strings.Index(field, "[")
		if left != -1 {
			right := strings.Index(field, "]")
			if right == -1 {
				return nil, errors.New("parse expr bad index")
			}
			index := field[left+1 : right]
			out = append(out, fieldAttr{
				isIndex:   true,
				fieldName: field[:left],
				index:     index,
			})
		} else {
			out = append(out, fieldAttr{
				fieldName: field,
			})
		}
	}
	return out, nil
}

// ReflectGetValue description
// @param obj
// @param attr
// @return *reflect.Value
// @return func(interface{})
// @return error
func ReflectGetValue(obj interface{}, attr string, tagNames ...string) (*reflect.Value, func(interface{}), error) {
	attr = strings.Replace(attr, "case", "Case", -1)
	// expr, err := parser.ParseExpr(attr)
	fields, err := parseExpr(attr)
	if err != nil {
		return nil, nil, err
	}
	if len(tagNames) == 0 {
		tagNames = append(tagNames, "json")
	}
	root := reflectValue(obj)
	// TODO 拆出去
	isValid := func(val *reflect.Value) bool {
		return val != nil && val.Kind() != reflect.Invalid
	}
	fieldSetter := func(o *reflect.Value) func(interface{}) {
		return func(val interface{}) {
			v := reflectValue(val)
			o.Set(*v)
		}
	}
	mapSetter := func(o *reflect.Value, key *reflect.Value) func(interface{}) {
		return func(val interface{}) {
			v := reflectValue(val)
			o.SetMapIndex(*key, *v)
		}
	}
	var getAttr func(o *reflect.Value, name string) (*reflect.Value, func(interface{}), error)
	getAttr = func(o *reflect.Value, name string) (*reflect.Value, func(interface{}), error) {
		names := []string{name}
		if name == "Case" {
			names = append(names, "case")
		}
		if o.Kind() == reflect.Map {
			for _, name := range names {
				key := reflect.ValueOf(name)
				val := o.MapIndex(key)
				if val.Kind() != reflect.Invalid {
					return &val, mapSetter(o, &key), nil
				}
			}
			return nil, nil, fmt.Errorf("can't find key: %v in map", name)
		}
		if o.Kind() == reflect.Interface || o.Kind() == reflect.Ptr {
			tmp := o.Elem()
			return getAttr(&tmp, name)
		}
		if o.Kind() != reflect.Struct {
			return nil, nil, fmt.Errorf("object kind must be map or struct, not %v", o.Kind().String())
		}
		val := o.FieldByName(name)

		if isValid(&val) {
			return &val, fieldSetter(&val), nil
		}
		valueType := reflect.TypeOf(o.Interface())
		for i := 0; i < valueType.NumField(); i++ {
			field := valueType.Field(i)
			found := false
			for _, name := range names {
				for _, tag := range tagNames {
					if strings.Split(field.Tag.Get(tag), ",")[0] == name {
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
		if isValid(&val) {
			return &val, fieldSetter(&val), nil
		}
		return nil, nil, fmt.Errorf("can't find attr %s", name)
	}
	var getIndex func(o *reflect.Value, index *reflect.Value) (*reflect.Value, func(interface{}), error)
	getIndex = func(o *reflect.Value, index *reflect.Value) (*reflect.Value, func(interface{}), error) {
		if o.Kind() == reflect.Interface || o.Kind() == reflect.Ptr {
			tmp := o.Elem()
			return getIndex(&tmp, index)
		}
		switch o.Kind() {
		case reflect.Array, reflect.Slice:
			i, err := cast.ToIntE(index.Interface())
			if err != nil {
				return nil, nil, fmt.Errorf("cast slice or array index to int err: %w", err)
			}
			if i < 0 || i >= o.Len() {
				return nil, nil, fmt.Errorf("slice or array index is out of range: %d", i)
			}
			tmp := o.Index(i)
			return &tmp, fieldSetter(&tmp), nil
		case reflect.Map:
			if index.Kind() == reflect.Interface {
				tmp := index.Elem()
				return getIndex(o, &tmp)
			}
			if index.Kind() == reflect.String {
				tmp := reflect.ValueOf(strings.Trim(index.String(), " \""))
				index = &tmp
			}
			tmp := o.MapIndex(*index)
			return &tmp, mapSetter(o, &tmp), nil
		}
		return nil, nil, errors.New("index object kind is not in (map, array, slice)")
	}

	v := root
	var setter func(interface{})
	for _, field := range fields {
		if field.isIndex {
			v, setter, err = getAttr(v, field.fieldName)
			if err != nil {
				break
			}
			tmp := reflect.ValueOf(field.index)
			v, setter, err = getIndex(v, &tmp)
		} else {
			v, setter, err = getAttr(v, field.fieldName)
		}
		if err != nil {
			break
		}
	}
	return v, setter, err

	//invalidErr := errors.New("invalid stmt")
	/*
		var walk func(node ast.Expr) (*reflect.Value, func(interface{}), error)
		walk = func(node ast.Expr) (*reflect.Value, func(interface{}), error) {
			switch node.(type) {
			case *ast.IndexExpr:
				indexExpr := node.(*ast.IndexExpr)
				index, _, err := walk(indexExpr.Index)
				if err != nil {
					return nil, nil, err
				}
				val, _, err := walk(indexExpr.X)
				if err != nil {
					return nil, nil, err
				}
				val, setter, err := getIndex(val, index)
				if err != nil {
					return nil, nil, fmt.Errorf("%w at offset %v", err, indexExpr.Pos())
				}
				return val, setter, nil
			case *ast.SelectorExpr:
				selectorExpr := node.(*ast.SelectorExpr)
				n := selectorExpr.Sel.Name
				val, _, err := walk(selectorExpr.X)
				if err != nil {
					return nil, nil, err
				}
				val, setter, err := getAttr(val, n)
				if err != nil {
					return nil, nil, fmt.Errorf("%w at offset %d", err, selectorExpr.Sel.Pos())
				}
				return val, setter, nil
			case *ast.BasicLit:
				bascLit := node.(*ast.BasicLit)
				switch bascLit.Kind {
				case token.INT:
					tmp := reflect.ValueOf(cast.ToInt(bascLit.Value))
					return &tmp, nil, nil
				case token.STRING:
					tmp := reflect.ValueOf(bascLit.Value)
					return &tmp, nil, nil
				}
				return nil, nil, fmt.Errorf("unknown basicLit: %s", bascLit.Kind.String())
			case *ast.Ident:
				identExpr := node.(*ast.Ident)
				val, setter, err := getAttr(root, identExpr.String())
				if err != nil {
					return nil, nil, fmt.Errorf("%w at offset %d", err, identExpr.Pos())
				}
				return val, setter, nil
			}
			return nil, nil, fmt.Errorf("unknown node type: %s at offset %d", reflect.TypeOf(node).String(), node.Pos())
		}
		return walk(expr)
	*/
}

// GetObjectValue get obj value
func GetObjectValue(obj interface{}, attr string) (interface{}, error) {
	result, _, err := ReflectGetValue(obj, attr)
	if err != nil {
		return nil, err
	}
	return result.Interface(), nil
}

// GetObjectValueAndSetter description
// @param obj
// @param attr
// @return interface{}
// @return func(interface{})
// @return error
func GetObjectValueAndSetter(obj interface{}, attr string) (interface{}, func(interface{}), error) {
	result, setter, err := ReflectGetValue(obj, attr)
	if err != nil {
		return nil, nil, err
	}
	return result.Interface(), setter, nil
}

/*
// GetObjectValue get value of object
func GetObjectValue(obj interface{}, attr string) (interface{}, error) {
	attr = strings.Trim(attr, " ")
	attrParts := strings.Split(attr, ".")
	result := obj
	for _, part := range attrParts {
		pos := strings.IndexByte(part, '[')
		if pos != -1 {
			// 简单的数组取值
			// TODO: 嵌套数组，多维数组
			name := strings.Trim(part[:pos], " ")
			end := strings.IndexByte(part[pos+1:], ']')
			if end == -1 {
				return nil, errors.New("invalid attr")
			}
			index := cast.ToInt(part[pos+1 : pos+end+1])
			tmp := ToAbstractMap(result).GetAbstractSlice(name)
			if tmp.Len() > index {
				result = tmp.Value(index)
			}
		} else {
			result = ToAbstractMap(result).GetValue(part)
		}
		if result == nil {
			return nil, errors.New("invalid attr")
		}
	}
	return result, nil
}
*/

// GetObjectValueString get value string of object
func GetObjectValueString(obj interface{}, attr string) (string, error) {
	val, err := GetObjectValue(obj, attr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%+v", val), nil
}

// GetFirstStandardVar description
// @param obj
// @param statement
// @return interface{}
// @return error
func GetFirstStandardVar(obj interface{}, statement string) (interface{}, error) {
	re := regexp.MustCompile(`\${([a-z|A-Z|0-9|\.]*)}`)
	matched := re.FindStringSubmatch(statement)
	if len(matched) == 0 {
		return nil, errors.New("can't find var")
	}
	return GetObjectValue(obj, matched[1])
}

// GetStandardVar description
// @param obj
// @param statement
// @return interface{}
// @return error
func GetStandardVar(obj interface{}, statement string) (interface{}, error) {
	re := regexp.MustCompile(`^\${([a-z|A-Z|0-9|\.|_|\[|\]]*)}$`)
	matched := re.FindStringSubmatch(statement)
	if len(matched) == 0 {
		return nil, ErrStandardVarNotMatched
	}
	return GetObjectValue(obj, matched[1])
}

// GetStandardVarAndSetter description
// @param obj
// @param statement
// @return interface{}
// @return func(interface{})
// @return error
func GetStandardVarAndSetter(obj interface{}, statement string) (interface{}, func(interface{}), error) {
	re := regexp.MustCompile(`^\${([a-z|A-Z|0-9|\.|_|\[|\]]*)}$`)
	matched := re.FindStringSubmatch(statement)
	if len(matched) == 0 {
		return nil, nil, errors.New("can't find var")
	}
	return GetObjectValueAndSetter(obj, matched[1])
}

// FormatObject format object to string
func FormatObject(format string, obj interface{}) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.WithContext(context.Background()).Errorf("%+v\n %s", r, PrintStack())
		}
	}()
	re := regexp.MustCompile(`\${([a-z|A-Z|0-9|\.|_|\[|\]]*)}`)
	return re.ReplaceAllStringFunc(format, func(s string) string {
		state := re.FindStringSubmatch(s)[1]
		val, e := GetObjectValueString(obj, state)
		if e != nil {
			err = e
			panic(err)
		}
		return val
	}), nil
}

// StringSliceContains description
// @param slice
// @param need
// @return bool
func StringSliceContains(slice []string, need string) bool {
	for _, val := range slice {
		if val == need {
			return true
		}
	}
	return false
}

// WaitUntil description
// @param ctx
// @param waitFunc
// @param interval
func WaitUntil(ctx context.Context, waitFunc func() bool, interval time.Duration) {
	if waitFunc() {
		return
	}
	if int(interval) == 0 {
		interval = time.Second
	}
	tick := time.NewTicker(interval)
	for {
		select {
		case <-tick.C:
			if waitFunc() {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// CreateZKPathP description
// @param ctx
// @param zkCli
// @param path
// @return error
func CreateZKPathP(ctx context.Context, zkCli *zk.Conn, path string) error {
	if path[0] != '/' {
		return errors.New("invalid path")
	}
	pathes := strings.Split(path[1:], "/")
	toCreate := ""
	for i := 0; i < len(pathes); i++ {
		toCreate += ("/" + pathes[i])
		_, err := zkCli.Create(toCreate, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists && err != zk.ErrNoChildrenForEphemerals {
			return err
		}
	}
	return nil
}

// PanicWhenError description
// @param err
func PanicWhenError(err error) {
	if err != nil {
		panic(err)
	}
}

// DeepCopyJSON description
// @param src
// @return map
func DeepCopyJSON(src map[string]interface{}) map[string]interface{} {
	dest := make(map[string]interface{})
	for key, value := range src {
		switch src[key].(type) {
		case map[string]interface{}:
			dest[key] = DeepCopyJSON(src[key].(map[string]interface{}))
		default:
			dest[key] = value
		}
	}
	return dest
}

func SliceIndex[T comparable](src []T, find T) int {
	for i, v := range src {
		if find == v {
			return i
		}
	}
	return -1
}

func IsNilInterface(i interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}

func RowDataToCol[K comparable, V any](data []map[K]V) map[K][]V {
	out := make(map[K][]V)
	lo.ForEach(data, func(item map[K]V, index int) {
		for k, v := range item {
			if _, ok := out[k]; !ok {
				out[k] = make([]V, index)
			}
			out[k] = append(out[k], v)
		}
	})
	return out
}
