BINARY  := markdown-serve
VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-w -s -X main.version=$(VERSION)"
BUILD_OPTIONS := -trimpath

.PHONY: all build clean docs

all: build

build:
	@echo ">>>> Compiling native binary ..."
	go build $(BUILD_OPTIONS) $(LDFLAGS) -o $(BINARY) .
	@echo ""
	@echo ">>>> Compiling cros-platform binaries ..."
	go-xbuild-go \
		-additional-files "mdsr.sh,mdsr.ps1" \
		-build-args '$(BUILD_OPTIONS) $(LDFLAGS)'

docs:
	markdown-toc-go -i docs/README.md \
        -o ./README.md --glossary docs/glossary.txt -f
	markdown-toc-go -i docs/ChangeLog.md -o ./ChangeLog.md --glossary docs/glossary.txt -f -no-credit

clean:
	rm -rf ./bin
	rm -f $(BINARY)
