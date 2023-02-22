package servercontext

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/LSDXXX/libs/constant"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

var ctxKey = key{}

type key struct{}

const GINServerContextKey = "server-context"

// Context description
type Context struct {
	logger         *logrus.Entry
	DB             *gorm.DB
	TraceID        string
	SpanID         string
	RequestID      string
	Token          string
	ProjectID      int
	DisablePushLog bool
	GinContext     *gin.Context
}

type request struct {
	Context traceInfo `json:"context"`
}

type traceInfo struct {
	TraceID string `json:"trace_id"`
	Token   string `json:"token"`
}

// ExtractFromHTTP description
// @param c
// @param req
// @return context.Context
func ExtractFromHTTP(c context.Context, ginc *gin.Context) context.Context {
	var ctx Context
	req := ginc.Request
	ctx.TraceID = req.Header.Get(constant.TraceIDHTTPHeader)
	ctx.SpanID = req.Header.Get(constant.SpanIDHTTPHeader)
	ctx.RequestID = req.Header.Get(constant.RequestIDHTTPHeader)
	ctx.Token = req.Header.Get(constant.TokenHTTPHeader)
	ctx.ProjectID = cast.ToInt(req.Header.Get(constant.ProjectIDHTTPHeader))
	if len(req.Header.Get(constant.DisablePushLogHeader)) != 0 {
		ctx.DisablePushLog = true
	}
	if projectId, ok := ginc.GetQuery("projectId"); ok {
		ctx.ProjectID = cast.ToInt(projectId)
	}
	if len(ctx.TraceID) == 0 {
		data, err := ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err == nil {
			var info request
			err = json.Unmarshal(data, &info)
			if err == nil {
				ctx.TraceID = info.Context.TraceID
				ctx.Token = info.Context.Token
			}
		}
		req.Body = io.NopCloser(bytes.NewBuffer(data))
	}
	newCtx := context.WithValue(context.Background(), ctxKey, &ctx)
	return newCtx
}

// ContextToHTTP description
// @param ctx
// @param req
func ContextToHTTP(ctx context.Context, req *http.Request) {
	c := Get(ctx)
	if c == nil {
		return
	}
	c.WriteHeader(req.Header)
}

// ContextToGrpc description
// @param ctx
// @return context.Context
func ContextToGrpc(ctx context.Context) context.Context {
	c := Get(ctx)
	if c == nil {
		c = &Context{}
	}
	if len(c.TraceID) == 0 {
		c.TraceID = uuid.New().String()
	}
	_, ok := metadata.FromIncomingContext(ctx)
	kvs := []string{constant.TraceIDHTTPHeader, c.TraceID,
		constant.SpanIDHTTPHeader, c.SpanID,
		constant.TokenHTTPHeader, c.Token,
		constant.RequestIDHTTPHeader, c.RequestID,
		constant.ProjectIDHTTPHeader, cast.ToString(c.ProjectID),
	}
	if ok {
		ctx = metadata.AppendToOutgoingContext(ctx, kvs...)
	} else {
		md := metadata.New(map[string]string{
			constant.TraceIDHTTPHeader:   c.TraceID,
			constant.SpanIDHTTPHeader:    c.SpanID,
			constant.TokenHTTPHeader:     c.Token,
			constant.RequestIDHTTPHeader: c.RequestID,
			constant.ProjectIDHTTPHeader: cast.ToString(c.ProjectID),
		})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func getMD(md metadata.MD, key string) string {
	v := md.Get(key)
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

// ExtractFromGrpc description
// @param ctx
// @return context.Context
func ExtractFromGrpc(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return context.WithValue(ctx, ctxKey, &Context{})
	}
	var c Context
	c.TraceID = getMD(md, constant.TraceIDHTTPHeader)
	c.SpanID = getMD(md, constant.SpanIDHTTPHeader)
	c.RequestID = getMD(md, constant.RequestIDHTTPHeader)
	c.Token = getMD(md, constant.TokenHTTPHeader)
	c.ProjectID = cast.ToInt(getMD(md, constant.ProjectIDHTTPHeader))
	return context.WithValue(ctx, ctxKey, &c)
}

// GetLogger description
// @receiver c
// @return *logrus.Entry
func (c *Context) GetLogger() *logrus.Entry {
	return c.logger
}

// SetLogger description
// @receiver c
// @param logger
func (c *Context) SetLogger(logger *logrus.Entry) {
	c.logger = logger
}

// WriteHeader description
// @receiver c
// @param header
func (c *Context) WriteHeader(header http.Header) {
	if len(c.TraceID) != 0 {
		header.Set(constant.TraceIDHTTPHeader, c.TraceID)
	}
	if len(c.SpanID) != 0 {
		header.Set(constant.SpanIDHTTPHeader, c.SpanID)
	}
	if len(c.Token) != 0 {
		header.Set(constant.TokenHTTPHeader, c.Token)
	}
	if len(c.RequestID) != 0 {
		header.Set(constant.RequestIDHTTPHeader, c.RequestID)
	}
	if c.ProjectID != 0 {
		header.Set(constant.ProjectIDHTTPHeader, cast.ToString(c.ProjectID))
	}
}

// GetTraceID description
// @receiver c
// @return string
func (c *Context) GetTraceID() string {
	return c.TraceID
}

// GetProjectID description
// @receiver c
// @return int
func (c *Context) GetProjectID() int {
	return c.ProjectID
}

// WithFields description
// @param c
// @param fields
// @return context.Context
func WithFields(c context.Context, fields map[string]interface{}) context.Context {
	cc := Get(c)
	if cc == nil {
		return context.WithValue(c, ctxKey, &Context{
			logger: logrus.WithFields(fields),
		})
	}
	return context.WithValue(c, ctxKey, &Context{
		logger: cc.logger.WithFields(fields),
	})
}

func getOrWith(ctx context.Context) (context.Context, *Context) {
	cc := Get(ctx)
	if cc == nil {
		cc = &Context{}
		ctx = context.WithValue(ctx, ctxKey, cc)
	}
	return ctx, cc
}

// WithTraceID description
// @param c
// @param traceID
// @return context.Context
func WithTraceID(c context.Context, traceID string) context.Context {
	ctx, cc := getOrWith(c)
	cc.TraceID = traceID
	return ctx
}

// WithProjectID description
// @param c
// @param projectID
// @return context.Context
func WithProjectID(c context.Context, projectID int) context.Context {
	ctx, cc := getOrWith(c)
	cc.ProjectID = projectID
	return ctx
}

// WithGinContext description
// @param c
// @param ginc
// @return context.Context
func WithGinContext(c context.Context, ginc *gin.Context) context.Context {
	ctx, cc := getOrWith(c)
	cc.GinContext = ginc
	return ctx
}

// Get description
// @param c
// @return *Context
func Get(c context.Context) *Context {
	v := c.Value(ctxKey)
	if v == nil {
		if v, ok := c.(*gin.Context); ok {
			if data := v.Value(GINServerContextKey); data != nil {
				out, ok := data.(*Context)
				if ok {
					return out
				}
			}
		}
		return nil
	}
	cc, ok := v.(*Context)
	if !ok {
		return nil
	}
	return cc
}

func GinWith(c *gin.Context, cc *Context) {
	c.Set(GINServerContextKey, cc)
}

// GetTraceID description
// @param c
// @return string
func GetTraceID(c context.Context) string {
	cc := Get(c)
	if cc != nil {
		return cc.GetTraceID()
	}
	return ""
}

// GetProjectID description
// @param c
// @return int
func GetProjectID(c context.Context) int {
	cc := Get(c)
	if cc != nil {
		return cc.GetProjectID()
	}
	return -1
}

// GetGinContext description
// @param c
// @return *gin.Context
func GetGinContext(c context.Context) *gin.Context {
	cc := Get(c)
	if cc != nil {
		return cc.GinContext
	}
	return nil
}
