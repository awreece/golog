/*
This package is designed to support better logging for Go. Specifically, this
project aims to support different levels of logging and the ability to
customize log output via custom implementations of the interfaces provided in
the package. In addition, all logged messages are wrapped in closures and are
only evaluated and rendered if they will be outputed.

The easiest way to start using this package is to use the Global 
PackageLogger and the exported global namespace wrapper functions. For
example:

	package mypackage

	import "github.com/awreece/golog"

	func Foo() {
		golog.Info("Hello, world")
		golog.Warningf("Error %d", 4)
		golog.Errorc(func() { return verySlowStringFunction() })
		golog.Fatal("Error opening file:", err)
	}

The Global PackageLogger output to default files set by flags. For example,
to log to stderr and to temp.log, invoke the binary with the additional
flags --golog.logfile=/dev/stderr --golog.logfile=temp.log.

This package also makes it easy to log to a testing harness in addition to
files. To do this, invoke StartTestLogging(t) at the start of every test
and StopTestLogging() at the end. For example:
	
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

While in test logging mode, calls to Fatal() (and DefaultLogger.FailNow())
will call testing.(*T).FailNow() rather than
exiting the program abruptly.

Another common way to use this pacakge is to create a local PackageLogger.
This can either be declared on the package level or passed in by value.

Advanced usage
This package is highly modular and configurable; different components can be
plugged in to modify the behavior. For example, to speed up logging an advanced
user could try creating a LocationLogger using the NoLocation function, or
even create a custom location function.

Advanced users can further take advantage of the modularity of the package to 
implement and control individual parts. For example, logging in XML format 
should be done by writing a proper LogOuter.

This package was designed to be highly modular, with different interfaces for
each logical component. The important types are:

-	A LogMessage is a logged message with associated metadata.

-	A LogOuter controls outputing a LogMessage.

-	A MultiLogOuter multiplexes an outputted message to a set of keyed
LogOuters. The associated MultiLogOuterFlag automatically add 
logfiles to the associated set of LogOuters.

-	A Logger decides whether or not to log a message, and if so renders 
the message and outputs it.

-	A LocationLogger is a wrapper for a Logger that generates a closure
to return a LogMessage with the associate metadata and is the first 
easily usable entrypoint into this package.

-	A PackageLogger has a set of functions designed be quickly useful
and is the expected entry point into this package.
*/
package golog
