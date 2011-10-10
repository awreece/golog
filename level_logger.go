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
	FailNow()
	SetMinLogLevel(int)
}

type LevelLogger interface {
	Log(level int, vals ...interface{})
	Logf(level int, f string, vals ...interface{})
	Logc(level int, closure func() string)
	LocationLogger
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
			Package:  funcParts[0],
			File:     path.Base(file),
			Function: funcParts[1],
			Line:     line,
		}
	}

	panic("Flow never reaches here, this mollifies the compiler")
}

func NewLevelLogger(l Logger, locFunc func(int) *LogLocation) LevelLogger {
	return &levelLoggerImpl{l, locFunc}
}

func (l *levelLoggerImpl) makeLogClosure(level int, msg func() string, skip int) func() *LogMessage {
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

func (l *levelLoggerImpl) LogDepth(level int, closure func() string, depth int) {
	l.Logger.Log(level, l.makeLogClosure(level, closure, depth + 1))
}

func (l *levelLoggerImpl) Log(level int, msg ...interface{}) {
	l.LogDepth(level, func() string { return fmt.Sprint(msg...) }, 1)
}

func (l *levelLoggerImpl) Logf(level int, f string, msg ...interface{}) {
	l.LogDepth(level, func() string { return fmt.Sprintf(f, msg...) }, 1)
}

func (l *levelLoggerImpl) Logc(level int, closure func() string) {
	l.LogDepth(level, closure, 1)
}
