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
	// Log the message at the level provided, formatting the message as if
	// via a call to fmt.Sprint (only rendering string if the message will
	// be logged).
	Log(level int, msg ...interface{})
	// Log the message at the level provided, formatting the message as if
	// via a call to fmt.Sprintf (only rendering the string if the message
	// will be logged).
	Logf(level int, fmt string, val ...interface{})
	// Log the message at the level provided. Only evaluates the closure if
	// the message will be logged.
	Logc(int, func() string)
	// Log the message at the INFO level, formatting the message as if via
	// a call to fmt.Sprint and only rendering the string if the message
	// will be logged.
	Info(msg ...interface{})
	// Log the message at the INFO level, formatting the message as if via
	// a call to fmt.Sprintf and only rendering the string if the message
	// will be logged.
	Infof(string, ...interface{})
	// Log the message at the INFO level, only evaluating the closure and 
	// rendering the string if the message will be logged.
	Infoc(func() string)
	// Log the message at the WARNING level, formatting the message as if via
	// a call to fmt.Sprint and only rendering the string if the message
	// will be logged.
	Warning(...interface{})
	// Log the message at the WARNING level, formatting the message as if via
	// a call to fmt.Sprintf and only rendering the string if the message
	// will be logged.
	Warningf(string, ...interface{})
	// Log the message at the WARNING level, only evaluating the closure and 
	// rendering the string if the message will be logged.
	Warningc(func() string)
	// Log the message at the ERROR level, formatting the message as if via
	// a call to fmt.Sprint and only rendering the string if the message
	// will be logged.
	Error(...interface{})
	// Log the message at the ERROR level, formatting the message as if via
	// a call to fmt.Sprintf and only rendering the string if the message
	// will be logged.
	Errorf(string, ...interface{})
	// Log the message at the ERROR level, only evaluating the closure and 
	// rendering the string if the message will be logged.
	Errorc(func() string)
	// Log the message at the FATAL level, formatting the message as if via
	// a call to fmt.Sprint and only rendering the string if the message
	// will be logged. Afterwards, calls a LogOuter.FailNow().
	Fatal(...interface{})
	// Log the message at the FATAL level, formatting the message as if via
	// a call to fmt.Sprintf and only rendering the string if the message
	// will be logged. Afterwards, calls a LogOuter.FailNow()
	Fatalf(string, ...interface{})
	// Log the message at the FATAL level, only evaluating the closure and 
	// rendering the string if the message will be logged. Afterwards,
	// calls a LogOuter.FailNow().
	Fatalc(func() string)
}

// A PackageLogger is the expected entry point into this package. It is a
// StringLogger that also provides the StartTestLogging and StopTestLogging,
// as well as the ability to add and remove LogOuters and modify minloglevel.
type PackageLogger struct {
	logger   LocationLogger
	outer    MultiLogOuter
	failFunc func()
}

func NewPackageLogger(outer MultiLogOuter, minloglevel int,
	failFunc func(), metadataFunc MetadataFunc) *PackageLogger {
	ret := &PackageLogger{failFunc: failFunc, outer: outer}

	ret.logger = NewLocationLogger(
		&loggerImpl{
			outer,
			defaultMinLogLevel,
			func() { ret.failFunc() },
		},
		metadataFunc)

	return ret
}

// Construct a new PackageLogger with sensible defaults.
func NewDefaultPackageLogger() *PackageLogger {
	return NewPackageLogger(
		NewDefaultMultiLogOuter(),
		defaultMinLogLevel,
		ExitError,
		MakeMetadataFunc(DefaultMetadata))
}

// Associates TestController with a the "testing LogOuter and updates
// l.FailNow() to call t.FailNow().
func (l *PackageLogger) StartTestLogging(t TestController) {
	l.outer.AddLogOuter("testing", NewTestLogOuter(t))
	// TODO(awreece) Save old failFunc so we can restore it properly.
	l.failFunc = func() { t.FailNow() }
}

// Removes the testing logger and restores l.FailNow() to its previous state.
func (l *PackageLogger) StopTestLogging() {
	l.outer.RemoveLogOuter("testing")
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

// Implements StringLogger.Info().
func (l *PackageLogger) Info(msg ...interface{}) {
	l.logger.LogDepth(INFO, printClosure(msg...), 1)
}

// Implement StringLogger.Infof().
func (l *PackageLogger) Infof(fmt string, vals ...interface{}) {
	l.logger.LogDepth(INFO, printfClosure(fmt, vals...), 1)
}

// Implement StringLogger.Infoc()
func (l *PackageLogger) Infoc(closure func() string) {
	l.logger.LogDepth(INFO, closure, 1)
}

// Implement StringLogger.Warning().
func (l *PackageLogger) Warning(msg ...interface{}) {
	l.logger.LogDepth(WARNING, printClosure(msg...), 1)
}

// Implement StringLogger.Warningf().
func (l *PackageLogger) Warningf(fmt string, vals ...interface{}) {
	l.logger.LogDepth(WARNING, printfClosure(fmt, vals...), 1)
}

// Implement StringLogger.Warningc().
func (l *PackageLogger) Warningc(closure func() string) {
	l.logger.LogDepth(WARNING, closure, 1)
}

// Implement StringLogger.Error().
func (l *PackageLogger) Error(msg ...interface{}) {
	l.logger.LogDepth(ERROR, printClosure(msg...), 1)
}

// Implement StringLogger.Errorf().
func (l *PackageLogger) Errorf(fmt string, vals ...interface{}) {
	l.logger.LogDepth(ERROR, printfClosure(fmt, vals...), 1)
}

// Implement StringLogger.Errorc().
func (l *PackageLogger) Errorc(closure func() string) {
	l.logger.LogDepth(ERROR, closure, 1)
}

// Implement StringLogger.Fatal().
func (l *PackageLogger) Fatal(msg ...interface{}) {
	l.logger.LogDepth(FATAL, printClosure(msg...), 1)
	l.logger.FailNow()
}

// Implement StringLogger.Fatalf()
func (l *PackageLogger) Fatalf(fmt string, vals ...interface{}) {
	l.logger.LogDepth(FATAL, printfClosure(fmt, vals...), 1)
	l.logger.FailNow()
}

// Implement StringLogger.Fatalc()
func (l *PackageLogger) Fatalc(closure func() string) {
	l.logger.LogDepth(FATAL, closure, 1)
	l.logger.FailNow()
}

// Implement StringLogger.Log()
func (l *PackageLogger) Log(level int, msg ...interface{}) {
	l.logger.LogDepth(level, printClosure(msg...), 1)
}

// Implement StringLogger.Logf()
func (l *PackageLogger) Logf(level int, format string, msg ...interface{}) {
	l.logger.LogDepth(level, printfClosure(format, msg...), 1)
}

// Implement StringLogger.Logc()
func (l *PackageLogger) Logc(level int, closure func() string) {
	l.logger.LogDepth(level, closure, 1)
}

// Export MutliLogOuter.AddLogOuter().
func (l *PackageLogger) AddLogOuter(key string, outer LogOuter) {
	l.outer.AddLogOuter(key, outer)
}

// Export MutliLogOuter.RemoveLogOuter().
func (l *PackageLogger) RemoveLogOuter(key string) {
	l.outer.RemoveLogOuter(key)
}

// Export Logger.SetMinLogLevel(). Note - this will affect MinLogLevel of
// underlying Logger.
func (l *PackageLogger) SetMinLogLevel(level int) {
	l.logger.SetMinLogLevel(level)
}
