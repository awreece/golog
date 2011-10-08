package golog

import (
	"flag"
	"fmt"
	"os"
)

type multiLogOuter struct {
	// TODO Add mutex.
	outers map[string]LogOuter
}

func (l *multiLogOuter) String() string {
	return fmt.Sprint("\"", l.outers, "\"")
}

func (l *multiLogOuter) Set(name string) bool {
	if file, err := os.Create(name); err != nil {
		os.Stderr.WriteString(
			fmt.Sprint("Error opening file for logging", name,
				": ", err))
		return false
	} else {
		l.AddLogOuter(name, NewFileLogOuter(file))
		return true
	}

	panic("Code never reaches here, this mollifies the compiler.")
}

func (l *multiLogOuter) AddLogOuter(key string, outer LogOuter) {
	// TODO Grab mutex.
	l.outers[key] = outer
}

func (l *multiLogOuter) RemoveLogOuter(key string) {
	// TODO Be Go1 compatible. :)
	l.outers[key] = nil, false
}

func (l *multiLogOuter) Output(m *LogMessage) {
	for _, outer := range l.outers {
		outer.Output(m)
	}
}

func (l *multiLogOuter) FailNow() {
	for _, outer := range l.outers {
		outer.FailNow()
	}
}

var defaultLogOuters multiLogOuter

func init() {
	flag.Var(&defaultLogOuters, "vlog.logfile", "Log to given file - can "+
		"be provided multiple times to log to multiple files")
}
