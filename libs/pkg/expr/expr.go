package expr

import (
	"regexp"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/pkg/errors"
	"gitlab.com/metakeule/fmtdate"
)

var (
	globalParams map[string]builtinParamConf
	paramRe      = regexp.MustCompile(`\$([a-z|A-Z|0-9|_|]*)`)
)

func init() {
	globalParams = map[string]builtinParamConf{
		"$timestampSec": {
			name:     "currentTime",
			replaced: "currentTime()",
			value: func() int64 {
				return time.Now().Unix()
			},
		},
		"$isoTimestamp": {
			name:     "isoCurrentTime",
			replaced: "isoCurrentTime()",
			value: func() string {
				return time.Now().Format(time.RFC3339)
			},
		},
		"TIMESTAMPSEC": {
			name:     "TIMESTAMPSEC",
			replaced: "",
			value: func(t string) int64 {
				res, err := time.Parse(time.RFC3339, t)
				if err != nil {
					return 0
				}
				return res.Unix()
			},
		},
		"ISOTIMESTAMP": {
			name:     "ISOTIMESTAMP",
			replaced: "",
			value: func(t int64) string {
				res := time.Unix(t, 0)
				return res.Format(time.RFC3339)
			},
		},
		"DATE_FORMAT": {
			name: "DATE_FORMAT",
			value: func(t int64, format string) string {
				tt := time.Unix(t, 0)
				return fmtdate.Format(format, tt)
			},
		},
	}
}

type Expr struct {
	program *vm.Program
}

func (e *Expr) Eval(env map[string]interface{}) (interface{}, error) {
	envMap := make(map[string]interface{})
	for k, v := range env {
		envMap[k] = v
	}
	for _, conf := range globalParams {
		envMap[conf.name] = conf.value
	}
	return expr.Run(e.program, envMap)
}

type builtinParamConf struct {
	name     string
	replaced string
	value    interface{}
}

func replace(code string) string {
	return paramRe.ReplaceAllStringFunc(code, func(s string) string {
		find, ok := globalParams[s]
		if ok && len(find.replaced) > 0 {
			return find.replaced
		}
		return s
	})
}

func Parse(code string) (*Expr, error) {
	code = replace(code)
	p, err := expr.Compile(code)
	if err != nil {
		return nil, errors.Wrap(err, "compile code")
	}
	ex := &Expr{
		program: p,
	}
	return ex, nil
}

func Eval(code string, env map[string]interface{}) (interface{}, error) {
	code = replace(code)
	return expr.Eval(code, env)
}
