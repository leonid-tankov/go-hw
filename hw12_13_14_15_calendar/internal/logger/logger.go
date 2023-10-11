package logger

import (
	"fmt"
	"io"
	"time"
)

const (
	FATAL = iota
	ERROR
	INFO
	DEBUG
)

type Logger struct {
	level  int
	output io.Writer
}

func New(level string, output io.Writer) *Logger {
	return &Logger{
		level:  getLevel(level),
		output: output,
	}
}

func (l Logger) Fatal(msg string, a ...any) {
	if l.level >= FATAL {
		panic(msg)
	}
}

func (l Logger) Error(msg string, a ...any) {
	if l.level >= ERROR {
		l.log("ERROR", msg, a...)
	}
}

func (l Logger) Info(msg string, a ...any) {
	if l.level >= INFO {
		l.log("INFO", msg, a...)
	}
}

func (l Logger) Debug(msg string, a ...any) {
	if l.level >= DEBUG {
		l.log("DEBUG", msg, a...)
	}
}

func (l Logger) log(level, msg string, a ...any) {
	log := fmt.Sprintf("%s: %s\n", fmt.Sprintf("%s [%s]", time.Now().Format(time.RFC3339), level), fmt.Sprintf(msg, a...))
	l.output.Write([]byte(log))
}

func getLevel(level string) int {
	switch level {
	case "FATAL":
		return 0
	case "ERROR":
		return 1
	case "INFO":
		return 2
	case "DEBUG":
		return 3
	default:
		panic(fmt.Sprintf("Invalid log level: %s", level))
	}
}
