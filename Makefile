# 
# Makefile for copyright-notice
# 
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
GOPATH?=`$(GOCMD) env GOPATH`

BINARY=copyright-notice
BINARY_DARWIN=$(BINARY)_darwin
BINARY_LINUX=$(BINARY)_linux
BINARY_PI=$(BINARY)_pi
BINARY_WINDOWS=$(BINARY).exe

TESTS=./...
COVERAGE_FILE=coverage.out

BUILD=build/
GO_VERSION=1.15

BUILD_DATE=`date`
BUILD_COMMIT=`git rev-parse HEAD`

.PHONY: all test test-ci build build-mac build-linux build-windows build-all coverage clean test-docker build-docker nightly toc

all: test build

build:
		$(GOBUILD) -o $(BINARY) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

build-mac:
		GOOS="darwin" GOARCH="amd64" $(GOBUILD) -o $(BINARY_DARWIN) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

build-linux:
		GOOS="linux" GOARCH="amd64" $(GOBUILD) -o $(BINARY_LINUX) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

build-pi:
		GOOS="linux" GOARCH="arm" GOARM="6" $(GOBUILD) -o $(BINARY_PI) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

build-windows:
		GOOS="windows" GOARCH="amd64" $(GOBUILD) -o $(BINARY_WINDOWS) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'"

build-all: build-mac build-linux build-windows

test:
		$(GOTEST) -v $(TESTS)

test-ci:
		$(GOTEST) -v -short ./...

coverage:
		$(GOTEST) -coverprofile=$(COVERAGE_FILE) $(TESTS)
		$(GOTOOL) cover -html=$(COVERAGE_FILE)

clean:
		$(GOCLEAN)
		rm -rf $(BINARY) $(BINARY_DARWIN) $(BINARY_LINUX) $(BINARY_PI) $(BINARY_WINDOWS) $(COVERAGE_FILE) ${BUILD} dist/*

test-docker:
		docker run --rm -v "${GOPATH}":/go -w /go/src/creativeprojects/copyright-notice golang:${GO_VERSION} $(GOTEST) -v $(TESTS)

build-docker: clean
		CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GOBUILD) -v -ldflags "-X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}' -X 'main.builtBy=make'" -o ${BUILD}$(BINARY) .
		cd ${BUILD}; docker build --pull --tag creativeprojects/copyright-notice .

nightly:
	go install github.com/goreleaser/goreleaser
	go mod tidy
	goreleaser --snapshot --skip-publish --rm-dist

toc:
	go install github.com/ekalinin/github-markdown-toc.go
	go mod tidy
	cat README.md | github-markdown-toc.go --hide-footer
