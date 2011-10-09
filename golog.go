package golog

import (
	"os"
)

var Global LevelLogger = &levelLoggerImpl{
	&loggerImpl{&defaultLogOuters, flag_minloglevel},
}

func exitNow() {
	os.Exit(1)
}

func Info(vals ...interface{}) {
	Global.Log(INFO, vals...)
}

func Infof(f string, args ...interface{}) {
	Global.Logf(INFO, f, args...)
}

func Infoc(closure func() string) {
	Global.Logc(INFO, closure)
}

func Warning(vals ...interface{}) {
	Global.Log(WARNING, vals...)
}

func Warningf(f string, args ...interface{}) {
	Global.Logf(WARNING, f, args...)
}

func Warningc(closure func() string) {
	Global.Logc(WARNING, closure)
}

func Error(vals ...interface{}) {
	Global.Log(ERROR, vals...)
}

func Errorf(f string, args ...interface{}) {
	Global.Logf(ERROR, f, args...)
}

func Errorc(closure func() string) {
	Global.Logc(ERROR, closure)
}

func Fatal(vals ...interface{}) {
	Global.Log(FATAL, vals...)
	exitNow()
}

func Fatalf(f string, args ...interface{}) {
	Global.Logf(FATAL, f, args...)
	exitNow()
}

func Fatalc(closure func() string) {
	Global.Logc(FATAL, closure)
	exitNow()
}

func StartTestLogging(t TestController) {
	defaultLogOuters.AddLogOuter("testing", NewTestLogOuter(t))
	// TODO Replace exitNow
}

func StopTestLogging() {
	defaultLogOuters.RemoveLogOuter("testing")
	// TODO Replace exitNow
}
