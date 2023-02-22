package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const DateFormat = "2006-01-02"

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(DateFormat))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}
	parsed, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

func (d Date) String() string {
	return d.Time.Format(DateFormat)
}

type Binder interface {
	Bind(src string) error
}

func BindStringToObject(src string, dst interface{}) error {
	var err error

	v := reflect.ValueOf(dst)
	t := reflect.TypeOf(dst)

	// We need to dereference pointers
	if t.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
		t = v.Type()
	}

	// For some optioinal args
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}

		v = reflect.Indirect(v)
		t = v.Type()
	}

	// The resulting type must be settable. reflect will catch issues like
	// passing the destination by value.
	if !v.CanSet() {
		return errors.New("destination is not settable")
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var val int64
		val, err = strconv.ParseInt(src, 10, 64)
		if err == nil {
			v.SetInt(val)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var val uint64
		val, err = strconv.ParseUint(src, 10, 64)
		if err == nil {
			v.SetUint(val)
		}
	case reflect.String:
		v.SetString(src)
		err = nil
	case reflect.Float64, reflect.Float32:
		var val float64
		val, err = strconv.ParseFloat(src, 64)
		if err == nil {
			v.SetFloat(val)
		}
	case reflect.Bool:
		var val bool
		val, err = strconv.ParseBool(src)
		if err == nil {
			v.SetBool(val)
		}
	case reflect.Struct:
		// if this is not of type Time or of type Date look to see if this is of type Binder.
		if dstType, ok := dst.(Binder); ok {
			return dstType.Bind(src)
		}

		if t.ConvertibleTo(reflect.TypeOf(time.Time{})) {
			// Don't fail on empty string.
			if src == "" {
				return nil
			}
			// Time is a special case of a struct that we handle
			parsedTime, err := time.Parse(time.RFC3339Nano, src)
			if err != nil {
				parsedTime, err = time.Parse(DateFormat, src)
				if err != nil {
					return fmt.Errorf("error parsing '%s' as RFC3339 or 2006-01-02 time: %s", src, err)
				}
			}
			// So, assigning this gets a little fun. We have a value to the
			// dereference destination. We can't do a conversion to
			// time.Time because the result isn't assignable, so we need to
			// convert pointers.
			if t != reflect.TypeOf(time.Time{}) {
				vPtr := v.Addr()
				vtPtr := vPtr.Convert(reflect.TypeOf(&time.Time{}))
				v = reflect.Indirect(vtPtr)
			}
			v.Set(reflect.ValueOf(parsedTime))
			return nil
		}

		if t.ConvertibleTo(reflect.TypeOf(Date{})) {
			// Don't fail on empty string.
			if src == "" {
				return nil
			}
			parsedTime, err := time.Parse(DateFormat, src)
			if err != nil {
				return fmt.Errorf("error parsing '%s' as date: %s", src, err)
			}
			parsedDate := Date{Time: parsedTime}

			// We have to do the same dance here to assign, just like with times
			// above.
			if t != reflect.TypeOf(Date{}) {
				vPtr := v.Addr()
				vtPtr := vPtr.Convert(reflect.TypeOf(&Date{}))
				v = reflect.Indirect(vtPtr)
			}
			v.Set(reflect.ValueOf(parsedDate))
			return nil
		}

		// We fall through to the error case below if we haven't handled the
		// destination type above.
		fallthrough
	default:
		// We've got a bunch of types unimplemented, don't fail silently.
		err = fmt.Errorf("can not bind to destination of type: %s", t.Kind())
	}
	if err != nil {
		return fmt.Errorf("error binding string parameter: %s", err)
	}
	return nil
}

type ServiceInterface[T any] interface {
	WithContext(context.Context) T
}

type InjectServices1[T ServiceInterface[T]] interface{}
type InjectServices2[T ServiceInterface[T], T1 ServiceInterface[T1]] interface{}
type InjectServices3[T ServiceInterface[T], T1 ServiceInterface[T1], T2 ServiceInterface[T2]] interface{}
type InjectServices4[T ServiceInterface[T],
	T1 ServiceInterface[T1], T2 ServiceInterface[T2], T3 ServiceInterface[T3]] interface{}
type InjectServices5[T ServiceInterface[T],
	T1 ServiceInterface[T1], T2 ServiceInterface[T2], T3 ServiceInterface[T3], T4 ServiceInterface[T4]] interface{}
type InjectServices6[T ServiceInterface[T],
	T1 ServiceInterface[T1], T2 ServiceInterface[T2], T3 ServiceInterface[T3],
	T4 ServiceInterface[T4], T5 ServiceInterface[T5]] interface{}
type InjectServices7[T ServiceInterface[T],
	T1 ServiceInterface[T1], T2 ServiceInterface[T2], T3 ServiceInterface[T3],
	T4 ServiceInterface[T4], T5 ServiceInterface[T5], T6 ServiceInterface[T6]] interface{}
