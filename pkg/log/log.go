package log

import (
	"go.uber.org/zap"
)

var l *zap.Logger

func init() {
	l, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(l)
}

func Errorw(msg string, args ...any) {
	zap.S().Errorw(msg, args...)
}

func Error(args ...any) {
	zap.S().Error(args...)
}

func Errorf(template string, args ...any) {
	zap.S().Errorf(template, args...)
}

func Info(args ...any) {
	zap.S().Info(args...)
}

func Infow(msg string, args ...any) {
	zap.S().Infow(msg, args...)
}

func Infof(template string, args ...any) {
	zap.S().Infof(template, args...)
}

func Fatal(args ...any) {
	zap.S().Fatal(args...)
}
