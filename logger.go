package golog

// Logger.Log uses the level to determine whether or not to output the
// arguments. Logger.Log will output the provided arguments exactly, without
// additional formatting such as adding a prefix etc. In addition, Logger.Log
// must be thread safe.
type Logger interface {
	// If the message is to be logged, evaluates the closure and outputs
	// the result.
	Log(level int, closure func() string)
}

// A FailLogger is a Logger with the addition FailNow() function, which flushes
// the Logger and performs some action. The action performed by FailNow() is 
// deliberately unspecified, but could include os.Exit(1) or 
// testing.(*T).FailNow().
type FailLogger interface {
	Logger
	FailNow()
}
