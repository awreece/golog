include $(GOROOT)/src/Make.inc

TARG=golog
GOFILES=\
	doc.go\
	golog.go\
	location_logger.go\
	logger.go\
	log_message.go\
	log_outer.go\
	multi_log_outer.go\
	package_logger.go\

# We trick godoc into not exporting our mock object by naming it
# mock_object_test.go. 
# TODO(awreece) This feels somewhat hacky.
MOCKFILES=\
	  mock_log_outer_test.go\

override GOTESTFILES+=$(MOCKFILES)

include $(GOROOT)/src/Make.pkg

CLEANFILES+=$(MOCKFILES)

mock_%_test.go: %.go
	mockgen --source=$< --destination=$@ --package=golog
