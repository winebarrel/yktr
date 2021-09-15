VERSION := v0.9.5
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

ifeq ($(GOOS),windows)
BIN := yktr.exe
else
BIN := yktr
endif

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/yktr

.PHONY: vet
vet:
	go vet

.PHONY: clean
clean:
	rm -rf yktr yktr.exe dist
