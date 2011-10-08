package golog

import (
	"flag"
	"fmt"
	"os"
)

// Prints everything this level and above. (Set to SILENT to disable).
var flag_minloglevel = flag.Int("vlog.minloglevel", INFO,
	"Log messages at or above this level. The "+
		"numbers of severity levels INFO, WARNING, "+
		"ERROR, and FATAL are 0, 1, 2, and 3, respectively")

// Logger.Log uses the level to determine whether or not to output the
// arguments. Logger.Log will output the provided arguments exactly, without
// additional formatting such as adding a prefix etc. In addition, Logger.Log
// must be thread safe.
type Logger interface {
	// If the message is to be logged, evaluates the closure and outputs
	// the result.
	Log(level int, closure func() string)
}

// A FailLogger is a Logger with the addition FailNow() function, which flushes
// the Logger and performs some action. The action performed by FailNow() is 
// deliberately unspecified, but could include os.Exit(1) or 
// testing.(*T).FailNow().
type FailLogger interface {
	Logger
	FailNow()
}

type loggerImpl struct {
	LogOuter
	// This is a reference type so we can add log writers before having 
	// parsed flag_minloglevel.
	minloglevel *int
	vmoduleLevelsMap
}

type logOuterList struct {
	// TODO Insert mutex here.
	outers []*loggerImpl
}

func (l *logOuterList) String() string {
	return fmt.Sprint("\"", l.outers, "\"")
}

func (l *logOuterList) Set(name string) bool {
	if file, err := os.Create(name); err != nil {
		os.Stderr.WriteString(
			fmt.Sprint("Error opening file for logging", name,
				": ", err))
		return false
	} else {
		l.AddDefaultLogFile(name, file)
		return true
	}

	panic("Code never reaches here, this mollifies the compiler.")
}

// TODO only require filename.
func (l *logOuterList) AddDefaultLogFile(filename string, file *os.File) {
	// TODO Grab mutex.
	l.outers = append(l.outers, &loggerImpl{
		LogOuter:         &fileLogOuter{file},
		minloglevel:      flag_minloglevel,
		vmoduleLevelsMap: flag_vmodule,
	})
}

func (l *logOuterList) AddDefaultLogTester(t TestController) {
	// TODO Grab Mutex
	l.outers = append(l.outers, &loggerImpl{
		LogOuter:         &testLogOuter{t},
		minloglevel:      flag_minloglevel,
		vmoduleLevelsMap: flag_vmodule,
	})
}

var LogOuters logOuterList

func init() {
	flag.Var(&LogOuters, "vlog.logfile", "Log to given file - can be"+
		" provided multiple times to log to multiple files")
}
