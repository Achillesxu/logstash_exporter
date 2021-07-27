.PHONY: help
help:
	@echo "lint             run lint"
	@echo "release-all      compile for all platforms "
	@echo "build            build"

PROJECT=logstash_exporter
VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILT_TIME=$(shell date -u '+%FT%T%z')

GO_VERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(word 4,$(GO_VERSION))))
GOARCH=$(word 2,$(subst /, ,$(word 4,$(GO_VERSION))))

LDFLAGS="-X main.BuildCommitSha=${GIT_COMMIT} -X main.BuildDate=${BUILT_TIME} -X main.BuildVersion=${VERSION}"

ARC_NAME=$(PROJECT)-$(VERSION)-$(GOOS)-$(GOARCH)
RELEASE_DIR=$(ARC_NAME)

ifeq "$(GOOS)" "windows"
	SUFFIX_EXE=".exe"
else
	SUFFIX_EXE=""
endif

DIST_DIR=dist
export GO111MODULE=on

.PHONY: release
release:
	rm -rf $(DIST_DIR)/$(RELEASE_DIR)
	mkdir -p $(DIST_DIR)/$(RELEASE_DIR)
	go clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) make build
	cp $(PROJECT)$(SUFFIX_EXE) $(DIST_DIR)/$(RELEASE_DIR)
	go clean

.PHONY: release-all
release-all:
	@$(MAKE) release GOOS=linux   GOARCH=amd64
	@$(MAKE) release GOOS=linux   GOARCH=arm64
	@$(MAKE) release GOOS=darwin  GOARCH=amd64
	@$(MAKE) release GOOS=windows  GOARCH=amd64

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags ${LDFLAGS}

.PHONY: lint
lint:
	gofmt -s -w .
	go vet
