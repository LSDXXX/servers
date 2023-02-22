package api

import (
	"time"

	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/servercontext"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
)

func ginTrace(c *gin.Context) {
	ctx := servercontext.ExtractFromHTTP(c.Request.Context(), c)
	ochttp.SetRoute(ctx, c.Request.URL.Path)
	ctx = servercontext.WithGinContext(ctx, c)
	c.Request = c.Request.WithContext(ctx)
	cc := servercontext.Get(ctx)
	servercontext.GinWith(c, cc)

	cc.WriteHeader(c.Writer.Header())
	c.Next()
}

type LogParams struct {
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration

	MillLatency int64
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
}

func ginLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		param := LogParams{}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.MillLatency = param.Latency.Milliseconds()

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path
		logger.Debugf("%s", log.JsonField(param))
	}

}
