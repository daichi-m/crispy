.DEFAULT_GOAL := all
SOURCES := $(shell find . -prune -o -name "*.$(go)" -not -name '*_test.$(go)' -print)

GO111MODULE ?= on
GO ?= go


.PHONY: setup
setup:
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) get -u

.PHONY: fmt
fmt:
	goimports -w .

.PHONY: tests
tests: 
	$(GO) test -race -covermode atomic -coverprofile coverage.txt .

.PHONY: build
build: setup fmt
	$(GO) build .

.PHONY: all
all: setup fmt tests build