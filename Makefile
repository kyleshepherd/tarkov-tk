# General
CMD      :=
PKG      := github.com/kyleshepherd/discord-tk-bot
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

# Docker
DOCKER_REGISTRY = gcr.io/tarkov-tk-dev/discord-bot

# Versioning
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_SHA    ?= $(shell git rev-parse --short HEAD)
GIT_TAG    ?= $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  ?= $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

# Binary Name
BIN_OUTDIR ?= ./build/bin
BIN_NAME   ?= tarkov-tk-bot-$(shell go env GOOS)-$(shell go env GOARCH)
ifeq ("$(GIT_TAG)","")
	BIN_VERSION = $(GIT_SHA)
endif
BIN_VERSION ?= ${GIT_TAG}

# Docker Tag from Git
DOCKER_TAG  ?= ${GIT_TAG}
ifeq ("$(DOCKER_TAG)","")
	DOCKER_TAG = $(GIT_SHA)
endif

# LDFlags
# LDFLAGS := -w -s
LDFLAGS += -X $(PKG)/internal/version.Timestamp=$(shell date +%s)
LDFLAGS += -X $(PKG)/internal/version.GitCommit=${GIT_COMMIT}
LDFLAGS += -X $(PKG)/internal/version.GitTreeState=${GIT_DIRTY}
LDFLAGS += -X $(PKG)/internal/version.Version=${BIN_VERSION}

# CGO
CGO ?= 0

# Go Build Flags
GOBUILDFLAGS :=
GOBUILDFLAGS += -o $(BIN_OUTDIR)/$(BIN_NAME)

.PHONY: info
info:
	@echo "Version:        ${BIN_VERSION}"
	@echo "Binary Name:    ${BIN_NAME}"
	@echo "Git Tag:        ${GIT_TAG}"
	@echo "Git Commit:     ${GIT_COMMIT}"
	@echo "Git Tree State: ${GIT_DIRTY}"

# Build a statically linked binary
.PHONY: static
static: CGO = 0
static: GOBUILDFLAGS += -a
static: GOBUILDFLAGS += -tags netgo -installsuffix netgo
static: GOBUILDFLAGS += -installsuffix netgo
static: LDFLAGS += -extldflags "-static"
static: build

# Build a binary
.PHONY: build
build: CMD = ./cmd/tarkov-tk-bot
build: GOBUILDFLAGS += -ldflags '$(LDFLAGS)'
build:
	@CGO_ENABLED=$(CGO) go build $(GOBUILDFLAGS) $(CMD)

# Build and run the application
.PHONY: run
run: GOBUILDFLAGS += -i
run: build
	@$(BIN_OUTDIR)/$(BIN_NAME)

# Build docker image
.PHONY: image
image:
	docker build \
		--force-rm \
		--build-arg BIN_VERSION=$(BIN_VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_SHA=$(GIT_SHA) \
		--build-arg GIT_TAG=$(GIT_TAG) \
		--build-arg GIT_DIRTY=$(GIT_DIRTY) \
		--build-arg GOPROXY \
		--build-arg GONOSUMDB \
		-f ./build/package/Dockerfile \
		-t $(DOCKER_REGISTRY):$(DOCKER_TAG) .

# Build docker image
.PHONY: imagex
imagex:
	docker buildx build \
		--platform linux/amd64 \
		--force-rm \
		--build-arg BIN_VERSION=$(BIN_VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_SHA=$(GIT_SHA) \
		--build-arg GIT_TAG=$(GIT_TAG) \
		--build-arg GIT_DIRTY=$(GIT_DIRTY) \
		--build-arg GOPROXY \
		--build-arg GONOSUMDB \
		-f ./build/package/Dockerfile \
		-t $(DOCKER_REGISTRY):$(DOCKER_TAG) .

# Run test suite
.PHONY: test
test:
ifeq ("$(wildcard $(shell which gocov))","")
	go get github.com/axw/gocov/gocov
endif
	gocov test -race ${PKG_LIST} | gocov report
