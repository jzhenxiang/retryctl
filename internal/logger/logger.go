package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents log severity.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Entry is a single structured log record.
type Entry struct {
	Timestamp string `json:"timestamp"`
	Level     Level  `json:"level"`
	Message   string `json:"message"`
	Fields    map[string]any `json:"fields,omitempty"`
}

// Logger writes structured JSON log entries.
type Logger struct {
	out io.Writer
}

// New returns a Logger writing to w. If w is nil, os.Stderr is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{out: w}
}

// log writes an entry at the given level.
func (l *Logger) log(level Level, msg string, fields map[string]any) {
	entry := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(l.out, `{"level":"error","message":"failed to marshal log entry"}\n`)
		return
	}
	fmt.Fprintf(l.out, "%s\n", data)
}

// Info logs an informational message.
func (l *Logger) Info(msg string, fields map[string]any) {
	l.log(LevelInfo, msg, fields)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, fields map[string]any) {
	l.log(LevelWarn, msg, fields)
}

// Error logs an error message.
func (l *Logger) Error(msg string, fields map[string]any) {
	l.log(LevelError, msg, fields)
}
