About this package
==================

This package is designed to support better logging for Go. Specifically, this
project aims to support different levels of logging and the ability to
customize log output via custom implementations of the interfaces provided in
the package. In addition, all logged messages are wrapped in closures and are
only evaluated and rendered if they will be outputed.

Using this package
==================

Introductory usage
------------------

The easiest way to start using this package is to use the `Global` 
`PackageLogger` and the exported global namespace wrapper functions. For
example:

	package mypackage

	import "github.com/awreece/golog"

	func Foo() {
		golog.Info("Hello, world")
		golog.Warningf("Error %d", 4)
		golog.Errorc(func() { return slowMakePrettyString() })
		golog.Fatal("Error opening file:", err)
	}

The `Global` `PackageLogger` output to default files set by flags. For example,
to log to `stderr` and to `temp.log`, invoke the binary with the additional
flags `--golog.logfile=/dev/stderr --golog.logfile=temp.log`.

This package also makes it easy to log to a testing harness in addition to
files. To do this, invoke `StartTestLogging(t)` at the start of every test
and `StopTestLogging()` at the end. For example:
	
	package mypackage
	
	import (
		"github.com/awreece/golog"
		"testing"
	)

	func TestFoo(t *testing.T) {
		golog.StartTestLogging(t); defer golog.StopTestLogging()

		// Test the Foo() function.
		Foo()
	}

While in test logging mode, calls to `Fatal()` (and `DefaultLogger.FailNow()`)
will call `testing.(*T).FailNow()` rather than
exiting the program abruptly.

Another common way to use this pacakge is to create a local `PackageLogger`.
This can either be declared on the package level or passed in by value.

Advanced usage
--------------
This package is highly modular and configurable; different components can be
plugged in to modify the behavior. For example, to speed up logging an advanced
user could try creating a `LevelLogger` using the `NoLocation` function, or
even create a custom location function.

Advanced users can further take advantage of the modularity of the package to 
implement and control individual parts. For example, logging in XML format 
should be done by writing a proper `LogOuter`.

Understanding this package
==========================

There are 3 important objects in this package.
	LogOuter: Outputs a LogMessage to (file, testing.T, network, xml, etc)
	Logger: Decides whether on not to generate output
	LevelLogger: Easier interface for Logger.

In practice, the user is encouraged to use the LevelLogger as an entrypoint into
the package. The provided Global LevelLogger is set up to have easy defaults
and to be easily configurable with flags and the AddLogFile and the
{Start,Stop}TestLogging functions. As an alternative, the user can create
package specific LevelLogger with their own presets or the default (flag based)
presets.

NOTE: The package is not quite stable. Most exported methods and types will
remain exported, but may change.
