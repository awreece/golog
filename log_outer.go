// This file definse the logOuter interface and several types of logOuter.
//
// emptyOuter = logOuter where both Out and Outf are noops
// lineOuter = logOuter where a newline is inserted after every call to
//			   Out and Outf
// fatalLineOuter = logOuter that logs message with inserted newline then
//					exits with call to os.EXIT(1)

package golog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// The location in the source of a log message. Any fields set to their zero
// values will be assumed to be absent. (The empty string is not a valid
// package/function/file, and 0 is not a valid line number).
type LogLocation struct {
	Package  string
	Function string
	File     string
	Line     int
}

type LogMessage struct {
	Level       int
	Nanoseconds int64
	Message     string
	Location    *LogLocation
}

type LogOuter interface {
	// Output a LogMessage (to a file, to stderr, to a tester, etc). Output
	// must be safe to call from multiple threads.
	Output(*LogMessage)
}

// Render a formatted LogLocation to the buffer.
func renderLogLocation(buf *bytes.Buffer, l *LogLocation) {
	if l == nil {
		return
	}
	packPresent := len(l.Package) > 0
	funcPresent := len(l.Function) > 0
	filePresent := len(l.Function) > 0
	linePresent := l.Line > 0

	// TODO(awreece) This logic is terrifying.
	if packPresent {
		buf.WriteString(l.Package)
	}
	if funcPresent {
		if packPresent {
			buf.WriteString(".")
		}
		buf.WriteString(l.Function)
	}
	if (packPresent || funcPresent) && (filePresent || linePresent) {
		buf.WriteString("/")
	}
	if filePresent {
		buf.WriteString(l.File)
	}
	if linePresent {
		if filePresent {
			buf.WriteString(":")
		}
		buf.WriteString(strconv.Itoa(l.Line))
	}
}

func formatLogMessage(m *LogMessage, insertNewline bool) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("L%d", m.Level))
	t := time.NanosecondsToLocalTime(m.Nanoseconds)
	buf.WriteString(t.Format(" 15:04:05.000000"))
	if m.Location != nil {
		buf.WriteString(" ")
		renderLogLocation(&buf, m.Location)
	}
	buf.WriteString("] ")
	buf.WriteString(m.Message)
	if insertNewline {
		buf.WriteString("\n")
	}
	return buf.String()
}

type writerLogOuter struct {
	// TODO Insert mutex?
	io.Writer
}

func (f *writerLogOuter) Output(m *LogMessage) {
	// TODO Grab mutex?
	// Make sure to insert a newline.
	f.Write([]byte(formatLogMessage(m, true)))
}

func NewWriterLogOuter(f io.Writer) LogOuter {
	return &writerLogOuter{f}
}

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

func NewTestLogOuter(t TestController) LogOuter {
	return &testLogOuter{t}
}
