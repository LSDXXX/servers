package app

import (
	"net/http"

	"github.com/LSDXXX/libs/api"
	"github.com/LSDXXX/libs/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

// HttpServer server
type HttpServer interface {
	Start(port int) error
}

type httpServer struct {
	engine *gin.Engine
}

func (s *httpServer) Start(port int) error {
	return http.ListenAndServe(":"+cast.ToString(port),
		&ochttp.Handler{
			Handler: s.engine,
			GetStartOptions: func(r *http.Request) trace.StartOptions {
				startOptions := trace.StartOptions{}
				if r.URL.Path == "/metrics" {
					startOptions.Sampler = trace.NeverSample()
				}
				return startOptions
			},
		},
	)
	// return s.engine.Run(":" + cast.ToString(port))
}

// NewHttpServer new
//  @param ginLog
//  @param routers
//  @return HttpServer
func NewHttpServer(ginLog *config.LogConfig, routers ...api.HttpRouter) HttpServer {
	ginLog.WithStdOut = false
	engine := api.DefaultEngine(ginLog)
	engine.RedirectTrailingSlash = false
	for _, router := range routers {
		router.Use(engine)
	}
	return &httpServer{
		engine: engine,
	}
}
