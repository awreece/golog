// This file definse the logOuter interface and several types of logOuter.
//
// emptyOuter = logOuter where both Out and Outf are noops
// lineOuter = logOuter where a newline is inserted after every call to
//			   Out and Outf
// fatalLineOuter = logOuter that logs message with inserted newline then
//					exits with call to os.EXIT(1)

package golog

import (
	"io"
	"os"
	"sync"
)

type LogOuter interface {
	// Output a LogMessage (to a file, to stderr, to a tester, etc). Output
	// must be safe to call from multiple threads.
	Output(*LogMessage)
}

type writerLogOuter struct {
	lock sync.Mutex
	io.Writer
}

func (f *writerLogOuter) Output(m *LogMessage) {
	f.lock.Lock()
	defer f.lock.Unlock()

	// TODO(awreece) Handle short write?
	// Make sure to insert a newline.
	f.Write([]byte(formatLogMessage(m, true)))
}

// Returns a LogOuter wrapping the io.Writer.
func NewWriterLogOuter(f io.Writer) LogOuter {
	return &writerLogOuter{io.Writer: f}
}

// Returns a LogOuter wrapping the file, or an error if the file cannot be
// opened.
func NewFileLogOuter(filename string) (LogOuter, os.Error) {
	if file, err := os.Create(filename); err != nil {
		return nil, err
	} else {
		return NewWriterLogOuter(file), nil
	}

	panic("Code never reaches here, this mollifies the compiler.")
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

func (t *testLogOuter) Output(m *LogMessage) {
	// Don't insert an additional log message since the tester inserts them
	// for us.
	t.Log(formatLogMessage(m, false))
}

// Return a LogOuter wrapping the TestControlller.
func NewTestLogOuter(t TestController) LogOuter {
	return &testLogOuter{t}
}
