package golog

import (
	"flag"
	"os"
)

// Prints everything this level and above. (Set to SILENT to disable).
var flag_minloglevel = flag.Int("golog.minloglevel", INFO,
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
	// Fail and halt standard control flow.
	FailNow()
	// All future calls to log with log only if the message is at 
	// level or higher.
	// TODO(awreece) Put in different interface to keep Logger agnostic?
	SetMinLogLevel(level int)
}

func ExitError() {
	os.Exit(1)
}

// Construct a new Logger that writes any messages of level minloglevel or
// higher to the given LogOuter. Calls to Logger.FailNow() call the provided
// failFunc closure.
func NewLogger(outer LogOuter, minloglevel int, failFunc func()) Logger {
	return &loggerImpl{outer, &minloglevel, failFunc}
}

// Return the default log outer. 
func NewDefaultLogger() Logger {
	// TODO Use only one global default logger, pass a ptr to it.
	return &loggerImpl{
		NewDefaultMultiLogOuter(),
		flag_minloglevel,
		ExitError,
	}
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

func (l *loggerImpl) SetMinLogLevel(level int) {
	l.minloglevel = &level
}
