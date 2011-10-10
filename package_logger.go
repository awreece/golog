package golog

type PackageLogger interface {
	LevelLogger
	StartTestLogging(TestController)
	StopTestLogging()
}

type packageLoggerImpl struct {
	LevelLogger
	failFunc func()
}

func NewDefaultPackageLogger() PackageLogger {
	return NewPackageLogger(
		NewDefaultMultiLogOuter(),
		flag_minloglevel,
		exitNow,
		FullLocation)
}

func NewPackageLogger(outer LogOuter, minloglevel_flag *int,
failFunc func(), locFunc func(skip int) *LogLocation) PackageLogger {
	ret := &packageLoggerImpl{failFunc: failFunc}

	ret.LevelLogger = NewLevelLogger(
		&loggerImpl{outer, minloglevel_flag, func() { ret.failFunc() }},
		FullLocation)

	return ret
}

func (l *packageLoggerImpl) StartTestLogging(t TestController) {
	defaultLogOuters.AddLogOuter("testing", NewTestLogOuter(t))
	l.failFunc = func() { t.FailNow() }
}

func (l *packageLoggerImpl) StopTestLogging() {
	defaultLogOuters.RemoveLogOuter("testing")
	l.failFunc = exitNow
}
