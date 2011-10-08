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
	"goprotobuf.googlecode.com/hg/proto"
	"os"
	"strconv"
	"time"
)

type LogOuter interface {
	Output(*LogMessage)
	// TODO This doesn't belong here.
	FailNow()
}

func formatLogMessage(m *LogMessage, insertNewline bool) string {
	var buf bytes.Buffer
	buf.WriteString(levelStrings[int(proto.GetInt32(m.Level))])
	t := time.NanosecondsToLocalTime(proto.GetInt64(m.Nanoseconds))
	buf.WriteString(t.Format(" 15:04:05.000000"))
	if m.Location != nil {
		buf.WriteString(" ")
		l := *m.Location
		if l.Package != nil {
			buf.WriteString(*l.Package)
		}
		if l.File != nil {
			buf.WriteString(*l.File)
		}
		if l.Function != nil {
			buf.WriteString(*l.Function)
		}
		if l.Line != nil {
			buf.WriteString(strconv.Itoa(
				int(proto.GetInt32(l.Line))))
		}
	}
	buf.WriteString("] ")
	buf.WriteString(proto.GetString(m.Message))
	if insertNewline {
		buf.WriteString("\n")
	}
	return buf.String()
}

type fileLogOuter struct {
	// TODO Insert mutex?
	*os.File
}

func (f *fileLogOuter) Output(m *LogMessage) {
	// TODO Grab mutex?
	// Make sure to insert a newline.
	f.WriteString(formatLogMessage(m, true))
	f.Sync()
}

func (f *fileLogOuter) FailNow() {
	// TODO Grab mutex?
	f.Close()
	os.Exit(1)
}

func NewFileLogOuter(f *os.File) LogOuter {
	return &fileLogOuter{f}
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
