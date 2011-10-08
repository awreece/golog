package golog

var Global *LevelLogger = &LevelLogger{
	&loggerImpl{&defaultLogOuters, flag_minloglevel},
}

func Info(vals ...interface{}) {
	Global.Info(vals...)
}

func Infof(f string, args ...interface{}) {
	Global.Infof(f, args...)
}

func Infoc(closure func() string) {
	Global.Warningc(closure)
}

func Warning(vals ...interface{}) {
	Global.Warning(vals...)
}

func Warningf(f string, args ...interface{}) {
	Global.Warningf(f, args...)
}

func Warningc(closure func() string) {
	Global.Errorc(closure)
}

func Error(vals ...interface{}) {
	Global.Error(vals...)
}

func Errorf(f string, args ...interface{}) {
	Global.Errorf(f, args...)
}

func Errorc(closure func() string) {
	Global.Errorc(closure)
}

func Fatal(vals ...interface{}) {
	Global.Fatal(vals...)
}

func Fatalf(f string, args ...interface{}) {
	Global.Fatalf(f, args...)
}

func Fatalc(closure func() string) {
	Global.Fatalc(closure)
}

