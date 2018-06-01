SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
OS=$(shell uname -s)

setup:
	go get -u golang.org/x/tools/cmd/cover
	go get -u gopkg.in/alecthomas/gometalinter.v2
ifeq ($(OS), Darwin)
	brew install dep
else
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure -vendor-only
.PHONY: setup

# Run all the tests
test:
	go test $(TEST_OPTIONS) -failfast -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m
.PHONY: test

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt


# Run all the lintersmake 
lint:
	gometalinter.v2 --vendor --deadline 2m ./...
.PHONY: lint


# Run all the tests and code checks
ci: build test lint
.PHONY: ci

# Build a beta version
build:
	go generate ./...
	go build
.PHONY: build

.DEFAULT_GOAL := build