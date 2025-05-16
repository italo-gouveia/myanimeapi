package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger with the specified configuration
func Init(level string) {
	log = logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

	// Set output format
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set output
	log.SetOutput(os.Stdout)
}

// Get returns the configured logger instance
func Get() *logrus.Logger {
	if log == nil {
		Init("info")
	}
	return log
}

// WithField adds a field to the log entry
func WithField(key string, value interface{}) *logrus.Entry {
	return Get().WithField(key, value)
}

// WithFields adds multiple fields to the log entry
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Get().WithFields(fields)
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	Get().Debug(args...)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	Get().Info(args...)
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	Get().Warn(args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	Get().Error(args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1
func Fatal(args ...interface{}) {
	Get().Fatal(args...)
}

// Panic logs a message at level Panic then panics
func Panic(args ...interface{}) {
	Get().Panic(args...)
}
