package js

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/util/structure"
	"github.com/robertkrimen/otto"
)

var (

	// ErrVMNotStartedYet .
	ErrVMNotStartedYet = errors.New("vm not started yet")

	// ErrVMAlreadyStarted .
	ErrVMAlreadyStarted = errors.New("vm already started")
)

const year = time.Hour * 24 * 365

// Script description
type Script struct {
	root             *otto.Otto
	maxWorkerRoutine int32
	vmCount          int32
	initialCount     int
	mu               sync.RWMutex
	ch               chan task
	started          int32
	maxIdle          time.Duration
}

type result struct {
	val otto.Value
	err error
}

type task struct {
	pre     func(*otto.Otto)
	timeout time.Duration
	args    map[string]interface{}
	resch   chan result
	code    *otto.Script
	ctx     context.Context
}

type option func(*Script)

// WithMaxRoutine description
// @param n
// @return option
func WithMaxRoutine(n int) option {
	return func(s *Script) {
		s.maxWorkerRoutine = int32(n)
	}
}

// WithInitialCount description
// @param n
// @return option
func WithInitialCount(n int) option {
	return func(s *Script) {
		s.initialCount = n
	}
}

// New description
// @param opts
// @return *Script
func New(opts ...option) *Script {
	s := &Script{
		root:             otto.New(),
		maxWorkerRoutine: 1024,
		vmCount:          0,
		initialCount:     32,
		ch:               make(chan task),
		maxIdle:          time.Minute * 10,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start description
// @receiver s
// @return error
func (s *Script) Start() error {
	if !atomic.CompareAndSwapInt32(&s.started, 0, 1) {
		return ErrVMAlreadyStarted
	}
	for i := 0; i < s.initialCount; i++ {
		go s.workLoop(true)
	}
	return nil
}

// RoutineCount description
// @receiver s
// @return int32
func (s *Script) RoutineCount() int32 {
	return atomic.LoadInt32(&s.vmCount)
}

func (s *Script) exec(vm *otto.Otto, t *task) (val otto.Value, err error) {
	ctx := t.ctx
	if t.pre != nil {
		t.pre(vm)
	}
	for k, v := range t.args {
		_ = vm.Set(k, v)
	}
	defer func() {
		if r := recover(); r != nil {
			log.WithContext(ctx).Errorf("exec js script panic: %v", r)
		}
	}()
	vm.Interrupt = make(chan func(), 1)
	cancel := make(chan struct{})
	defer close(cancel)
	if t.timeout == 0 {
		go func() {
			select {
			case <-time.After(t.timeout):
				vm.Interrupt <- func() {
					err = fmt.Errorf("script time out, code: %s", t.code.String())
					panic("script time out")
				}
			case <-cancel:
				return
			}
		}()
	}
	val, err = vm.Run(t.code)
	return
}

func (s *Script) workLoop(permanent bool) {
	vm := s.root.Copy()
	idle := s.maxIdle
	if permanent {
		idle = year
	}
	ticker := time.NewTicker(idle)
	atomic.AddInt32(&s.vmCount, 1)
	defer atomic.AddInt32(&s.vmCount, -1)
	for {
		select {
		case t := <-s.ch:
			val, err := s.exec(vm, &t)
			t.resch <- result{
				err: err,
				val: val,
			}
			ticker.Reset(s.maxIdle)
		case <-ticker.C:
			return
		}
	}
}

// Compile description
// @receiver s
// @param src
// @return *otto.Script
// @return error
func (s *Script) Compile(src interface{}) (*otto.Script, error) {
	vm := s.root.Copy()
	return vm.Compile("", src)
}

// SetFunc description
// @receiver s
// @param name
// @param f
// @return error
func (s *Script) SetFunc(name string, f func(call otto.FunctionCall) otto.Value) error {
	return s.root.Set(name, f)
}

// ToGoValue description
// @param val
// @return interface{}
func ToGoValue(val otto.Value) interface{} {
	if val.IsBoolean() {
		out, _ := val.ToBoolean()
		return out
	} else if val.IsNumber() {
		out, _ := val.ToFloat()
		return out
	} else if val.IsString() {
		out, _ := val.ToString()
		return out
	} else if val.IsObject() {
		out, _ := val.Export()
		return structure.ToAbstractMap(out)
	}
	return nil
}

// RunCode description
// @receiver s
// @param ctx
// @param code
// @param args
// @param pre
// @param timeout
// @return otto.Value
// @return error
func (s *Script) RunCode(ctx context.Context, code *otto.Script, args map[string]interface{},
	pre func(*otto.Otto), timeout time.Duration) (otto.Value, error) {
	if atomic.LoadInt32(&s.started) != 1 {
		return otto.NullValue(), ErrVMNotStartedYet
	}
	resch := make(chan result)
	t := task{
		pre:     pre,
		timeout: timeout,
		args:    args,
		ctx:     ctx,
		code:    code,
		resch:   resch,
	}
	select {
	case s.ch <- t:
		res := <-t.resch
		return res.val, res.err
	default:
		if atomic.AddInt32(&s.vmCount, 1) < s.maxWorkerRoutine {
			go s.workLoop(false)
		}
		s.ch <- t
		res := <-t.resch
		return res.val, res.err
	}
}
