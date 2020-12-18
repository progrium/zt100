package logging

import (
	"unsafe"

	"go.uber.org/zap"
)

func Debug(log DebugLogger, args ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Debug(args...)
	}
}

func Debugf(log DebugLogger, template string, args ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Debugf(template, args...)
	}
}

func Debugw(log DebugLogger, msg string, keysAndValues ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Debugw(msg, keysAndValues...)
	}
}

func Info(log InfoLogger, args ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Info(args...)
	}
}

func Infof(log InfoLogger, template string, args ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Infof(template, args...)
	}
}

func Infow(log InfoLogger, msg string, keysAndValues ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Infow(msg, keysAndValues...)
	}
}

func Error(log ErrorLogger, args ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Error(args...)
	}
}

func Errorf(log ErrorLogger, template string, args ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Errorf(template, args...)
	}
}

func Errorw(log ErrorLogger, msg string, keysAndValues ...interface{}) {
	if !isNilValue(log) {
		withCallerSkip(log).Errorw(msg, keysAndValues...)
	}
}

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

func withCallerSkip(log interface{}) Logger {
	zlog, ok := log.(withOptioner)
	if !ok {
		return log.(Logger)
	}
	return zlog.WithOptions(zap.AddCallerSkip(1))
}

type withOptioner interface {
	WithOptions(opts ...zap.Option) Logger
}
