NAME     := sorry
VERSION  := 1.1.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: mac
mac:
	rm -rf build/darwin-amd64
	make build OS=darwin ARCH=amd64

.PHONY: build
build:
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o build/$(OS)-$(ARCH)/$(NAME)

bin: $(SRCS)
	go build -a $(LDFLAGS) -o $(NAME)

.PHONY: clean
clean:
	rm -rf build/*

.PHONY: cross
cross: clean linux win64 linux386 win386 mac

.PHONY: linux
linux:
	rm -rf build/linux-amd64
	make build OS=linux ARCH=amd64
.PHONY: win64
win64:
	rm -rf build/windows-amd64
	make build OS=windows ARCH=amd64
.PHONY: linux386
linux386:
	rm -rf build/linux-386
	make build OS=linux ARCH=386
.PHONY: win386
win386:
	rm -rf build/windows-386
	make build OS=windows ARCH=386
