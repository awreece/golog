// The vlog interface enables varied levels of logging. This package was
// designed to emulate the google-glog library for C++ logging.

// Usage:
//
// Examples:
//    vlog.Infof("Invalid foo %d", foo)
//	  vlog.Vlogf(4, "Invalid read %v from pipe", err)
//    vlog.Fatal("Invalid read from channel inputBytes")
// 	  vlog.Infoc(func() string { return fmt.Sprint(foo, bar) })
//
// Use the minloglevel flag to control the logging threshold (see flag
// description). Use the v flag to control the default logging threshold for 
// Vlog logging (again, see flag description). The vmodule flag sets logging
// thresholds for vlog logging on a per-package level.
//
// This package provides wrappers for the existing log.Logger library. It
// extends it in two ways:
// 	1) It provides the ability to log to INFO, WARNING, ERROR, or FATAL and
//	   toggle a global threshold for logging. Logging to a level will produce
//	   a log message of the form "{LEVEL} {TIME} {FILE:line}: {message}", for 
//     example: "INFO 00:36:15 vlog.go:24: Hello, world". In addition, an
//     attempt is made to minimize the amount of logic executed in the case
//     where not log statement would be produced (this is done via using the
//	   emptyOuter interface, which is a no-op and thus doesn't render the
//     log string and might possible be optimized out by a compiler).
//
//  2) It gives the user the ability to provide his/her own levels of logging,
//     using the VLOG(n) call and the -v flag. It uses similar logic to the 
//     previous feature, although the threshold is reversed: a call to VLOG(n)
//     will only produce output if n is *below* the global threshold, rather
//     than above.

package vlog

import (
	"fmt"
)

// TODO: Cite http://google-glog.googlecode.com/svn/trunk/doc/glog.html
// 
// With VLOG, the lower the verbose level, the more likely messages are to be 
// logged. For example, if --v==1, VLOG(1) will log, but VLOG(2) will not log. 
// This is opposite of the severity level, where INFO is 0, and ERROR is 2. 
// --minloglevel of 1 will log WARNING and above. Though you can specify any 
// integers for both VLOG macro and --v flag, the common values for them are 
// small positive integers. For example, if you write VLOG(0), you should specify
// --v=-1 or lower to silence it. This is less useful since we may not want 
// verbose logs by default in most cases. The VLOG macros always log at the INFO
//  log level (when they log at all).
//
// Verbose logging can be controlled from the command line on a per-module basis:
//
//   --vmodule=mapreduce=2,file=1,gfs*=3 --v=0
// will:
// 
// a. Print VLOG(2) and lower messages from mapreduce.{h,cc}
// b. Print VLOG(1) and lower messages from file.{h,cc}
// c. Print VLOG(3) and lower messages from files prefixed with "gfs"
// d. Print VLOG(0) and lower messages from elsewhere

func Vlog(level int, m ...interface{}) {
	c := func() string { return fmt.Sprint(m...) }

	// TODO get caller package and time here!
	internalVLogc(level, c, 1)
}

func Vlogf(level int, f string, v ...interface{}) {
	c := func() string { return fmt.Sprintf(f, v...) }

	// TODO get caller package and time here!
	internalVLogc(level, c, 1)
}

func Vlogc(level int, c func() string) {
	// TODO get caller package and time here!
	internalVLogc(level, c, 1)
}

type LogFunc func(message ...interface{})
type LogfFunc func(fmt string, vals ...interface{})
type LogcFunc func(c func() string)

func makeLogFunc(level int) LogFunc {
	return func(v ...interface{}) {
		c := func() string { return fmt.Sprint(v...) }

		// TODO Get time here!
		internalLogc(level, c, 1)
	}
}

func makeLogfFunc(level int) LogfFunc {
	return func(f string, v ...interface{}) {
		c := func() string { return fmt.Sprintf(f, v...) }

		// TODO Get time here!
		internalLogc(level, c, 1)
	}
}

func makeLogcFunc(level int) LogcFunc {
	return func(c func() string) {
		// TODO Get time here!
		internalLogc(level, c, 1)
	}
}

var Info LogFunc = makeLogFunc(INFO)
var Infof LogfFunc = makeLogfFunc(INFO)
var Infoc LogcFunc = makeLogcFunc(INFO)
var Warning LogFunc = makeLogFunc(WARNING)
var Warningf LogfFunc = makeLogfFunc(WARNING)
var Warningc LogcFunc = makeLogcFunc(WARNING)
var Error LogFunc = makeLogFunc(ERROR)
var Errorf LogfFunc = makeLogfFunc(ERROR)
var Errorc LogcFunc = makeLogcFunc(ERROR)
var Fatal LogFunc = makeLogFunc(FATAL)
var Fatalf LogfFunc = makeLogfFunc(FATAL)
var Fatalc LogcFunc = makeLogcFunc(FATAL)
