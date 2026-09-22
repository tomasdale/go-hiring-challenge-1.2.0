package logger

import "go.uber.org/zap"

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type ZapLogger struct {
	*zap.Logger
}

func NewZapLogger() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{Logger: logger}, nil
}

func NewNopLogger() Logger {
	return &ZapLogger{Logger: zap.NewNop()}
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}
