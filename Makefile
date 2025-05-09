GO ?= go
EXECUTABLE := reviewbot
GOFILES := $(shell find . -type f -name "*.go")

## build: build the reviewbot binary
build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	$(GO) build -v -o bin/$@ ./cmd/$(EXECUTABLE)

## build_linux_amd64: build the reviewbot binary for linux amd64
build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -o release/linux/amd64/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_linux_arm64: build the reviewbot binary for linux arm64
build_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -a -o release/linux/arm64/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_linux_arm: build the reviewbot binary for linux arm
build_linux_arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GO) build -a -o release/linux/arm/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_mac_intel: build the reviewbot binary for mac intel
build_mac_intel:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -a -o release/mac/intel/$(EXECUTABLE) ./cmd/$(EXECUTABLE)

## build_windows_64: build the reviewbot binary for windows 64
build_windows_64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -a -o release/windows/intel/$(EXECUTABLE).exe ./cmd/$(EXECUTABLE)