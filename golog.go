package golog

var Global FailLogger = &loggerImpl{&defaultLogOuters, flag_minloglevel}
