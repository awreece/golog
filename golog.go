// Better logging for Go.
package golog

import (
	"flag"
)

var Global *PackageLogger

func init() {
	Global = &PackageLogger{
		failFunc: ExitError,
		outer:    NewDefaultMultiLogOuter(),
	}

	logger := NewLogger(
		Global.outer,
		defaultMinLogLevel,
		func() { Global.failFunc() })

	// Prints everything this level and above. 
	flag.Var(logger, "golog.minloglevel",
		"Log messages at or above this level. The "+
			"numbers of severity levels INFO, WARNING, "+
			"ERROR, and FATAL are 0, 1, 2, and 3, respectively")

	Global.logger = NewLocationLogger(logger, FullLocation)
}

// Wrapper for Global.Info().
func Info(msg ...interface{}) {
	Global.logger.LogDepth(INFO, printClosure(msg...), 1)
}

// Wrapper for Global.Infof().
func Infof(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(INFO, printfClosure(fmt, vals...), 1)
}

// Wrapper for Global.Infoc().
func Infoc(closure func() string) {
	Global.logger.LogDepth(INFO, closure, 1)
}

// Wrapper for Global.Warning().
func Warning(msg ...interface{}) {
	Global.logger.LogDepth(WARNING, printClosure(msg...), 1)
}

// Wrapper for Global.Warningf().
func Warningf(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(WARNING, printfClosure(fmt, vals...), 1)
}

// Wrapper for Global.Warningc().
func Warningc(closure func() string) {
	Global.logger.LogDepth(WARNING, closure, 1)
}

// Wrapper for Global.Error().
func Error(msg ...interface{}) {
	Global.logger.LogDepth(ERROR, printClosure(msg...), 1)
}

// Wrapper for Global.Errorf().
func Errorf(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(ERROR, printfClosure(fmt, vals...), 1)
}

// Wrapper for Global.Errorc().
func Errorc(closure func() string) {
	Global.logger.LogDepth(ERROR, closure, 1)
}

// Wrapper for Global.Fatal().
func Fatal(msg ...interface{}) {
	Global.logger.LogDepth(FATAL, printClosure(msg...), 1)
	Global.logger.FailNow()
}

// Wrapper for Global.Fatalf().
func Fatalf(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(FATAL, printfClosure(fmt, vals...), 1)
	Global.logger.FailNow()
}

// Wrapper for Global.Fatalc().
func Fatalc(closure func() string) {
	Global.logger.LogDepth(FATAL, closure, 1)
	Global.logger.FailNow()
}

// Wrapper for Global.StartTestLogging().
func StartTestLogging(t TestController) {
	Global.StartTestLogging(t)
}

// Wrapper for Global.StopTestLogging().
func StopTestLogging() {
	Global.StopTestLogging()
}

// Wrapper for Global.AddLogOuter().
func AddLogOuter(key string, outer LogOuter) {
	Global.AddLogOuter(key, outer)
}

// Wrapper for Global.RemoveLogOuter().
func RemoveLogOuter(key string) {
	Global.RemoveLogOuter(key)
}

// Wrapper for Global.SetMinLogLevel().
func SetMinLogLevel(level int) {
	Global.SetMinLogLevel(level)
}
