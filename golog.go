package golog

const (
	INFO = iota
	WARNING
	ERROR
	FATAL
)

var Global PackageLogger = NewDefaultPackageLogger()

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
	Global.FailNow()
}

func Fatalf(f string, args ...interface{}) {
	Global.Logf(FATAL, f, args...)
	Global.FailNow()
}

func Fatalc(closure func() string) {
	Global.Logc(FATAL, closure)
	Global.FailNow()
}

func Log(level int, vals ...interface{}) {
	Global.Log(level, vals...)
}

func Logf(level int, f string, args ...interface{}) {
	Global.Logf(level, f, args...)
}

func Logc(level int, closure func() string) {
	Global.Logc(level, closure)
}

func StartTestLogging(t TestController) {
	Global.StartTestLogging(t)
}

func StopTestLogging() {
	Global.StopTestLogging()
}
