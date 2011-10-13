package golog

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sync"
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

// REVIEW(korfuri) I'm a bit doubtful of how useful being able to
// remove a log outer is. It seems that it complicates the code a lot
// here by requiring much more synchronisation.

// REVIEW(korfuri) This being said, I'm not a big fan of the
// mutex-based approach here. A goroutine-based solution seems more
// go-ish. Creating a multiLogOuterImpl can start a goroutine that
// will select on several channels : one channel for messages (that
// will be broadcasted to all other LogOuters), one channel to get new
// LogOuters (and eventually one to remove them from the list of
// outers), and one channel to destroy the multiLogOuterImpl. This way
// you get implicit synchronisation instead of explicit.

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
	lock   sync.Mutex
	outers map[string]LogOuter
}

func (l *multiLogOuterImpl) String() string {
	l.lock.Lock()
	defer l.lock.Unlock()

	var buf bytes.Buffer
	buf.WriteString("\"")

	var first bool = true
	for filename, _ := range l.outers {
		if first {
			first = false
		} else {
			buf.WriteString(",")
		}
		buf.WriteString(filename)
	}

	buf.WriteString("\"")
	return buf.String()
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
	l.lock.Lock()
	defer l.lock.Unlock()

	l.outers[key] = outer
}

func (l *multiLogOuterImpl) RemoveLogOuter(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.outers[key] = nil, false
}

func (l *multiLogOuterImpl) Output(m *LogMessage) {
	l.lock.Lock()
	defer l.lock.Unlock()

	for _, outer := range l.outers {
		outer.Output(m)
	}
}

var defaultLogOuters MultiLogOuterFlag = NewMultiLogOuter()

// Create a new MultiLogOuter initialized with a mapping of "default" to the 
// default MultiLogOuter.
func NewDefaultMultiLogOuter() MultiLogOuterFlag {
	return &multiLogOuterImpl{
		outers: map[string]LogOuter{"default": defaultLogOuters},
	}
}

// Create an empty new MutliLogOuter.
func NewMultiLogOuter() MultiLogOuterFlag {
	return &multiLogOuterImpl{outers: make(map[string]LogOuter)}
}

func init() {
	flag.Var(defaultLogOuters, "golog.logfile",
		"Log to given file - can be provided multiple times to log "+
			"to multiple files")
}
