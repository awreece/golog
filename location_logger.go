package golog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type LocationLogger interface {
	LogDepth(level int, closure func() string, depth int)
	Log(int, ...interface{})
	Logf(int, string, ...interface{})
	Logc(int, func() string)
	FailNow()
	SetMinLogLevel(int)
}

type locationLoggerImpl struct {
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
			Package:  funcParts[0],
			File:     path.Base(file),
			Function: funcParts[1],
			Line:     line,
		}
	}

	panic("Flow never reaches here, this mollifies the compiler")
}

func NewLocationLogger(l Logger, locFunc func(int) *LogLocation) LocationLogger {
	return &locationLoggerImpl{l, locFunc}
}

func (l *locationLoggerImpl) makeLogClosure(level int, msg func() string, skip int) func() *LogMessage {
	// Evaluate this early.
	ns := time.Nanoseconds()
	location := l.getLocation(skip + 1)

	return func() *LogMessage {
		return &LogMessage{
			Level:       level,
			Message:     msg(),
			Nanoseconds: ns,
			Location:    location,
		}
	}
}

func (l *locationLoggerImpl) LogDepth(level int, closure func() string, depth int) {
	l.Logger.Log(level, l.makeLogClosure(level, closure, depth+1))
}

func (l *locationLoggerImpl) Log(level int, msg ...interface{}) {
	l.LogDepth(level, func() string { return fmt.Sprint(msg...) }, 1)
}

func (l *locationLoggerImpl) Logf(level int, format string, msg ...interface{}) {
	l.LogDepth(level, func() string { return fmt.Sprintf(format, msg...) }, 1)
}

func (l *locationLoggerImpl) Logc(level int, closure func() string) {
	l.LogDepth(level, closure, 1)
}
