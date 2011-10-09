package golog

type PackageLogger struct {
	pack string
	LevelLogger
}

func NewDefaultPackageLogger(pack string) *PackageLogger {
	return &PackageLogger{pack, &levelLoggerImpl{
		DefaultLogger,
		func(int) *LogLocation { return &LogLocation{Package: pack} },
	}}
}
