include $(GOROOT)/src/Make.inc

TARG=golog
GOFILES=\
	doc.go\
	golog.go\
	location_logger.go\
	logger.go\
	log_outer.go\
	mock_log_outer.go\
	multi_log_outer.go\
	package_logger.go\

include $(GOROOT)/src/Make.pkg
