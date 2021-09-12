VERSION := v0.4.1
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/yktr

.PHONY: package
package: clean build
	gzip yktr -c > yktr_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum yktr_$(VERSION)_$(GOOS)_$(GOARCH).gz > yktr_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum

.PHONY: clean
clean:
	rm -f yktr
