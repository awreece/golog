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
	ret := &packageLoggerImpl{ failFunc: exitNow }

	ret.LevelLogger = NewLevelLogger(
		&loggerImpl{
			NewDefaultMultiLogOuter(),
			flag_minloglevel,
			func() { ret.failFunc() },
		}, FullLocation)

	return ret
}

func (l *packageLoggerImpl) StartTestLogging(t TestController) {
	defaultLogOuters.AddLogOuter("testing", NewTestLogOuter(t))
	l.failFunc = func () { t.FailNow() }
}

func (l *packageLoggerImpl) StopTestLogging() {
	defaultLogOuters.RemoveLogOuter("testing")
	l.failFunc = exitNow
}
