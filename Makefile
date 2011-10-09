include $(GOROOT)/src/Make.inc

TARG=golog
GOFILES=\
	golog.go\
	level_logger.go\
	logger.go\
	log_outer.go\
	multi_log_outer.go\
	package_logger.go\
	vmodules.go\

include $(GOROOT)/src/Make.pkg
