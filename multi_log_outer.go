package golog

import (
	"flag"
	"fmt"
	"os"
)

type multiLogOuter struct {
	// TODO Insert mutex here.
	outers []*loggerImpl
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
	// TODO Grab mutex.
	l.outers = append(l.outers, &loggerImpl{
		LogOuter:         &fileLogOuter{file},
		minloglevel:      flag_minloglevel,
		vmoduleLevelsMap: flag_vmodule,
	})
}

func (l *multiLogOuter) AddDefaultLogTester(t TestController) {
	// TODO Grab Mutex
	l.outers = append(l.outers, &loggerImpl{
		LogOuter:         &testLogOuter{t},
		minloglevel:      flag_minloglevel,
		vmoduleLevelsMap: flag_vmodule,
	})
}

func (l *multiLogOuter) Println(s string) {
}

func (l *multiLogOuter) FailNow() {
}

var defaultLogOuters multiLogOuter

func init() {
	flag.Var(&defaultLogOuters, "vlog.logfile", "Log to given file - can "+
		"be provided multiple times to log to multiple files")
}
