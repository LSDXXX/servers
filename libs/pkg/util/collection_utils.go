package util

import (
	"reflect"
	"strconv"
)

func ListToMap[T any](list []T, key string) map[string]T {
	res := make(map[string]T)
	arr := ToSlice[T](list)
	for _, row := range arr {
		immutable := reflect.ValueOf(row)
		var val string
		if immutable.FieldByName(key).Type().Kind() == reflect.Int {
			val = strconv.FormatInt(immutable.FieldByName(key).Int(), 10)
		} else {
			val = immutable.FieldByName(key).String()
		}
		res[val] = row
	}
	return res
}
func GetFieldsFromList[T any](arr []T, field, ignore string) []string {
	res := make([]string, 0)
	for _, v := range arr {
		immutable := reflect.ValueOf(v)
		var value string
		if immutable.FieldByName(field).Type().Kind() == reflect.Int {
			value = strconv.FormatInt(immutable.FieldByName(field).Int(), 10)
		} else {
			value = immutable.FieldByName(field).String()
		}
		if value != ignore {
			res = append(res, value)
		}
	}
	return res
}

func ToSlice[T any](arr []T) []T {
	ret := make([]T, 0)
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		ret = append(ret, arr...)
		return ret
	}
	l := v.Len()
	for i := 0; i < l; i++ {
		ret = append(ret, v.Index(i).Interface().(T))
	}
	return ret
}

func StringToInt(arr []string) []int {
	res := make([]int, 0)
	for _, v := range arr {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		res = append(res, i)
	}
	return res
}

func Intersect(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

func Contains(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}
	return false
}
