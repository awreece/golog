// This file definse the logOuter interface and several types of logOuter.
//
// emptyOuter = logOuter where both Out and Outf are noops
// lineOuter = logOuter where a newline is inserted after every call to
//			   Out and Outf
// fatalLineOuter = logOuter that logs message with inserted newline then
//					exits with call to os.EXIT(1)

package golog

import (
	"os"
)

type LogOuter interface {
	// Println guarantees a newline and flushes output.
	Println(string)
	FailNow()
}

type fileLogOuter struct {
	// TODO Insert mutex?
	*os.File
}

func (f *fileLogOuter) Println(s string) {
	// TODO Grab mutex?
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
	// TODO Grab mutex?
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
