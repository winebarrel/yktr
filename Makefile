VERSION := v0.7.4
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/yktr

.PHONY: vet
vet:
	go vet

.PHONY: package
package: clean vet build
ifeq ($(GOOS),windows)
	zip yktr_$(VERSION)_$(GOOS)_$(GOARCH).zip yktr.exe
	sha1sum yktr_$(VERSION)_$(GOOS)_$(GOARCH).zip > yktr_$(VERSION)_$(GOOS)_$(GOARCH).zip.sha1sum
else
	gzip yktr -c > yktr_$(VERSION)_$(GOOS)_$(GOARCH).gz
	sha1sum yktr_$(VERSION)_$(GOOS)_$(GOARCH).gz > yktr_$(VERSION)_$(GOOS)_$(GOARCH).gz.sha1sum
endif

.PHONY: clean
clean:
	rm -f yktr yktr.exe
