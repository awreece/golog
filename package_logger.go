package golog

import (
	"fmt"
)

const (
	INFO int = iota
	WARNING
	ERROR
	FATAL
)

type StringLogger interface {
	Log(...interface{})
	Logf(string,...interface{})
	Logc(func()string)
	Info(...interface{})
	Infof(string,...interface{})
	Infoc(func()string)
	Warning(...interface{})
	Warningf(string,...interface{})
	Warningc(func()string)
	Error(...interface{})
	Errorf(string,...interface{})
	Errorc(func()string)
	Fatal(...interface{})
	Fatalf(string,...interface{})
	Fatalc(func()string)
}

// A PackageLogger 
type PackageLogger struct {
	LocationLogger
	MultiLogOuter
	failFunc func()
}

func newPackageLoggerCommon(outer MultiLogOuter, minloglevel_flag *int,
failFunc func(), locFunc func(skip int) *LogLocation) *PackageLogger {
	ret := &PackageLogger{failFunc: failFunc, MultiLogOuter: outer}

	ret.LocationLogger = NewLocationLogger(
		&loggerImpl{outer, minloglevel_flag, func() { ret.failFunc() }},
		FullLocation)

	return ret
}

func NewDefaultPackageLogger() *PackageLogger {
	return newPackageLoggerCommon(
		NewDefaultMultiLogOuter(),
		flag_minloglevel,
		ExitError,
		FullLocation)
}

func NewPackageLogger(outer MultiLogOuter, minloglevel int,
failFunc func(), locFunc func(skip int) *LogLocation) *PackageLogger {
	return newPackageLoggerCommon(outer, &minloglevel, failFunc, locFunc)
}

// Associates TestController with a the "testing LogOuter and updates
// l.FailNow() to call t.FailNow().
func (l *PackageLogger) StartTestLogging(t TestController) {
	l.MultiLogOuter.AddLogOuter("testing", NewTestLogOuter(t))
	// TODO(awreece) Save old failFunc so we can restore it properly.
	l.failFunc = func() { t.FailNow() }
}

// Removes the testing logger and restores l.FailNow() to its previous state.
func (l *PackageLogger) StopTestLogging() {
	l.MultiLogOuter.RemoveLogOuter("testing")
	// TODO(awreece) Restored to saved failFunc.
	l.failFunc = ExitError
}

// Returns a closure that formats the message via a call to fmt.Sprint.
func printClosure(msg ...interface{}) func() string {
	return func() string {
		return fmt.Sprint(msg...)
	}
}

// Returns a closure that formats the message via a call to fmt.Sprintf.
func printfClosure(format string, vals ...interface{}) func() string {
	return func() string {
		return fmt.Sprintf(format, vals...)
	}
}

// Log the message at level INFO, only formatting if message will be logged.
func (l *PackageLogger) Info(msg ...interface{}) {
	l.LogDepth(INFO, printClosure(msg...), 1)
}

// Log the message at level INFO, only fomatting if message will be logged.
func (l *PackageLogger) Infof(fmt string, vals ...interface{}) {
	l.LogDepth(INFO, printfClosure(fmt, vals...), 1)
}

// Log the message at level INFO, only evaluating the closure if the message
// will be logged.
func (l *PackageLogger) Infoc(closure func() string) {
	l.LogDepth(INFO, closure, 1)
}

// Log the message at level WARNING, only formatting if message will be logged.
func (l *PackageLogger) Warning(msg ...interface{}) {
	l.LogDepth(WARNING, printClosure(msg...), 1)
}

// Log the message at level WARNING, only fomatting if message will be logged.
func (l *PackageLogger) Warningf(fmt string, vals ...interface{}) {
	l.LogDepth(WARNING, printfClosure(fmt, vals...), 1)
}

// Log the message at level WARNING, only evaluating the closure if the message
// will be logged.
func (l *PackageLogger) Warningc(closure func() string) {
	l.LogDepth(WARNING, closure, 1)
}

// Log the message at level ERROR, only formatting if message will be logged.
func (l *PackageLogger) Error(msg ...interface{}) {
	l.LogDepth(ERROR, printClosure(msg...), 1)
}

// Log the message at level ERROR, only formatting if message will be logged.
func (l *PackageLogger) Errorf(fmt string, vals ...interface{}) {
	l.LogDepth(ERROR, printfClosure(fmt, vals...), 1)
}

// Log the message at level ERROR, only evaluating the closure if the message
// will be logged.
func (l *PackageLogger) Errorc(closure func() string) {
	l.LogDepth(ERROR, closure, 1)
}

// Log the message at level FATAL, only formatting if the message will be
// logged. Also call l.FailNow().
func (l *PackageLogger) Fatal(msg ...interface{}) {
	l.LogDepth(FATAL, printClosure(msg...), 1)
	l.FailNow()
}

// Log the message at level FATAL, only formatting if the message will be
// logged. Also call l.FailNow().
func (l *PackageLogger) Fatalf(fmt string, vals ...interface{}) {
	l.LogDepth(FATAL, printfClosure(fmt, vals...), 1)
	l.FailNow()
}

// Log the message at level FATAL, only evaluating the clousre if the message
// will be logged. Also call l.FailNow().
func (l *PackageLogger) Fatalc(closure func() string) {
	l.LogDepth(FATAL, closure, 1)
	l.FailNow()
}

func (l *PackageLogger) Log(level int, msg ...interface{}) {
	l.LogDepth(level, printClosure(msg...), 1)
}

func (l *PackageLogger) Logf(level int, format string, msg ...interface{}) {
	l.LogDepth(level, printfClosure(format, msg...), 1)
}

func (l *PackageLogger) Logc(level int, closure func() string) {
	l.LogDepth(level, closure, 1)
}
