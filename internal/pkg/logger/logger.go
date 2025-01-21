package logger

import "go.uber.org/zap"

//go:generate mockgen -source=logger.go -package=mock_logger -destination=../../mocks/pkg/logger/mock_logger.go
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger zap.SugaredLogger
}

func New() *Logger {
	zap, err := zap.NewDevelopment()

	if err != nil {
		// handle err
	}

	return &Logger{logger: *zap.Sugar()}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.logger.Debug(message, args)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info(message, args)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn(message, args)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.logger.Error(message, args)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.logger.Fatal(message, args)
}