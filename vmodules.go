package vlog

import (
	"bytes"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

// Sets default verbose logging level.
var flag_vloglevel = flag.Int("vlog.v", 0, "Show all Vlog(l) messages for l "+
	"less or equal the value of this flag. Note that this is the opposite "+
	"behavior of minloglevel")

type vmoduleLevelsMap map[string]int

func (m vmoduleLevelsMap) String() string {
	buf := new(bytes.Buffer)
	buf.WriteByte('"')
	for pattern, level := range m {
		// TODO: Be better about inserting the comma.
		buf.WriteString(fmt.Sprint(pattern, "=", level, ","))
	}
	buf.WriteByte('"')
	return buf.String()
}

func (m vmoduleLevelsMap) Set(v string) bool {
	for _, setting := range strings.Split(v, ",") {
		ss := strings.SplitN(setting, "=", 2)
		if len(ss) < 2 {
			return false
		}

		pattern := ss[0]

		if level, ok := strconv.Atoi(ss[1]); ok == nil {
			m[pattern] = level
		} else {
			// TODO: On error, flush map or leave half modified?
			return false
		}
	}

	return true
}

// Per-module verbose level. The argument has to contain a comma-separated list
// of <module name>=<log level>. <module name> is a glob pattern (e.g., gfs* for
// all modules whose name starts with "gfs"), matched against the filename base 
// (that is, name ignoring .cc/.h./-inl.h). <log level> overrides any value given 
// by --v. 
//
// WARNING: Globbing support not provided yet!
// TODO: Remove warning when globbing support added.
var flag_vmodule vmoduleLevelsMap = make(map[string]int)

func init() {
	flag.Var(flag_vmodule, "vlog.vmodule", "Per-module verbose level. Example "+
		"usage: --vlog.vmodule=mapreduce=2,file=1,gfs*=3")
}

// Returns the vlog level for the given module.
func (m vmoduleLevelsMap) vmoduleLevel(module string) int {
	// TODO iterate in a sane order
	for pattern, level := range m {
		if match, err := filepath.Match(pattern, module); err == nil &&
		match {
			return level
		}
	}

	return *flag_vloglevel
}

// TODO: Remove this for final libary.
func SetVerbose(level int) {
	*flag_vloglevel = level
}
