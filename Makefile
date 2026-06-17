BINARY  := markdown-serve
VERSION := $(shell cat VERSION)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: all build clean docs

all: build

build:
	go build $(LDFLAGS) -o $(BINARY) .

docs:
	markdown-toc-go -i docs/README.md \
        -o ./README.md --glossary docs/glossary.txt -f

clean:
	rm -f $(BINARY)
