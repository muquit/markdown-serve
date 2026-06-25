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

# check if GITHUB_TOKEN is set and valid, fail the build otherwise
check_github_token:
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
        echo "*** ERROR: GITHUB_TOKEN is not set"; \
        exit 1; \
    fi
	@status=$$(curl -s -o /tmp/check_github_token.$$$$.json -w '%{http_code}' \
        -H "Authorization: token $(GITHUB_TOKEN)" https://api.github.com/user); \
    if [ "$$status" != "200" ]; then \
        echo "*** ERROR: GITHUB_TOKEN is not valid (HTTP $$status)"; \
        cat /tmp/check_github_token.$$$$.json; \
        rm -f /tmp/check_github_token.$$$$.json; \
        exit 1; \
    fi; \
    jq '{login, name, type}' < /tmp/check_github_token.$$$$.json; \
    rm -f /tmp/check_github_token.$$$$.json
	@curl -sI -H "Authorization: token $(GITHUB_TOKEN)" \
        https://api.github.com/user | grep -i x-oauth-scopes

release: check_github_token
	@echo "*** Releasing on github ..."
	go-xbuild-go -release


clean:
	rm -rf ./bin
	rm -f $(BINARY)
