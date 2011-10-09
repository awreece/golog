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
	// TODO comment this
	// Skip 0 refers to the function calling getLocation.
	getLocation func(skip int) *LogLocation
}

func NoLocation (skip int) *LogLocation { return nil }

func NewLevelLogger(f FailLogger) LevelLogger {
	return &levelLoggerImpl{f, NoLocation}
}

var DefaultLevelLogger LevelLogger = &levelLoggerImpl{DefaultLogger, NoLocation}

func (l *levelLoggerImpl) makeLogClosure(level int, msg func() string) func() *LogMessage {
	// Evaluate this early.
	ns := time.Nanoseconds()
	// TODO Be less brittle.
	// Skip over makeLogClosure, logCommon, and Log
	location := l.getLocation(3)

	return func() *LogMessage {
		return &LogMessage{
			Level:       level,
			Message:     msg(),
			Nanoseconds: ns,
			Location: location,
		}
	}
}

func (l *levelLoggerImpl) logCommon(level int, closure func() string) {
	l.Log(level, l.makeLogClosure(level, closure))
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
