package golog

type PackageLogger struct {
	LevelLogger
}

func NewDefaultPackageLogger(pack string) *PackageLogger {
	return &PackageLogger{&levelLoggerImpl{
		DefaultLogger,
		func(int) *LogLocation { return &LogLocation{Package: pack} },
	}}
}
