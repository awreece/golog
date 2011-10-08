package golog

import (
	"fmt"
)

type LevelLogger struct {
	// TODO Can we get away with just a Logger?
	FailLogger
}

func (l *LevelLogger) logCommon(level int, closure func() string) {
	// TODO Add prefix, timestamp, etc.
	l.Logc(level, closure)
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
