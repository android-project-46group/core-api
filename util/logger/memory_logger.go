package logger

import (
	"context"
	"fmt"
)

type memoryLogger struct {
	level Level
	Logs  []string
}

func NewMemoryLogger(level Level) Logger {
	return &memoryLogger{
		level: level,
		Logs:  []string{},
	}
}

func (m *memoryLogger) Critical(ctx context.Context, msg string) {
	m.Print(ctx, Critical, msg)
}

func (m *memoryLogger) Error(ctx context.Context, msg string) {
	m.Print(ctx, Error, msg)
}

func (m *memoryLogger) Warn(ctx context.Context, msg string) {
	m.Print(ctx, Warn, msg)
}

func (m *memoryLogger) Info(ctx context.Context, msg string) {
	m.Print(ctx, Info, msg)
}

func (m *memoryLogger) Debug(ctx context.Context, msg string) {
	m.Print(ctx, Degub, msg)
}

func (m *memoryLogger) Criticalf(ctx context.Context, msg string, a ...interface{}) {
	m.Print(ctx, Critical, fmt.Sprintf(msg, a...))
}

func (m *memoryLogger) Errorf(ctx context.Context, msg string, a ...interface{}) {
	m.Print(ctx, Error, fmt.Sprintf(msg, a...))
}

func (m *memoryLogger) Warnf(ctx context.Context, msg string, a ...interface{}) {
	m.Print(ctx, Warn, fmt.Sprintf(msg, a...))
}

func (m *memoryLogger) Infof(ctx context.Context, msg string, a ...interface{}) {
	m.Print(ctx, Info, fmt.Sprintf(msg, a...))
}

func (m *memoryLogger) Debugf(ctx context.Context, msg string, a ...interface{}) {
	m.Print(ctx, Degub, fmt.Sprintf(msg, a...))
}

func (m *memoryLogger) Print(ctx context.Context, level Level, msg string) {
	if !shouldPrint(m.level, level) {
		return
	}

	m.Logs = append(m.Logs, msg)
}

func (m *memoryLogger) SetLevel(lv Level) {}
