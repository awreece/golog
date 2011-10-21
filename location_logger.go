package golog

import (
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
	getMetadata MakeMetadataFunc
}


// Returns a new LocationLogger wrapping the associated logger, and using
// the provided function to generate the metadata. For example:
//	log := NewLocationLogger(NewDefaultLogger(), NoLocation)
func NewLocationLogger(l Logger, metadataFunc MakeMetadataFunc) LocationLogger {
	return &locationLoggerImpl{l, metadataFunc}
}

// Returns a LocationLogger wrapping the DefaultLogger. 
func NewDefaultLocationLogger() LocationLogger {
	return NewLocationLogger(NewDefaultLogger(), FullLocation)
}

func (l *locationLoggerImpl) makeLogClosure(level int, msg func() string, skip int) func() *LogMessage {
	// Evaluate this early.
	ns := time.Nanoseconds()
	// TODO add ns to metadata.
	metadata := l.getMetadata(skip + 1)

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
