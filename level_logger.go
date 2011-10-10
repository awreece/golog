package golog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type LevelLogger interface {
	FailNow()
	Log(level int, vals ...interface{})
	Logf(level int, f string, vals ...interface{})
	Logc(level int, closure func() string)
}

type levelLoggerImpl struct {
	Logger
	// TODO comment this
	// Skip 0 refers to the function calling getLocation.
	getLocation func(skip int) *LogLocation
}

func NoLocation(skip int) *LogLocation { return nil }

func FullLocation(skip int) *LogLocation {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	} else {
		// TODO Make sure this is compiler agnostic.
		funcParts := strings.SplitN(runtime.FuncForPC(pc).Name(), ".", 2)
		return &LogLocation{
			Package: funcParts[0],
			File: path.Base(file),
			Function: funcParts[1],
			Line: line,
		}
	}

	panic("Flow never reaches here, this mollifies the compiler")
}

func NewLevelLogger(l Logger, locFunc func(int) *LogLocation) LevelLogger {
	return &levelLoggerImpl{l, locFunc}
}

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
			Location:    location,
		}
	}
}

func (l *levelLoggerImpl) logCommon(level int, closure func() string) {
	l.Logger.Log(level, l.makeLogClosure(level, closure))
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
