package golog

import (
	"path"
	"runtime"
	"strings"
	"time"
)

// A LocationLogger wraps useful methods for outputting strings by creating
// a LogMessage with the relevant metadata.
type LocationLogger interface {
	Logger
	LogDepth(level int, closure func() string, depth int)
}

type locationLoggerImpl struct {
	Logger
	// TODO(awreece) comment this
	// Skip 0 refers to the function calling getLocation.
	getLocation func(skip int) *LogLocation
}

// Return a nil LogLocation.
func NoLocation(skip int) *LogLocation { return nil }

// Walks up the stack skip frames and returns the LogLocation for that frame.
// TODO(awreece) Provide a arg to select which fields to produce?
// REVIEW(korfuri) You can generalize further, see my comment in log_outer.go
func FullLocation(skip int) *LogLocation {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	} else {
		// TODO(awreece) Make sure this is compiler agnostic.
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

// Returns a new LocationLogger wrapping the associated logger, and using
// the provided function to generate LogLocations. The locFunc should walk
// up the stack skip frames and generate the LogLocation for that function
// call. For example:
//	log := NewLocationLogger(NewDefaultLogger(), NoLocation)
func NewLocationLogger(l Logger, locFunc func(int) *LogLocation) LocationLogger {
	return &locationLoggerImpl{l, locFunc}
}

// Returns a LocationLogger wrapping the DefaultLogger. 
func NewDefaultLocationLoger() LocationLogger {
	return NewLocationLogger(NewDefaultLogger(), FullLocation)
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
	l.Log(level, l.makeLogClosure(level, closure, depth+1))
}
