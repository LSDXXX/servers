package util

import (
	"fmt"
	"testing"
)

func TestGetStandardObj(t *testing.T) {
	tmp := map[string]interface{}{
		"msg": map[string]interface{}{
			"payload": map[string]interface{}{
				"method": "test",
			},
		},
	}
	fmt.Println(GetStandardVar(tmp, "${msg.payload.method}"))
}

func TestReflectGetValue(t *testing.T) {
	tmp := map[string]interface{}{
		"info": map[string]interface{}{
			"name": "test",
		},
	}
	type Test struct {
		Value []string `json:"value"`
	}
	tt := Test{
		Value: []string{"22"},
	}
	ss := "123"
	_, setter, _ := ReflectGetValue(&tt, "value[0]")
	setter(ss)
	fmt.Println(tmp, tt.Value)
}

func TestFormatVal(t *testing.T) {
	fmt.Println(FormatObject("${a}+${b}", map[string]string{
		"a": "123",
		"b": "456",
	}))
}
