package zap

import (
	"fmt"
	"io"
	"net/url"
	"time"

	api "github.com/progrium/zt100/pkg/misc/logging"
	"go.uber.org/zap"
)

func NewLogger(w io.WriteCloser, options ...zap.Option) *Logger {
	logger := newLogger(w, options...)
	return &Logger{logger.Sugar()}
}

func NewRedirectedLogger(w io.WriteCloser, options ...zap.Option) (*Logger, func()) {
	logger := newLogger(w, options...)
	undo, err := zap.RedirectStdLogAt(logger, zap.DebugLevel)
	if err != nil {
		panic(err)
	}
	return &Logger{logger.Sugar()}, undo
}

func newLogger(w io.WriteCloser, options ...zap.Option) *zap.Logger {
	sinkName := fmt.Sprintf("logger-%d", time.Now().Unix())
	zap.RegisterSink(sinkName, func(u *url.URL) (zap.Sink, error) {
		return sink{w}, nil
	})
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{fmt.Sprintf("%s://", sinkName)}
	logger, _ := config.Build(options...)
	return logger
}

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) With(args ...interface{}) api.Logger {
	return l.With(args...)
}

func (l *Logger) WithOptions(opts ...zap.Option) api.Logger {
	return &Logger{l.SugaredLogger.Desugar().WithOptions(opts...).Sugar()}
}

type sink struct {
	io.WriteCloser
}

func (w sink) Sync() error {
	return nil
}
