package expr

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type cases struct {
	expr   string
	env    map[string]interface{}
	assert func(t *testing.T, v interface{}) bool
}

func TestExpr(t *testing.T) {
	testCases := map[string]cases{
		"time": {
			expr: `ISOTIMESTAMP($timestampSec-7200)`,
			assert: func(t *testing.T, v interface{}) bool {
				fmt.Println(v)
				return true
			},
		},
		"format": {
			expr: `DATE_FORMAT($timestampSec, 'YYYY.MM.DD')`,
			assert: func(t *testing.T, v interface{}) bool {
				now := time.Now()
				s := fmt.Sprintf("%04d.%02d.%02d", now.Year(), now.Month(), now.Day())
				return assert.Equal(t, s, v)
			},
		},
		"$param": {
			expr: `15*field`,
			assert: func(t *testing.T, v interface{}) bool {
				return assert.Equal(t, 30, v)
			},
			env: map[string]interface{}{
				"field": 2,
			},
		},
	}

	for name, c := range testCases {
		fmt.Printf("-------- test case %s --------\n", name)
		ex, err := Parse(c.expr)
		if err != nil {
			t.Fatal(err)
		}
		res, err := ex.Eval(c.env)
		if err != nil {
			t.Fatal(err)
		}
		if !c.assert(t, res) {
			t.Fatal(err)
		}
	}
}
