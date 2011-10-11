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

# We have to manually add mock_log_outer here because mockgen is not in a sane
# state and is unable to produce mock_log_outer.go (for progress on bug, see
# http://code.google.com/p/gomock/issues/detail?id=7). We generated
# mock_log_outer.go with the command:
# 	mockgen --source=log_outer.go --destination=mock_log_outer.go --package=golog

include $(GOROOT)/src/Make.pkg
