package golog

import (
	"flag"
	"fmt"
	"os"
)

type MultiLogOuter interface {
	LogOuter
	AddLogOuter(key string, outer LogOuter)
	RemoveLogOuter(key string)
}

type MultiLogOuterFlag interface {
	MultiLogOuter
	flag.Value
}

type multiLogOuterImpl struct {
	// TODO Add mutex.
	outers map[string]LogOuter
}

func (l *multiLogOuterImpl) String() string {
	// TODO better string
	return fmt.Sprint("\"", l.outers, "\"")
}

func (l *multiLogOuterImpl) Set(name string) bool {
	if outer, err := NewFileLogOuter(name); err != nil {
		os.Stderr.WriteString(
			fmt.Sprint("Error opening file for logging", name,
				": ", err))
		return false
	} else {
		l.AddLogOuter(name, outer)
		return true
	}

	panic("Code never reaches here, this mollifies the compiler.")
}

func (l *multiLogOuterImpl) AddLogOuter(key string, outer LogOuter) {
	// TODO Grab mutex.
	l.outers[key] = outer
}

func (l *multiLogOuterImpl) RemoveLogOuter(key string) {
	// TODO Grab mutex.
	// TODO Be Go1 compatible. :)
	l.outers[key] = nil, false
}

func (l *multiLogOuterImpl) Output(m *LogMessage) {
	// TODO Grab mutex.
	for _, outer := range l.outers {
		outer.Output(m)
	}
}

var defaultLogOuters MultiLogOuterFlag = NewMultiLogOuter()

func NewDefaultMultiLogOuter() MultiLogOuterFlag {
	return &multiLogOuterImpl{
		outers: map[string]LogOuter{"default": defaultLogOuters},
	}
}

func NewMultiLogOuter() MultiLogOuterFlag {
	return &multiLogOuterImpl{make(map[string]LogOuter)}
}

func init() {
	flag.Var(defaultLogOuters, "golog.logfile",
		"Log to given file - can be provided multiple times to log "+
			"to multiple files")
}
