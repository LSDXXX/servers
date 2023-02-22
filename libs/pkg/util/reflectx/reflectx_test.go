package reflectx

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ST struct {
	Name  string                 `json:"name"`
	Value string                 `json:"value"`
	Data  map[string]interface{} `json:"data"`
}

type Test struct {
	Name         string                       `json:"name"`
	PName        *string                      `json:"pname"`
	StringSlice  []string                     `json:"string_slice"`
	StringSliceP []*string                    `json:"string_slicep"`
	StructSlice  []ST                         `json:"struct_slice"`
	IntSlice     []int                        `json:"int_slice"`
	FloatSlice   []float64                    `json:"float_slice"`
	STS          map[string]interface{}       `json:"sts"`
	MapString    map[string]string            `json:"map_string"`
	MapMapString map[string]map[string]string `json:"map_map_string"`
	StructP      *ST                          `json:"st"`
	MapStruct    map[string]ST                `json:"map_struct"`
}

type testCase struct {
	attr   string
	value  interface{}
	assert func(t *testing.T, v Test) bool
}

func TestSet(t *testing.T) {
	milli := time.Now().UnixMilli()
	fmt.Print(milli)

	testCases := map[string]testCase{
		"test_field_string": {
			attr:  "name",
			value: "test",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.Name, "test", v.Name)
			},
		},
		"test_field_pstring": {
			attr:  "pname",
			value: "test",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, *v.PName, "test")
			},
		},
		"test_st": {
			attr:  "st.name",
			value: "haha",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.StructP.Name, "haha")
			},
		},
		"test_field_stringslice": {
			attr:  "string_slice",
			value: []string{"1"},
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, []string{"1"}, v.StringSlice)
			},
		},
		"test_field_stringslicen": {
			attr:  "string_slice[2]",
			value: "1",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.StringSlice[2], "1")
			},
		},
		"test_field_structslice": {
			attr: "struct_slice",
			value: []map[string]interface{}{
				{
					"value": "haha",
				},
			},
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.StructSlice[0].Value, "haha")
			},
		},
		"test_string_slicep": {
			attr:  "string_slicep",
			value: []interface{}{"10", "20"},
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, *v.StringSliceP[0], "10")
			},
		},
		"test_string_slicepn": {
			attr:  "string_slicep[1]",
			value: "10",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, *v.StringSliceP[1], "10")
			},
		},
		"test_sts": {
			attr:  "sts.info.name",
			value: "test",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.STS["info"].(map[string]interface{})["name"].(string), "test")
			},
		},
		"test_sts_s": {
			attr:  "sts.info.name[0]",
			value: "test",
			assert: func(t *testing.T, v Test) bool {
				fmt.Println(v.STS)
				return assert.Equal(t,
					v.STS["info"].(map[string]interface{})["name"].([]interface{})[0],
					"test")
			},
		},
		"test_int_slice": {
			attr:  "int_slice",
			value: []interface{}{0, 2, 3},
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.IntSlice[1], 2)
			},
		},
		"test_float_slice": {
			attr:  "float_slice",
			value: []interface{}{0, 2, 3},
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.FloatSlice[1], float64(2))
			},
		},
		"test_map_string": {
			attr:  "map_string.test",
			value: "haha",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.MapString["test"], "haha")
			},
		},

		"test_map_map_string": {
			attr:  "map_map_string.test.test",
			value: "haha",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.MapMapString["test"]["test"], "haha")
			},
		},
		"test_map_struct": {
			attr:  "map_struct.test.name",
			value: "haha",
			assert: func(t *testing.T, v Test) bool {
				return assert.Equal(t, v.MapStruct["test"].Name, "haha")
			},
		},
	}
	for k, v := range testCases {
		var tt Test
		fmt.Printf("case %s started  ..\n", k)
		err := SetAttr(&tt, v.attr, v.value)
		if err != nil {
			t.Fatal(fmt.Errorf("test case %s, occur err: %+v", k, err))
		}
		if !v.assert(t, tt) {
			t.Fatal(fmt.Errorf("test case %s, asset fail", k))
		}
		fmt.Printf("case %s passed \n", k)
	}
}

func TestParseExpr(t *testing.T) {
	fmt.Println(ParseExpr("test.haha[1][2].test2", false))
}

func TestGetAttr(t *testing.T) {
	type test struct {
		Value  string            `json:"value"`
		Data   map[string]string `json:"data"`
		Array  []string          `json:"array"`
		Arrays [][]string        `json:"arrays"`
	}

	type cases struct {
		attr   string
		assert func(t *testing.T, v interface{}) bool
	}
	v := test{
		Value: "haha",
		Data: map[string]string{
			"123":   "tt",
			"test2": "ttt",
		},
		Array: []string{"1", "2", "3"},
		Arrays: [][]string{
			{"1", "2", "3"},
		},
	}
	testCases := map[string]cases{
		/*
			"test_field": {
				attr: "value",
				assert: func(t *testing.T, res interface{}) bool {
					return assert.Equal(t, v.Value, res)
				},
			},
		*/
		"test_data": {
			attr: "data.123",
			assert: func(t *testing.T, res interface{}) bool {
				return assert.Equal(t, v.Data["123"], res)
			},
		},
		/*
			"test_array": {
				attr: "array[2]",
				assert: func(t *testing.T, res interface{}) bool {
					return assert.Equal(t, v.Array[2], res)
				},
			},
		*/
		"test_arrays": {
			attr: "arrays[0][2]",
			assert: func(t *testing.T, res interface{}) bool {
				return assert.Equal(t, v.Arrays[0][2], res)
			},
		},
	}
	for name, c := range testCases {
		fmt.Printf("start test case: %s\n", name)
		res, _, err := GetAttr(v, c.attr)
		if err != nil {
			t.Fatal(err)
		}
		if !c.assert(t, res.Interface()) {
			t.Fail()
		}
		fmt.Printf("test case %s result success \n", name)
	}

}

func TestWrapper(t *testing.T) {
	wrapper := GenerateTypeWrapper(reflect.TypeOf(Test{}), "json")
	_ = wrapper
}
