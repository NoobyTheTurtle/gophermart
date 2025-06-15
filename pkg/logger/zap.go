package logger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func New() (*ZapLogger, error) {
	var cfg zap.Config

	cfg = zap.NewDevelopmentConfig()
	cfg.DisableCaller = true

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	sugar := l.Sugar()

	return &ZapLogger{
		logger: sugar,
	}, nil
}

func (z *ZapLogger) Info(message string, args ...any) {
	z.logger.Infow(message, args...)
}

func (z *ZapLogger) Error(message string, args ...any) {
	z.logger.Errorw(message, args...)
}

func (z *ZapLogger) Warn(message string, args ...any) {
	z.logger.Warnw(message, args...)
}

func (z *ZapLogger) Debug(message string, args ...any) {
	z.logger.Debugw(message, args...)
}

func (z *ZapLogger) Fatal(message string, args ...any) {
	z.logger.Fatalw(message, args...)
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}
