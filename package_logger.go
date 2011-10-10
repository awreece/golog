package golog

type PackageLogger interface {
	LevelLogger
	StartTestLogging(TestController)
	StopTestLogging()
	AddLogOuter(key string, outer LogOuter)
	RemoveLogOuter(key string)
}

type packageLoggerImpl struct {
	LevelLogger
	MultiLogOuter
	failFunc func()
}

func NewDefaultPackageLogger() PackageLogger {
	return NewPackageLogger(
		NewDefaultMultiLogOuter(),
		flag_minloglevel,
		exitNow,
		FullLocation)
}

func NewPackageLogger(outer MultiLogOuter, minloglevel_flag *int,
failFunc func(), locFunc func(skip int) *LogLocation) PackageLogger {
	ret := &packageLoggerImpl{failFunc: failFunc, MultiLogOuter: outer}

	ret.LevelLogger = NewLevelLogger(
		&loggerImpl{outer, minloglevel_flag, func() { ret.failFunc() }},
		FullLocation)

	return ret
}

func (l *packageLoggerImpl) StartTestLogging(t TestController) {
	l.MultiLogOuter.AddLogOuter("testing", NewTestLogOuter(t))
	l.failFunc = func() { t.FailNow() }
}

func (l *packageLoggerImpl) StopTestLogging() {
	l.MultiLogOuter.RemoveLogOuter("testing")
	l.failFunc = exitNow
}
