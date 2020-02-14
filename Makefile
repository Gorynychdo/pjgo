.PHONY: build
build:
		go build -x -v ./cmd/pjgo

.PHONY: v2
v2:
		go build -x -v ./cmd/pjgo2

.DEFAULT_GOAL := build