package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/abhinandanchivate/employee-crud-api/internal/config"
)

type LogLevel string

const (
	DebugLevel LogLevel = "DEBUG"
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
	FatalLevel LogLevel = "FATAL"
)

type Logger struct {
	level  LogLevel
	format string
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

func New() *Logger {
	cfg := config.Get()
	return &Logger{
		level:  LogLevel(cfg.Logging.Level),
		format: cfg.Logging.Format,
	}
}

func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DebugLevel: 0,
		InfoLevel:  1,
		WarnLevel:  2,
		ErrorLevel: 3,
		FatalLevel: 4,
	}
	return levels[level] >= levels[l.level]
}

func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}

	if l.format == "json" {
		jsonData, _ := json.Marshal(entry)
		fmt.Fprintln(os.Stdout, string(jsonData))
	} else {
		fmt.Printf("[%s] %s: %s", entry.Timestamp, entry.Level, entry.Message)
		if fields != nil && len(fields) > 0 {
			fmt.Printf(" %v", fields)
		}
		fmt.Println()
	}

	if level == FatalLevel {
		os.Exit(1)
	}
}

func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.log(DebugLevel, message, fields)
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.log(InfoLevel, message, fields)
}

func (l *Logger) Warn(message string, fields map[string]interface{}) {
	l.log(WarnLevel, message, fields)
}

func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.log(ErrorLevel, message, fields)
}

func (l *Logger) Fatal(message string, fields map[string]interface{}) {
	l.log(FatalLevel, message, fields)
}
