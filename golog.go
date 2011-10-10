package golog

const (
	INFO = iota
	WARNING
	ERROR
	FATAL
)

var Global *PackageLogger = NewDefaultPackageLogger()


func Info(msg ...interface{}) {
	Global.LogDepth(INFO, printClosure(msg...), 1)
}

func Infof(fmt string, vals ...interface{}) {
	Global.LogDepth(INFO, printfClosure(fmt, vals...), 1)
}

func Infoc(closure func() string) {
	Global.LogDepth(INFO, closure, 1)
}

func Warning(msg ...interface{}) {
	Global.LogDepth(WARNING, printClosure(msg...), 1)
}

func Warningf(fmt string, vals ...interface{}) {
	Global.LogDepth(WARNING, printfClosure(fmt, vals...), 1)
}

func Warningc(closure func() string) {
	Global.LogDepth(WARNING, closure, 1)
}

func Error(msg ...interface{}) {
	Global.LogDepth(ERROR, printClosure(msg...), 1)
}

func Errorf(fmt string, vals ...interface{}) {
	Global.LogDepth(ERROR, printfClosure(fmt, vals...), 1)
}

func Errorc(closure func() string) {
	Global.LogDepth(ERROR, closure, 1)
}

func Fatal(msg ...interface{}) {
	Global.LogDepth(FATAL, printClosure(msg...), 1)
	Global.FailNow()
}

func Fatalf(fmt string, vals ...interface{}) {
	Global.LogDepth(FATAL, printfClosure(fmt, vals...), 1)
}

func Fatalc(closure func() string) {
	Global.LogDepth(FATAL, closure, 1)
}

func StartTestLogging(t TestController) {
	Global.StartTestLogging(t)
}

func StopTestLogging() {
	Global.StopTestLogging()
}

func AddLogOuter(key string, outer LogOuter) {
	Global.AddLogOuter(key, outer)
}

func RemoveLogOuter(key string) {
	Global.RemoveLogOuter(key)
}

func SetMinLogLevel(level int) {
	Global.SetMinLogLevel(level)
}
