SWIG_DIR := ~/src/pjproject-2.9/pjsip-apps/src/swig/
CGO_CXXFLAGS := -I/usr/local/include

.PHONY: pjsua-install
pjsua-install:
	rm -rf internal/pkg/pjsua2\
	&& mkdir internal/pkg/pjsua2\
	&& cd internal/pkg/pjsua2\
	&& cp $(SWIG_DIR)/pjsua2.i .\
	&& cp $(SWIG_DIR)/symbols.i .\
	&& swig -go -cgo -intgosize 64 $(CGO_CXXFLAGS) -c++ pjsua2.i

.PHONY: build
build:
	go build -x -v -i ./cmd/pjgo

.DEFAULT_GOAL := build