package api

import "github.com/gin-gonic/gin"

// GinHandleFunc handle func
//  @param *gin.Context
type GinHandleFunc func(*gin.Context)

// HandlerConfig config
type HandlerConfig struct {
	preInterceptors   []GinHandleFunc
	afterInterceptors []GinHandleFunc
}

// HandlerOption opt
//  @param *HandlerConfig
type HandlerOption func(*HandlerConfig)

// HandlerWithPreProcessor pre
//  @param h
//  @return HandlerOption
func HandlerWithPreProcessor(h GinHandleFunc) HandlerOption {
	return func(conf *HandlerConfig) {
		conf.preInterceptors = append(conf.preInterceptors, h)
	}
}

// HandlerWithAfterProcessor after
//  @param h
//  @return HandlerOption
func HandlerWithAfterProcessor(h GinHandleFunc) HandlerOption {
	return func(conf *HandlerConfig) {
		conf.preInterceptors = append(conf.afterInterceptors, h)
	}
}

// DoPreInterceptors do pre
//  @receiver h
//  @param ctx
func (h HandlerConfig) DoPreInterceptors(ctx *gin.Context) {
	for _, f := range h.preInterceptors {
		f(ctx)
	}
}

// DoAfterInterceptors do after
//  @receiver h
//  @param ctx
func (h HandlerConfig) DoAfterInterceptors(ctx *gin.Context) {
	for _, f := range h.afterInterceptors {
		f(ctx)
	}
}
