package api

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/infra"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/servercontext"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/LSDXXX/libs/pkg/wsmanager"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	stdgrpc "google.golang.org/grpc"
)

var (
	// DummyServer fake
	DummyServer    = "http://127.0.0.1:80"
	routers        []HttpRouter
	streamHandlers []StreamMessageHandler
	grpcServices   []GrpcService
	discovery      infra.Discovery
	onceInit       sync.Once
)

func init() {
	// http.DefaultClient.Timeout = constant.RPCTimeout
}

// StreamMessageHandler handler
type StreamMessageHandler interface {
	Topic() string
	Process(message []byte) error
	GroupID() string
}

// HttpRouter router
type HttpRouter interface {
	Use(*gin.Engine)
}

// GrpcService service
type GrpcService interface {
	Use(*stdgrpc.Server)
}

// RegisterSteamMessageHandler register
//
//	@param handler
func RegisterSteamMessageHandler(handler StreamMessageHandler) {
	streamHandlers = append(streamHandlers, handler)
}

// GetStreamMessageHandler get
//
//	@return []StreamMessageHandler
func GetStreamMessageHandler() []StreamMessageHandler {
	return streamHandlers
}

// RegisterHttpRouter register
//
//	@param router
func RegisterHttpRouter(router HttpRouter) {
	routers = append(routers, router)
}

// GetHttpRouters get
//
//	@return []HttpRouter
func GetHttpRouters() []HttpRouter {
	return routers
}

// RegisterGrpcService register
//
//	@param service
func RegisterGrpcService(service GrpcService) {
	grpcServices = append(grpcServices, service)
}

// GetGrpcServices get
//
//	@return []GrpcService
func GetGrpcServices() []GrpcService {
	return grpcServices
}

// DefaultEngine default gin engine
//
//	@param ginLog
//	@return *gin.Engine
func DefaultEngine(ginLog *config.LogConfig) *gin.Engine {
	logger := log.NewLogger(ginLog)
	// gin.DefaultWriter = io.MultiWriter(logger.Out)
	// r := gin.Default()
	// logger.Debug("gin log init")
	r := gin.New()
	// r.Use(gin.RecoveryWithWriter(io.MultiWriter(log.WithContext(context.Background()).Writer())))
	if ginLog.Level == "debug" {
		r.Use(ginLogger(logger))
	}

	r.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.WithContext(context.Background()).Errorf("%+v\n %s", r, util.PrintStack())
				panic(r)
			}
		}()
		c.Next()
	})

	r.Use(ginTrace)

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code":    0,
			"message": "ok",
		})
	})
	return r
}

// WithDiscovery opt
//
//	@param name
//	@return context.Context
//	@return *http.Request
//	@return func(context.Context, *http.Request) error
func WithDiscovery(name string) func(context.Context, *http.Request) error {
	onceInit.Do(func() {
		util.PanicWhenError(container.Resolve(&discovery))
	})
	return func(ctx context.Context, req *http.Request) error {
		address, err := discovery.GetAddress(ctx, name)
		if err != nil {
			return errors.Wrap(err, "get address of "+name)
		}
		// log.WithContext(ctx).Debugf("get address %s: %s", name, address.String())
		servercontext.ContextToHTTP(ctx, req)
		req.URL.Host = address.String()
		return nil
	}
}

func getSig(appKey string, content string) string {
	content = url.QueryEscape(content)
	// fmt.Printf("encode all: %v\n", content)
	key := []byte(appKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func InitOpencensus() {
	err := view.Register(
		// Gin (HTTP) stats
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	)
	if err != nil {
		panic(err)
	}

	// trace only 10%  request
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(0.1)})

	// register exporter: jaeger
	// trace.RegisterExporter()

	// register exporter: prometheus
	// view.RegisterExporter()
}

// WithSign description
// @param appKey
// @return context.Context
// @return *http.Request
// @return func(context.Context, *http.Request) error
func WithSign(appKey string) func(context.Context, *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		query := req.URL.RawQuery
		if req.Method == "POST" {
			data, _ := ioutil.ReadAll(req.Body)
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(data))
			if len(data) > 0 {
				sig := getSig(appKey, string(data))
				query = query + "&sig=" + url.QueryEscape(sig)
			}
		}
		req.URL.RawQuery = query
		sig := getSig(appKey, query)
		req.Header.Set("sig", sig)
		return nil
	}
}

// Init init api
//
//	@return error
func Init(srv string) error {
	InitOpencensus()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var conf *config.Config
	util.PanicWhenError(container.Resolve(&conf))

	container.Singleton(wsmanager.New)

	return nil
}
