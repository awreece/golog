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
	"strconv"
	"time"
)

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
	Output(*LogMessage)
}

func renderLogLocation(buf *bytes.Buffer, l *LogLocation) {
	packPresent := len(l.Package) > 0
	funcPresent := len(l.Function) > 0
	filePresent := len(l.Function) > 0
	linePresent := l.Line > 0

	// TODO This logic is terrifying.
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

func NewWriterLogOuter(f io.WriteCloser) LogOuter {
	return &writerLogOuter{f}
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
