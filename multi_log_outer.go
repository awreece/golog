package golog

import (
	"flag"
	"fmt"
	"os"
)

type multiLogOuter struct {
	outers []LogOuter
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
		l.AddDefaultLogFile(name, file)
		return true
	}

	panic("Code never reaches here, this mollifies the compiler.")
}

// TODO only require filename.
func (l *multiLogOuter) AddDefaultLogFile(filename string, file *os.File) {
	l.outers = append(l.outers, &fileLogOuter{file})
}

func (l *multiLogOuter) AddDefaultLogTester(t TestController) {
	l.outers = append(l.outers, &testLogOuter{t})
}

func (l *multiLogOuter) Println(s string) {
	for _, outer := range l.outers {
		outer.Println(s)
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
