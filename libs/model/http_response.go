package model

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/LSDXXX/libs/pkg/errorcode"
	"github.com/pkg/errors"
)

// Response http 恢复
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Msg     string      `json:"msg,omitempty"`
	ErrData interface{} `json:"errData,omitempty"`
	Data    interface{} `json:"data,omitempty"`

	StatusCode int               `json:"-"`
	Headers    map[string]string `json:"-"`
}

// ResponseOption opts
type ResponseOption func(*Response)

// WithCode with code
func WithCode(n int) ResponseOption {
	return func(res *Response) {
		res.Code = n
	}
}

// WithMessage with message
func WithMessage(msg string) ResponseOption {
	return func(res *Response) {
		res.Message = msg
	}
}

// WithCodeMessage with code and message
func WithCodeMessage(code int, msg string) ResponseOption {
	return func(res *Response) {
		res.Code = code
		res.Message = msg
	}
}

// WithData with data
func WithData(data interface{}) ResponseOption {
	return func(res *Response) {
		res.Data = data
	}
}

func WithMsg(msg string) ResponseOption {
	return func(res *Response) {
		res.Message = msg
	}
}

// WithError with error
func WithError(err error, data ...interface{}) ResponseOption {
	return func(res *Response) {
		if err == nil {
			return
		}
		res.Code = errorcode.Code(err)
		res.Message = err.Error()
		if len(data) == 1 {
			res.ErrData = data[0]
		}
	}
}

// NewResponse create http response
func NewResponse(opts ...ResponseOption) *Response {
	res := &Response{
		Message:    "ok",
		Code:       0,
		StatusCode: 200,
		Headers:    make(map[string]string),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func FailResponse(opts ...ResponseOption) *Response {
	res := &Response{
		Message: "fail",
		Code:    -1,
		Headers: make(map[string]string),
		Data:    nil,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// ReadFromBody read from bocy
//
//	@param rc
//	@return *Response
//	@return error
func ReadFromBody(rc io.ReadCloser) (*Response, error) {
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "read request body")
	}
	var out Response
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal json: "+string(data))
	}
	return &out, nil
}

// Message entity
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GenericResponse generic entity
type GenericResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Msg     string `json:"msg,omitempty"`
	Data    T      `json:"Data,omitempty"`
}

// GenericReadFromBody read
func GenericReadFromBody[T any](rc io.ReadCloser) (*GenericResponse[T], error) {
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "read request body")
	}
	var out GenericResponse[T]
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal json: "+string(data))
	}
	return &out, nil
}
