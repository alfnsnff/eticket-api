package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Debug(args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

type LogrusLogger struct {
	entry *logrus.Entry
}

func NewLogrusLogger(serviceName string) Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{}) // Structured logging
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)

	return &LogrusLogger{
		entry: l.WithField("service", serviceName),
	}
}

func (l *LogrusLogger) Info(args ...interface{})  { l.entry.Info(args...) }
func (l *LogrusLogger) Error(args ...interface{}) { l.entry.Error(args...) }
func (l *LogrusLogger) Warn(args ...interface{})  { l.entry.Warn(args...) }
func (l *LogrusLogger) Debug(args ...interface{}) { l.entry.Debug(args...) }

func (l *LogrusLogger) WithField(k string, v interface{}) Logger {
	return &LogrusLogger{entry: l.entry.WithField(k, v)}
}

func (l *LogrusLogger) WithFields(f map[string]interface{}) Logger {
	return &LogrusLogger{entry: l.entry.WithFields(f)}
}

func (l *LogrusLogger) WithError(err error) Logger {
	return &LogrusLogger{entry: l.entry.WithError(err)}
}
