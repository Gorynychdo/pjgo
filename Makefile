.PHONY: build
build:
	go build -x -v ./cmd/pjgo2

.PHONY: v2
demo:
	go build -x -v ./cmd/pjgo

.PHONY: pjsua
pjsua:
	go install ./internal/pjsua2

.DEFAULT_GOAL := build