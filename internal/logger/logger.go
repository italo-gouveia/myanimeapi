// internal/logger/logger.go
// Package logger provides structured logging capabilities for the MyAnimeAPI application.
// It uses the standard log package as a base but adds structured logging capabilities
// with fields and levels.
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel string

const (
	// LogLevelDebug represents debug level logs
	LogLevelDebug LogLevel = "DEBUG"
	// LogLevelInfo represents info level logs
	LogLevelInfo LogLevel = "INFO"
	// LogLevelWarning represents warning level logs
	LogLevelWarning LogLevel = "WARNING"
	// LogLevelError represents error level logs
	LogLevelError LogLevel = "ERROR"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Logger provides structured logging capabilities
type Logger struct {
	*log.Logger
	fields map[string]interface{}
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
		fields: make(map[string]interface{}),
	}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := &Logger{
		Logger: l.Logger,
		fields: make(map[string]interface{}),
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	newLogger.fields[key] = value
	return newLogger
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := &Logger{
		Logger: l.Logger,
		fields: make(map[string]interface{}),
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// log writes a structured log entry
func (l *Logger) log(level LogLevel, message string) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Fields:    l.fields,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		l.Logger.Printf("Failed to marshal log entry: %v", err)
		return
	}

	l.Logger.Println(string(jsonBytes))
}

// Debug logs a debug level message
func (l *Logger) Debug(message string) {
	l.log(LogLevelDebug, message)
}

// Info logs an info level message
func (l *Logger) Info(message string) {
	l.log(LogLevelInfo, message)
}

// Warning logs a warning level message
func (l *Logger) Warning(message string) {
	l.log(LogLevelWarning, message)
}

// Error logs an error level message
func (l *Logger) Error(message string) {
	l.log(LogLevelError, message)
}

// Debugf logs a formatted debug level message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info level message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warningf logs a formatted warning level message
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warning(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error level message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Default logger instance
var defaultLogger = New()

// SetDefaultLogger sets the default logger instance
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}

// WithField adds a field to the default logger
func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

// WithFields adds multiple fields to the default logger
func WithFields(fields map[string]interface{}) *Logger {
	return defaultLogger.WithFields(fields)
}

// Debug logs a debug level message using the default logger
func Debug(message string) {
	defaultLogger.Debug(message)
}

// Info logs an info level message using the default logger
func Info(message string) {
	defaultLogger.Info(message)
}

// Warning logs a warning level message using the default logger
func Warning(message string) {
	defaultLogger.Warning(message)
}

// Error logs an error level message using the default logger
func Error(message string) {
	defaultLogger.Error(message)
}

// Debugf logs a formatted debug level message using the default logger
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof logs a formatted info level message using the default logger
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warningf logs a formatted warning level message using the default logger
func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

// Errorf logs a formatted error level message using the default logger
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}
