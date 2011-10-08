package golog

import (
	"bytes"
	"fmt"
	"time"
)

const (
	INFO = iota
	WARNING
	ERROR
	FATAL
)

var levelStrings []string = []string{"I", "W", "E", "F"}

type LevelLogger struct {
	// TODO Can we get away with just a Logger?
	FailLogger
}

// Formats the message with metadata. The format is: 
// LEVEL HH:MM:SS:NANOSC LOC] MESSAGE
func makeLogClosure(level int, ns int64, location string, msg func() string) func() string {
	return func() string {
		var buf bytes.Buffer
		buf.WriteString(levelStrings[level])
		t := time.NanosecondsToLocalTime(ns)
		buf.WriteString(t.Format(" 15:04:05.000000 "))
		// TODO Write file and line?
		buf.WriteString(location)
		buf.WriteString("] ")
		buf.WriteString(msg())
		return buf.String()
	}
}

func (l *LevelLogger) logCommon(level int, closure func() string) {
	ns := time.Nanoseconds()
	// TODO Get location?
	l.Logc(level, makeLogClosure(level, ns, "", closure))
}

func (l *LevelLogger) Info(vals ...interface{}) {
	l.logCommon(INFO, func() string { return fmt.Sprint(vals...) })
}

func (l *LevelLogger) Infof(f string, args ...interface{}) {
	l.logCommon(INFO, func() string { return fmt.Sprintf(f, args...) })
}

func (l *LevelLogger) Infoc(closure func() string) {
	l.logCommon(INFO, closure)
}

func (l *LevelLogger) Warning(vals ...interface{}) {
	l.logCommon(WARNING, func() string { return fmt.Sprint(vals...) })
}

func (l *LevelLogger) Warningf(f string, args ...interface{}) {
	l.logCommon(WARNING, func() string { return fmt.Sprintf(f, args...) })
}

func (l *LevelLogger) Warningc(closure func() string) {
	l.logCommon(WARNING, closure)
}

func (l *LevelLogger) Error(vals ...interface{}) {
	l.logCommon(ERROR, func() string { return fmt.Sprint(vals...) })
}

func (l *LevelLogger) Errorf(f string, args ...interface{}) {
	l.logCommon(ERROR, func() string { return fmt.Sprintf(f, args...) })
}

func (l *LevelLogger) Errorc(closure func() string) {
	l.logCommon(ERROR, closure)
}

func (l *LevelLogger) Fatal(vals ...interface{}) {
	l.logCommon(FATAL, func() string { return fmt.Sprint(vals...) })
	l.FailNow()
}

func (l *LevelLogger) Fatalf(f string, args ...interface{}) {
	l.logCommon(FATAL, func() string { return fmt.Sprintf(f, args...) })
	l.FailNow()
}

func (l *LevelLogger) Fatalc(closure func() string) {
	l.logCommon(FATAL, closure)
	l.FailNow()
}
