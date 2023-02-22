package api

import (
	"net/http/pprof"
	"path"

	"github.com/gin-gonic/gin"
)

func RegisterProfile(rootPath string) {
	h := &profileHandler{rootPath: rootPath}
	RegisterHttpRouter(h)
}

type profileHandler struct {
	rootPath string
}

func (p *profileHandler) Use(engine *gin.Engine) {
	engine.Any(path.Join(p.rootPath, "/profile/pprof"), func(c *gin.Context) {
		pprof.Index(c.Writer, c.Request)
	})
	engine.Any(path.Join(p.rootPath, "/profile/pprof/cmdline"), func(c *gin.Context) {
		pprof.Cmdline(c.Writer, c.Request)
	})
	engine.Any(path.Join(p.rootPath, "/profile/pprof/profile"), func(c *gin.Context) {
		pprof.Profile(c.Writer, c.Request)
	})
	engine.Any(path.Join(p.rootPath, "/profile/pprof/symbol"), func(c *gin.Context) {
		pprof.Symbol(c.Writer, c.Request)
	})
	engine.Any(path.Join(p.rootPath, "/profile/pprof/trace"), func(c *gin.Context) {
		pprof.Trace(c.Writer, c.Request)
	})
}
