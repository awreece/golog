About this package
==================

This package is designed to support better logging for Go. Specifically, this
project aims to support different levels of logging and the ability to
customize log output via custom implementations of the interfaces provided in
the package. In addition, all logged messages are wrapped in closures and are
only evaluated and rendered if they will be outputed.

You can install this package via:

	goinstall github.com/awreece/golog

For additional documenation, see `doc.go`
