package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/pkg/servercontext"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	tf  = "2006-01-02 15:04:05.000"
	ccf = func(frame *runtime.Frame) string {
		file := path.Base(frame.File)
		return fmt.Sprintf(" [%s:%d %s]", file, frame.Line, frame.Function)
	}
	formatter = &nested.Formatter{
		NoColors:              true,
		CallerFirst:           true,
		TimestampFormat:       tf,
		CustomCallerFormatter: ccf,
	}

	formatter0 = &nested.Formatter{
		HideKeys:              true,
		NoColors:              true,
		CallerFirst:           true,
		TimestampFormat:       tf,
		CustomCallerFormatter: ccf,
	}

	rootLogger *log.Logger
	tagLogger  *log.Logger
)

func init() {
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

// WithContext description
// @param ctx
// @return *log.Entry
func WithContext(ctx context.Context) *log.Entry {
	c := servercontext.Get(ctx)
	if c != nil && c.GetLogger() != nil {
		return c.GetLogger()
	}
	var logger *log.Entry
	if rootLogger == nil {
		logger = log.WithContext(ctx)
	} else {
		logger = rootLogger.WithContext(ctx)
	}
	if c == nil {
		return logger
	}

	logger = logger.WithFields(map[string]interface{}{
		"traceId": c.TraceID,
		// "spanId":  c.SpanID,
	})
	c.SetLogger(logger)
	return logger
}

func getLevel(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	}
	return log.DebugLevel
}

// NewLogger description
// @param conf
// @return *log.Logger
func NewLogger(conf *config.LogConfig) *log.Logger {
	logger := log.New()
	if conf.WithCaller {
		logger.SetReportCaller(true)
	}
	if conf.HiddenKey {
		logger.SetFormatter(formatter0)
	} else {
		logger.SetFormatter(formatter)
	}
	if conf.WithStdOut {
		logger.SetOutput(io.MultiWriter(&conf.Output, os.Stdout))
	} else {
		logger.SetOutput(&conf.Output)
	}
	logger.SetLevel(getLevel(conf.Level))
	return logger
}

// InitGlobalLog description
// @param conf
func InitGlobalLog(conf *config.LogConfig) {
	if conf.WithCaller {
		log.SetReportCaller(true)
	}
	log.SetFormatter(formatter)
	log.SetOutput(io.MultiWriter(&conf.Output, os.Stdout))
	log.SetLevel(getLevel(conf.Level))
	rootLogger = NewLogger(conf)
	conf.HiddenKey = true
	tagLogger = NewLogger(conf)
}

type jsonField struct {
	val interface{}
}

func JsonField(v interface{}) jsonField {
	return jsonField{
		val: v,
	}
}

func (j jsonField) String() string {
	out, err := json.Marshal(j.val)
	if err != nil {
		return errors.Wrap(err, "json marshal").Error()
	}
	return string(out)
}

type EntryProxy struct {
	Tag   string
	entry *log.Entry
	once  sync.Once
}

func (p *EntryProxy) withField() {
	p.once.Do(func() {
		p.entry = tagLogger.WithField(p.Tag, p.Tag)
	})
}

func (p *EntryProxy) Info(args ...any) {
	p.withField()
	p.entry.Info(args)
}

func (p *EntryProxy) Infof(format string, args ...any) {
	p.withField()
	p.entry.Infof(format, args)
}

func (p *EntryProxy) Debug(args ...any) {
	p.withField()
	p.entry.Debug(args)
}

func (p *EntryProxy) Debugf(format string, args ...any) {
	p.withField()
	p.entry.Debugf(format, args)
}

func (p *EntryProxy) Error(args ...any) {
	p.withField()
	p.entry.Error(args)
}

func (p *EntryProxy) Errorf(format string, args ...any) {
	p.withField()
	p.entry.Errorf(format, args)
}

func CreateTagLogger(tag string) *EntryProxy {
	entry := EntryProxy{Tag: tag}
	return &entry
}
