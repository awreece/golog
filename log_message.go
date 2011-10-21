package golog

import (
	"bytes"
	"fmt"
	"os"
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
// Walks up the stack skip frames and returns the metatdata for that frame.
type MetadataFunc func(skip int) map[string]string

var NoLocation MetadataFunc = func(skip int) map[string]string {
	// TODO: Add timestamp?
	return make(map[string]string)
}

type LocationFlag int

const (
	None LocationFlag = 1 << iota
	Package
	Function
	File
	Line
	Hostname
	Default    = Package | Function | File | Line
	All        = Package | Function | File | Line | Hostname
	requiresPC = Package | Function | File | Line
)


// Returns a function the computes the specified fields of metadata for the log
// message.
// 
// flags is the set of locations to add to the metadata. For example, 
//	MakeMetadataFunc(File | Line | Hostname)
func MakeMetadataFunc(flags LocationFlag) MetadataFunc {
	return func(skip int) map[string]string {
		ret := NoLocation(skip + 1)

		// TODO(awreece) Refactor.
		if flags|requiresPC > 0 {
			// Don't get the pc unless we have to.
			if pc, file, line, ok := runtime.Caller(skip + 1); ok {
				// Don't get FuncForPC unless we have to.
				if flags|Package > 0 || flags|Function > 0 {
					// TODO(awreece) Make sure this is 
					// compiler agnostic.
					funcParts := strings.SplitN(
						runtime.FuncForPC(pc).Name(),
						".", 2)
					if flags|Package > 0 {
						ret["package"] = funcParts[0]
					}
					if flags|Function > 0 {
						ret["function"] = funcParts[1]
					}
				}

				if flags|File > 0 {
					ret["file"] = path.Base(file)
				}
				if flags|Line > 0 {
					ret["line"] = strconv.Itoa(line)
				}

			}
		}
		if flags|Hostname > 0 {
			if host, err := os.Hostname(); err == nil {
				ret["hostname"] = host
			}
		}
		return ret
	}
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
