// Better logging for Go.
package golog

var Global *PackageLogger = NewDefaultPackageLogger()

func Info(msg ...interface{}) {
	Global.logger.LogDepth(INFO, printClosure(msg...), 1)
}

func Infof(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(INFO, printfClosure(fmt, vals...), 1)
}

func Infoc(closure func() string) {
	Global.logger.LogDepth(INFO, closure, 1)
}

func Warning(msg ...interface{}) {
	Global.logger.LogDepth(WARNING, printClosure(msg...), 1)
}

func Warningf(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(WARNING, printfClosure(fmt, vals...), 1)
}

func Warningc(closure func() string) {
	Global.logger.LogDepth(WARNING, closure, 1)
}

func Error(msg ...interface{}) {
	Global.logger.LogDepth(ERROR, printClosure(msg...), 1)
}

func Errorf(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(ERROR, printfClosure(fmt, vals...), 1)
}

func Errorc(closure func() string) {
	Global.logger.LogDepth(ERROR, closure, 1)
}

func Fatal(msg ...interface{}) {
	Global.logger.LogDepth(FATAL, printClosure(msg...), 1)
	Global.logger.FailNow()
}

func Fatalf(fmt string, vals ...interface{}) {
	Global.logger.LogDepth(FATAL, printfClosure(fmt, vals...), 1)
	Global.logger.FailNow()
}

func Fatalc(closure func() string) {
	Global.logger.LogDepth(FATAL, closure, 1)
	Global.logger.FailNow()
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
