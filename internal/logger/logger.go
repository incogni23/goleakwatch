package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents log levels
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger interface for pluggable logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
	SetLevel(level Level)
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new Field
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// DefaultLogger is a simple logger implementation
type DefaultLogger struct {
	level  Level
	output io.Writer
	fields []Field
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger(level Level, output io.Writer) *DefaultLogger {
	if output == nil {
		output = os.Stderr
	}

	return &DefaultLogger{
		level:  level,
		output: output,
		fields: make([]Field, 0),
	}
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	l.log(DEBUG, msg, fields...)
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	l.log(INFO, msg, fields...)
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	l.log(WARN, msg, fields...)
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, fields ...Field) {
	l.log(ERROR, msg, fields...)
}

// WithFields returns a new logger with additional fields
func (l *DefaultLogger) WithFields(fields ...Field) Logger {
	newLogger := &DefaultLogger{
		level:  l.level,
		output: l.output,
		fields: make([]Field, len(l.fields)+len(fields)),
	}

	copy(newLogger.fields, l.fields)
	copy(newLogger.fields[len(l.fields):], fields)

	return newLogger
}

// SetLevel sets the log level
func (l *DefaultLogger) SetLevel(level Level) {
	l.level = level
}

// log is the internal logging method
func (l *DefaultLogger) log(level Level, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	levelStr := level.String()

	// Build the log message
	logMsg := "[" + timestamp + "] " + levelStr + ": " + msg

	// Add fields if any
	if len(l.fields) > 0 || len(fields) > 0 {
		logMsg += " | "
		allFields := make([]Field, 0, len(l.fields)+len(fields))
		allFields = append(allFields, l.fields...)
		allFields = append(allFields, fields...)

		for i, field := range allFields {
			if i > 0 {
				logMsg += ", "
			}
			logMsg += fmt.Sprintf("%s=%v", field.Key, field.Value)
		}
	}

	logMsg += "\n"

	fmt.Fprint(l.output, logMsg)
}

// NoOpLogger is a logger that does nothing
type NoOpLogger struct{}

// Debug does nothing
func (l *NoOpLogger) Debug(msg string, fields ...Field) {}

// Info does nothing
func (l *NoOpLogger) Info(msg string, fields ...Field) {}

// Warn does nothing
func (l *NoOpLogger) Warn(msg string, fields ...Field) {}

// Error does nothing
func (l *NoOpLogger) Error(msg string, fields ...Field) {}

// WithFields returns the same logger
func (l *NoOpLogger) WithFields(fields ...Field) Logger {
	return l
}

// SetLevel does nothing
func (l *NoOpLogger) SetLevel(level Level) {}

// Global logger instance
var globalLogger Logger = NewDefaultLogger(INFO, os.Stderr)

// SetGlobalLogger sets the global logger
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger
func GetGlobalLogger() Logger {
	return globalLogger
}

// Debug logs a debug message using the global logger
func Debug(msg string, fields ...Field) {
	globalLogger.Debug(msg, fields...)
}

// Info logs an info message using the global logger
func Info(msg string, fields ...Field) {
	globalLogger.Info(msg, fields...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, fields ...Field) {
	globalLogger.Warn(msg, fields...)
}

// Error logs an error message using the global logger
func Error(msg string, fields ...Field) {
	globalLogger.Error(msg, fields...)
}

// WithFields returns a new logger with additional fields
func WithFields(fields ...Field) Logger {
	return globalLogger.WithFields(fields...)
}
