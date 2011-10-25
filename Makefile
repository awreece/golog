include $(GOROOT)/src/Make.inc

TARG=golog
GOFILES=\
	doc.go\
	golog.go\
	location_logger.go\
	logger.go\
	log_message.go\
	log_outer.go\
	mock_log_outer.go\
	multi_log_outer.go\
	package_logger.go\

include $(GOROOT)/src/Make.pkg

mock_log_outer.go:
	mockgen --source=log_outer.go --destination=mock_log_outer.go --package=golog

CLEANFILES+=\
	    mock_log_outer.go\
