VERSION := v0.9.0
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

.PHONY: package
package: clean vet build
ifeq ($(GOOS),darwin)
	codesign -s $(CODESIGN_ID) -o runtime -v ./yktr
endif
	zip yktr_$(VERSION)_$(GOOS)_$(GOARCH).zip $(BIN)
	sha1sum yktr_$(VERSION)_$(GOOS)_$(GOARCH).zip > yktr_$(VERSION)_$(GOOS)_$(GOARCH).zip.sha1sum

.PHONY: notarize
notarize:
ifeq ($(GOOS),darwin)
	xcrun altool \
		--notarize-app \
		--primary-bundle-id com.github.winebarrel.yktr \
		--username sugawara@winebarrel.jp \
		--password "$(ALTOOL_PASSWORD)" \
		--file yktr_$(VERSION)_$(GOOS)_$(GOARCH).zip
endif

.PHONY: clean
clean:
	rm -f yktr yktr.exe
