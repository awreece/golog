package golog

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var defaultMinLogLevel int = ERROR

var defaultLogger LoggerFlag = NewDefaultLogger()

func init() {
	// Prints everything this level and above. 
	flag.Var(defaultLogger, "golog.minloglevel",
		"Log messages at or above this level. The "+
			"numbers of severity levels INFO, WARNING, "+
			"ERROR, and FATAL are 0, 1, 2, and 3, respectively")
}

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

// A Logger that can be used as a flag to set minloglevel. For example,
//	var packageLogger LoggerFlag = NewDefaultLogger()
//	
//	func init() {
//		flag.Var(defaultLogger, "minloglevel", "Log messages at or "+
//			"above this level")
//	}
type LoggerFlag interface {
	Logger
	flag.Value
}

func ExitError() {
	os.Exit(1)
}

// Construct a new Logger that writes any messages of level minloglevel or
// higher to the given LogOuter. Calls to Logger.FailNow() call the provided
// failFunc closure.
func NewLogger(outer LogOuter, minloglevel int, failFunc func()) LoggerFlag {
	return &loggerImpl{outer, minloglevel, failFunc}
}

// Return a default initialized log outer. 
func NewDefaultLogger() LoggerFlag {
	return &loggerImpl{
		NewDefaultMultiLogOuter(),
		defaultMinLogLevel,
		ExitError,
	}
}

type loggerImpl struct {
	LogOuter
	minloglevel int
	failFunc    func()
}

func (l *loggerImpl) Log(level int, closure func() *LogMessage) {
	if level >= l.minloglevel {
		l.Output(closure())
	}
}

func (l *loggerImpl) FailNow() {
	// TODO Flush log outer?
	l.failFunc()
}

func (l *loggerImpl) SetMinLogLevel(level int) {
	l.minloglevel = level
}

func (l *loggerImpl) Set(val string) bool {
	if ival, err := strconv.Atoi(val); err != nil {
		l.minloglevel = ival
		return true
	} else {
		fmt.Println("Error setting flag: ", err)
	}

	return false
}

func (l *loggerImpl) String() string {
	return fmt.Sprint(l.minloglevel)
}
