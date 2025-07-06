GO ?= go
EXECUTABLE := reviewbot
GOFILES := $(shell find . -type f -name "*.go")

VERSION_PACKAGE := github.com/loveRyujin/ReviewBot/pkg/version

## 定义 VERSION 语义化版本号
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif

## 检查代码仓库是否是 dirty（默认dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
    GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
    -X $(VERSION_PACKAGE).gitVersion=$(VERSION) \
    -X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) \
    -X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) \
    -X $(VERSION_PACKAGE).buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')


## build: build the reviewbot binary
build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	$(GO) mod tidy -v 
	$(GO) build -v -ldflags "$(GO_LDFLAGS)" -o bin/$@ ./cmd/$(EXECUTABLE)

## build_linux_amd64: build the reviewbot binary for linux amd64
build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -ldflags "$(GO_LDFLAGS)" -o release/linux/amd64/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_linux_arm64: build the reviewbot binary for linux arm64
build_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -a -ldflags "$(GO_LDFLAGS)" -o release/linux/arm64/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_linux_arm: build the reviewbot binary for linux arm
build_linux_arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GO) build -a -ldflags "$(GO_LDFLAGS)" -o release/linux/arm/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_mac_intel: build the reviewbot binary for mac intel
build_mac_intel:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -a -ldflags "$(GO_LDFLAGS)" -o release/mac/intel/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_windows_64: build the reviewbot binary for windows 64
build_windows_64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -a -ldflags "$(GO_LDFLAGS)" -o release/windows/intel/$(EXECUTABLE).exe ./cmd/$(EXECUTABLE)