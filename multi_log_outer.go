package golog

import (
	"flag"
	"fmt"
	"os"
)

type MultiLogOuter struct {
	// TODO Add mutex.
	outers map[string]LogOuter
}

func (l *MultiLogOuter) String() string {
	return fmt.Sprint("\"", l.outers, "\"")
}

func (l *MultiLogOuter) Set(name string) bool {
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

func (l *MultiLogOuter) AddLogOuter(key string, outer LogOuter) {
	// TODO Grab mutex.
	l.outers[key] = outer
}

func (l *MultiLogOuter) RemoveLogOuter(key string) {
	// TODO Be Go1 compatible. :)
	l.outers[key] = nil, false
}

func (l *MultiLogOuter) Output(m *LogMessage) {
	// TODO Grab mutex.
	for _, outer := range l.outers {
		outer.Output(m)
	}
}

func (l *MultiLogOuter) FailNow() {
	// TODO Grab mutex.
	// TODO: This doesn't belong here. We want a MuliFailLogger.
	for _, outer := range l.outers {
		outer.FailNow()
	}
}

var defaultLogOuters MultiLogOuter = MultiLogOuter{make(map[string]LogOuter)}

func init() {
	flag.Var(&defaultLogOuters, "vlog.logfile", "Log to given file - can "+
		"be provided multiple times to log to multiple files")
}
