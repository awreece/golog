package golog

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogMessage struct {
	Level       int
	Nanoseconds int64
	Message     string
	// A map from the type of metadata to the metadata, if present.
	// By convention, fields in this map will be entirely lowercase and
	// single word.
	Metadata map[string]string
}

// TODO(awreece) comment this
// Skip 0 refers to the function calling this function.
type MakeMetadataFunc func(skip int) map[string]string

// Return a nil LogLocation.
func NoLocation(skip int) map[string]string { return make(map[string]string) }

// Walks up the stack skip frames and returns the metatdata for that frame.
// TODO(awreece) Provide a arg to select which fields to produce?
func FullLocation(skip int) map[string]string {
	ret := make(map[string]string)
	// TODO add timestamp?

	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return ret
	} else {
		// TODO(awreece) Make sure this is compiler agnostic.
		funcParts := strings.SplitN(runtime.FuncForPC(pc).Name(), ".", 2)

		ret["package"] = funcParts[0]
		ret["file"] = path.Base(file)
		ret["function"] = funcParts[1]
		ret["line"] = strconv.Itoa(line)

		return ret
	}

	panic("Flow never reaches here, this mollifies the compiler")
}

// Render the formatted metadata to the buffer. If all present, format is 
// "{time} {pack}.{func}/{file}:{line}". If some fields omitted, intelligently
// delimits the remaining fields.
func renderMetadata(buf *bytes.Buffer, m *LogMessage) {
	if m == nil {
		// TODO Panic here?
		return
	}

	t := time.NanosecondsToLocalTime(m.Nanoseconds)
	buf.WriteString(t.Format(" 15:04:05.000000"))

	packName, packPresent := m.Metadata["package"]
	file, filePresent := m.Metadata["file"]
	funcName, funcPresent := m.Metadata["function"]
	line, linePresent := m.Metadata["line"]

	if packPresent || filePresent || funcPresent || linePresent {
		buf.WriteString(" ")
	}

	// TODO(awreece) This logic is terrifying.
	if packPresent {
		buf.WriteString(packName)
	}
	if funcPresent {
		if packPresent {
			buf.WriteString(".")
		}
		buf.WriteString(funcName)
	}
	if (packPresent || funcPresent) && (filePresent || linePresent) {
		buf.WriteString("/")
	}
	if filePresent {
		buf.WriteString(file)
	}
	if linePresent {
		if filePresent {
			buf.WriteString(":")
		}
		buf.WriteString(line)
	}
}

// Format the message as a string, optionally inserting a newline.
// Format is: "L{level} {time} {pack}.{func}/{file}:{line}] {message}"
func formatLogMessage(m *LogMessage, insertNewline bool) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("L%d", m.Level))
	renderMetadata(&buf, m)
	buf.WriteString("] ")
	buf.WriteString(m.Message)
	if insertNewline {
		buf.WriteString("\n")
	}
	return buf.String()
}
