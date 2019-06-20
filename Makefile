NAME     := sorry
VERSION  := 1.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: mac
mac:
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
cross: linux win64 linux386 win386 mac

.PHONY: linux
linux:
	make build OS=linux ARCH=amd64
.PHONY: win64
win64:
	make build OS=windows ARCH=amd64
.PHONY: linux386
linux386:
	make build OS=linux ARCH=386
.PHONY: win386
win386:
	make build OS=windows ARCH=386
