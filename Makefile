include $(GOROOT)/src/Make.inc

TARG=golog
GOFILES=\
	outers.go\
	golog.go\
	internal.go\
	vmodules.go\

include $(GOROOT)/src/Make.pkg
