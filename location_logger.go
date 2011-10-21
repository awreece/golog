package golog

import (
	"path"
	"runtime"
	"strconv"
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
	getLocation func(skip int) map[string]string
}

// Return a nil LogLocation.
func NoLocation(skip int) map[string]string { return make(map[string]string) }

// Walks up the stack skip frames and returns the LogLocation for that frame.
// TODO(awreece) Provide a arg to select which fields to produce?
func FullLocation(skip int) map[string]string {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		// TODO add timestamp.
		return make(map[string]string)
	} else {
		// TODO(awreece) Make sure this is compiler agnostic.
		funcParts := strings.SplitN(runtime.FuncForPC(pc).Name(), ".", 2)
		// TODO add timestamp.
		return map[string]string{
			"package":  funcParts[0],
			"file": path.Base(file),
			"function": funcParts[1],
			"line": strconv.Itoa(line),
		}
	}

	panic("Flow never reaches here, this mollifies the compiler")
}

// Returns a new LocationLogger wrapping the associated logger, and using
// the provided function to generate LogLocations. The locFunc should walk
// up the stack skip frames and generate the LogLocation for that function
// call. For example:
//	log := NewLocationLogger(NewDefaultLogger(), NoLocation)
func NewLocationLogger(l Logger, locFunc func(int) map[string]string) LocationLogger {
	return &locationLoggerImpl{l, locFunc}
}

// Returns a LocationLogger wrapping the DefaultLogger. 
func NewDefaultLocationLoger() LocationLogger {
	return NewLocationLogger(NewDefaultLogger(), FullLocation)
}

func (l *locationLoggerImpl) makeLogClosure(level int, msg func() string, skip int) func() *LogMessage {
	// Evaluate this early.
	ns := time.Nanoseconds()
	metadata := l.getLocation(skip + 1)
	// TODO add ns to metadata.

	return func() *LogMessage {
		return &LogMessage{
			Level:       level,
			Message:     msg(),
			Nanoseconds: ns,
			Metadata : metadata,
		}
	}
}

func (l *locationLoggerImpl) LogDepth(level int, closure func() string, depth int) {
	l.Log(level, l.makeLogClosure(level, closure, depth+1))
}
