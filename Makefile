include $(GOROOT)/src/Make.inc

TARG=golog
GOFILES=\
	golog.go\
	level_logger.go\
	logger.go\
	log_message.pb.go\
	log_outer.go\
	multi_log_outer.go\
	vmodules.go\

include $(GOROOT)/src/Make.pkg
include $(GOROOT)/src/pkg/goprotobuf.googlecode.com/hg/Make.protobuf
