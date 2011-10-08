// This file definse the logOuter interface and several types of logOuter.
//
// emptyOuter = logOuter where both Out and Outf are noops
// lineOuter = logOuter where a newline is inserted after every call to
//			   Out and Outf
// fatalLineOuter = logOuter that logs message with inserted newline then
//					exits with call to os.EXIT(1)

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

type LogOuter interface {
	// Println guarantees a newline and flushes output.
	Println(string)
	FailNow()
}

type fileLogOuter struct {
	*os.File
}

func (f *fileLogOuter) Println(s string) {
	l := len(s)
	if l > 0 {
		if s[l-1] == '\n' {
			f.WriteString(s)
		} else {
			f.WriteString(s + "\n")
		}
	}

	f.Sync()
}

func (f *fileLogOuter) FailNow() {
	f.Close()
	os.Exit(1)
}

// We want to allow an abitrary testing framework.
type TestController interface {
	// We will assume that testers insert newlines in manner similar to 
	// the FEATURE of testing.T where it inserts extra newlines. >.<
	Log(...interface{})
	FailNow()
}

type testLogOuter struct {
	TestController
}

func (t* testLogOuter) Println(s string) {
	l := len(s)
	if l > 0 {
		// Since testers insert newlines, we strip the newline
		// in our string.
		if s[l -1] == '\n' {
			t.Log(s[:l-1])
		} else {
			t.Log(s)
		}
	}
}

type internalLogOuter struct {
	LogOuter
	// This is a reference type so we can add log writers before having 
	// parsed flag_minloglevel.
	minloglevel *int
	vmoduleLevelsMap
}

type logOuterList struct{
	// TODO Insert mutex here.
	outers []*internalLogOuter
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
	l.outers = append(l.outers, &internalLogOuter{
		LogOuter: &fileLogOuter{file},
		minloglevel: flag_minloglevel,
		vmoduleLevelsMap: flag_vmodule,
	})
}

func (l *logOuterList) AddDefaultLogTester(t TestController) {
	// TODO Grab Mutex
	l.outers = append(l.outers, &internalLogOuter{
		LogOuter: &testLogOuter{t},
		minloglevel: flag_minloglevel,
		vmoduleLevelsMap: flag_vmodule,
	})
}

var LogOuters logOuterList

func init() {
	flag.Var(&LogOuters, "vlog.logfile", "Log to given file - can be"+
		" provided multiple times to log to multiple files")
}
