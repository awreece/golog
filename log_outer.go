// This file definse the logOuter interface and several types of logOuter.
//
// emptyOuter = logOuter where both Out and Outf are noops
// lineOuter = logOuter where a newline is inserted after every call to
//			   Out and Outf
// fatalLineOuter = logOuter that logs message with inserted newline then
//					exits with call to os.EXIT(1)

package golog

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"sync"
)

type LogOuter interface {
	// Output a LogMessage (to a file, to stderr, to a tester, etc). Output
	// must be safe to call from multiple threads.
	Output(*LogMessage)
}

type writerLogOuter struct {
	lock sync.Mutex
	io.Writer
}

func (f *writerLogOuter) Output(m *LogMessage) {
	f.lock.Lock()
	defer f.lock.Unlock()

	// TODO(awreece) Handle short write?
	// Make sure to insert a newline.
	f.Write([]byte(formatLogMessage(m, true)))
}

// Returns a LogOuter wrapping the io.Writer.
func NewWriterLogOuter(f io.Writer) LogOuter {
	return &writerLogOuter{io.Writer: f}
}

// Returns a LogOuter wrapping the file, or an error if the file cannot be
// opened.
func NewFileLogOuter(filename string) (LogOuter, error) {
	// TODO(awreece) Permissions?
	if file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		return nil, err
	} else {
		return NewWriterLogOuter(file), nil
	}

	panic("Code never reaches here, this mollifies the compiler.")
}

// We want to allow an abitrary testing framework.
type TestController interface {
	// We will assume that testers insert newlines in manner similar to 
	// the FEATURE of testing.T where it inserts extra newlines. >.<
	Log(...interface{})
	FailNow()
}

type testLogOuter struct {
	TestController
}

func (t *testLogOuter) Output(m *LogMessage) {
	// Don't insert an additional log message since the tester inserts them
	// for us.
	t.Log(formatLogMessage(m, false))
}

// Return a LogOuter wrapping the TestControlller.
func NewTestLogOuter(t TestController) LogOuter {
	return &testLogOuter{t}
}

type udpLogOuter struct {
	conn  net.PacketConn
	raddr net.Addr
}

// Returns a LogOuter that forwards LogMessages in json format to UDP network
// address. TODO(awreece): Use protobuf?
func NewUDPLogOuter(raddr string) (LogOuter, error) {
	var addr *net.UDPAddr
	var err error
	var conn net.PacketConn

	if addr, err = net.ResolveUDPAddr("udp", raddr); err != nil {
		return nil, err
	}

	if conn, err = net.DialUDP("udp", nil, addr); err != nil {
		return nil, err
	}

	return &udpLogOuter{conn, addr}, nil
}

func (o *udpLogOuter) Output(m *LogMessage) {
	// TODO(awreece): Add Hostname if not present?
	if bytes, err := json.Marshal(m); err == nil {
		// TODO(awreece) Handle error?
		o.conn.WriteTo(bytes, o.raddr)
	}
}
