package golog

import (
	"flag"
)

// Prints everything this level and above. (Set to SILENT to disable).
var flag_minloglevel = flag.Int("vlog.minloglevel", INFO,
	"Log messages at or above this level. The "+
		"numbers of severity levels INFO, WARNING, "+
		"ERROR, and FATAL are 0, 1, 2, and 3, respectively")

// Logger.Log uses the level to determine whether or not to output the
// arguments. Logger.Log will output the provided arguments exactly, without
// additional formatting such as adding a prefix etc. In addition, Logger.Log
// must be thread safe. The FailNow() function flushes the Logger and performs
// some action. The action performed by FailNow() is deliberately unspecified, 
// but could include os.Exit(1) or testing.(*T).FailNow(), etc.
type Logger interface {
	// If the message is to be logged, evaluates the closure and outputs
	// the result.
	Log(level int, closure func() *LogMessage)
	FailNow()
}

var DefaultLogger Logger = &loggerImpl{
	defaultLogOuters,
	flag_minloglevel,
	exitNow,
}

func NewLogger(outer LogOuter, minloglevel int, failFunc func()) Logger {
	return &loggerImpl{outer, &minloglevel, failFunc}
}

func NewDefaultLogger() Logger {
	return DefaultLogger
}

type loggerImpl struct {
	LogOuter
	// This is a reference type so we can add log writers before having 
	// parsed flag_minloglevel.
	minloglevel *int
	failFunc    func()
}

func (l *loggerImpl) Log(level int, closure func() *LogMessage) {
	if level >= *l.minloglevel {
		l.Output(closure())
	}
}

func (l *loggerImpl) FailNow() {
	// TODO Flush log outer?
	l.failFunc()
}
