package golog

// Logger.Log{,f,c} methods use the level to determine whether or not to 
// output the arguments. Logger.Log{,f,c} will output the provided arguments
// exactly, without additional formatting such as adding a prefix etc.
// Logger.Log{,f,c} must be thread safe.
type Logger interface {
	// If the message is to be logged, outputs the values as if via a call
	// to fmt.Sprint(vals...).
	Log(level int, vals ...interface{})
	// If the message is to be logged, outputs the values as if via a call
	// to fmt.Sprintf(fmt, args...).
	Logf(level int, fmt string, args ...interface{})
	// If the message is to be logged, evaluates the closure and outputs
	// the result.
	Logc(level int, closure func() string)
}
