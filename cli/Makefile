.DEFAULT_GOAL := build

export PROJECT := "navy-seal"
export PROJECT_BIN_DIR := "./dist/bin"
export UNAME := $(shell uname)
ifeq ($(findstring MINGW64_NT,$(UNAME)),MINGW64_NT)
    # Found
    DETECTED_OS=Windows64
endif


showos:
	@echo " the OS is [${DETECTED_OS}]"
	@echo " the native OS is [$(OS)]"
lint:
	golangci-lint run

go-fetch:
	go mod download
	go mod tidy

up:
	go get -u ./...
	go mod tidy

clean:
	/bin/rm -rfv "dist/" "${PROJECT}"

go-dlv: go-prepare
	dlv debug \
		--headless --listen=:2345 \
		--api-version=2 --log \
		--allow-non-terminal-interactive \
		${PACKAGE} -- --debug

go-debug: clean
	go run *.go --debug

go-prepare: go-fetch clean
	go generate -x ./...


ifeq (${DETECTED_OS},Windows64)
build: go-prepare go-fetch
	CGO_ENABLED=0 \
	$(shell mkdir -p ${PROJECT_BIN_DIR}) \
	go build \
		-ldflags '-s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-buildvcs=false \
		-trimpath \
		-o ${PROJECT_BIN_DIR}/${PROJECT}
endif


ifneq (${DETECTED_OS},Windows64)
build: go-prepare go-fetch
	CGO_ENABLED=0 \
	$(shell mkdir -p ${PROJECT_BIN_DIR}) \
	go build \
		-ldflags '-d -s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-buildvcs=false \
		-trimpath \
		-o ${PROJECT_BIN_DIR}/${PROJECT}
endif


