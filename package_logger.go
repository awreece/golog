package golog

type PackageLogger struct {
	LocationLogger
	MultiLogOuter
	failFunc func()
}

func newPackageLoggerCommon(outer MultiLogOuter, minloglevel_flag *int,
failFunc func(), locFunc func(skip int) *LogLocation) *PackageLogger {
	ret := &PackageLogger{failFunc: failFunc, MultiLogOuter: outer}

	ret.LocationLogger = NewLocationLogger(
		&loggerImpl{outer, minloglevel_flag, func() { ret.failFunc() }},
		FullLocation)

	return ret
}

func NewDefaultPackageLogger() *PackageLogger {
	return newPackageLoggerCommon(
		NewDefaultMultiLogOuter(),
		flag_minloglevel,
		ExitError,
		FullLocation)
}

func NewPackageLogger(outer MultiLogOuter, minloglevel int,
failFunc func(), locFunc func(skip int) *LogLocation) *PackageLogger {
	return newPackageLoggerCommon(outer, &minloglevel, failFunc, locFunc)
}

func (l *PackageLogger) StartTestLogging(t TestController) {
	l.MultiLogOuter.AddLogOuter("testing", NewTestLogOuter(t))
	// TODO(awreece) Safe old failFunc so we can restore it properly.
	l.failFunc = func() { t.FailNow() }
}

func (l *PackageLogger) StopTestLogging() {
	l.MultiLogOuter.RemoveLogOuter("testing")
	l.failFunc = ExitError
}
func (l *PackageLogger) Info(msg ...interface{}) {
	l.LogDepth(INFO, printClosure(msg...), 1)
}
func (l *PackageLogger) Infof(fmt string, vals ...interface{}) {
	l.LogDepth(INFO, printfClosure(fmt, vals...), 1)
}
func (l *PackageLogger) Infoc(closure func() string) {
	l.LogDepth(INFO, closure, 1)
}
func (l *PackageLogger) Warning(msg ...interface{}) {
	l.LogDepth(WARNING, printClosure(msg...), 1)
}
func (l *PackageLogger) Warningf(fmt string, vals ...interface{}) {
	l.LogDepth(WARNING, printfClosure(fmt, vals...), 1)
}
func (l *PackageLogger) Warningc(closure func() string) {
	l.LogDepth(WARNING, closure, 1)
}
func (l *PackageLogger) Error(msg ...interface{}) {
	l.LogDepth(ERROR, printClosure(msg...), 1)
}
func (l *PackageLogger) Errorf(fmt string, vals ...interface{}) {
	l.LogDepth(ERROR, printfClosure(fmt, vals...), 1)
}
func (l *PackageLogger) Errorc(closure func() string) {
	l.LogDepth(ERROR, closure, 1)
}
func (l *PackageLogger) Fatal(msg ...interface{}) {
	l.LogDepth(FATAL, printClosure(msg...), 1)
	l.FailNow()
}
func (l *PackageLogger) Fatalf(fmt string, vals ...interface{}) {
	l.LogDepth(FATAL, printfClosure(fmt, vals...), 1)
	l.FailNow()
}
func (l *PackageLogger) Fatalc(closure func() string) {
	l.LogDepth(FATAL, closure, 1)
	l.FailNow()
}
