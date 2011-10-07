package golog

import (
	"bytes"
	"fmt"
	"math"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	INFO = iota
	WARNING
	ERROR
	FATAL
)

// TODO: Pass in time as an arg so we can check it earlier?
// TODO: Accept func, package as args or rename file to location to be more 
// reusable?
func makePrefix(buf *bytes.Buffer, msg string, file string, line int) {
	// TODO Cite logging package formatHeader
	t := time.NanosecondsToLocalTime(time.Nanoseconds())
	buf.WriteString(msg)
	buf.WriteString(t.Format(" 15:04:05.000000 "))

	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))
	buf.WriteString("] ")
}

var levelStrings []string = []string{"I", "W", "E", "F"}

func internalLogc(level int, c func() string, skip int) {
	closure := func() string {
		// skip + 3 because we have to skip over the call to this
		// closure, the call to output, and the call to
		// internalLogc.
		_, file, line, ok := runtime.Caller(skip + 3)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = path.Base(file)
		}

		var buf bytes.Buffer
		makePrefix(&buf, levelStrings[level], file, line)
		buf.WriteString(c())
		return buf.String()
	}

	// TODO Dirty hack. We pass math.MinInt32 so they 
	// cannot disable output by setting vlog incorrectly.
	output(level, int(math.MinInt32), "package", closure)
}

func internalVLogc(vloglevel int, c func() string, skip int) {
	var pack string
	// skip + 3 because we have to skip over the call to this
	// closure, the call to output, and the call to
	// internalLogc.
	pc, file, line, ok := runtime.Caller(skip + 3)
	if !ok {
		pack = "???"
		file = "???"
		line = 0

	} else {
		// TODO We really want offset from root.
		file = path.Base(file)
		// TODO Make this compiler agnostic.
		// TODO Use more standard syntax for go name.
		// In 8g, a function name is "package . [type .] name", so the
		// string piece before the first dot is our package name.
		pack = strings.SplitN(runtime.FuncForPC(pc).Name(), ".", 3)[0]
	}

	closure := func() string {
		var buf bytes.Buffer

		// VLOG messages get logged at INFO level
		stem := fmt.Sprintf("V(%d)", vloglevel)
		// 3 + skip because we have to skip over the call to this
		// closure, the call to output, and the call to
		// internalLogc.
		makePrefix(&buf, stem, file, line)
		buf.WriteString(c())
		return buf.String()
	}

	outputv(vloglevel, pack, closure)
}

func outputv(vloglevel int, pack string, closure func() string) {
	var closureString string
	var closureEvaled bool
	// TODO Grab mutex.
	for _, outer := range LogOuters.outers {
		if vloglevel <= outer.vmoduleLevel(pack) {
			if !closureEvaled {
				closureString, closureEvaled = closure(), true
			}
			outer.Println(closureString)
		}
	}
}

func output(level int, vloglevel int, pack string, closure func() string) {
	var closureString string
	var closureEvaled bool
	// TODO Grab mutex.
	for _, outer := range LogOuters.outers {
		if level >= *outer.minloglevel{
			if !closureEvaled {
				closureString, closureEvaled = closure(), true
			}
			outer.Println(closureString)
		}

		if level == FATAL {
			outer.FailNow()
		}
	}
}
