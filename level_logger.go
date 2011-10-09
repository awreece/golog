package golog

import (
	"fmt"
	"time"
)

const (
	INFO = iota
	WARNING
	ERROR
	FATAL
)

type LevelLogger interface {
	Log(level int, vals ...interface{})
	Logf(level int, f string, vals ...interface{})
	Logc(level int, closure func() string)
}

type levelLoggerImpl struct {
	FailLogger
}

func NewLevelLogger(f FailLogger) LevelLogger {
	return &levelLoggerImpl{f}
}

var DefaultLevelLogger LevelLogger = &levelLoggerImpl{DefaultLogger}

func makeLogClosure(level int, msg func() string) func() *LogMessage {
	// Evaluate this early.
	ns := time.Nanoseconds()

	return func() *LogMessage {
		return &LogMessage{
			Level:       level,
			Message:     msg(),
			Nanoseconds: ns,
		}
	}
}

func (l *levelLoggerImpl) logCommon(level int, closure func() string) {
	l.Log(level, makeLogClosure(level, closure))
}

func (l *levelLoggerImpl) Log(level int, msg ...interface{}) {
	l.logCommon(level, func() string { return fmt.Sprint(msg...) })
}

func (l *levelLoggerImpl) Logf(level int, f string, msg ...interface{}) {
	l.logCommon(level, func() string { return fmt.Sprintf(f, msg...) })
}

func (l *levelLoggerImpl) Logc(level int, closure func() string) {
	l.logCommon(level, closure)
}
