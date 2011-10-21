package golog


import (
	"bytes"
	"fmt"
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
