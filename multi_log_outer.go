package golog

import (
	"flag"
	"fmt"
	"os"
)

// A MultiLogOuter is a LogOuter with multiple keyed LogOuters. All functions
// should be safe to call in a multi-threaded environment.
type MultiLogOuter interface {
	LogOuter
	// Add the LogOuter, associating it with the key.
	AddLogOuter(key string, outer LogOuter)
	// Remove the LogOuter associated with the key.
	RemoveLogOuter(key string)
}

// A MultiLogOuter than can also be used as a flag for setting logfiles. 
// For example, it is possible to use a logger other than default via:
// 	var myOuter MultiLogOuterFlag = NewMultiLogOuter()
// 	
// 	func init() {
// 		flag.Var(myOuter, 
// 			"mypack.logfile", 
// 			"Log to file - can be provided multiple times")
// 	}
type MultiLogOuterFlag interface {
	MultiLogOuter
	flag.Value
}

type multiLogOuterImpl struct {
	// TODO(awreece) Add mutex.
	outers map[string]LogOuter
}

func (l *multiLogOuterImpl) String() string {
	// TODO(awreece) better string
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
	// TODO(awreece) Grab mutex.
	l.outers[key] = outer
}

func (l *multiLogOuterImpl) RemoveLogOuter(key string) {
	// TODO(awreece) Grab mutex.
	// TODO(awreece) Be Go1 compatible. :)
	l.outers[key] = nil, false
}

func (l *multiLogOuterImpl) Output(m *LogMessage) {
	// TODO(awreece) Grab mutex.
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
