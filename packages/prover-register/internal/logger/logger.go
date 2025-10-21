package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	output *json.Encoder
}

func NewJSONLogger() *Logger {
	return &Logger{
		output: json.NewEncoder(os.Stderr),
	}
}

func (l *Logger) log(level string, msg string, fields ...interface{}) {
	entry := map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": level,
		"msg":   msg,
	}

	// Add fields as key-value pairs
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			entry[key] = fields[i+1]
		}
	}

	l.output.Encode(entry)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log("info", msg, fields...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log("error", msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.log("debug", msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log("warn", msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log("fatal", msg, fields...)
	os.Exit(1)
}

// Printf for compatibility with standard logger interface
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}
