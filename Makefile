VERSION = $(shell cat VERSION)
GO ?= go
GOFMT ?= $(GO)fmt
DOCKER_ORG = bluesteelabm

FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
DEFAULT_GOPATH = $(shell echo $$GOPATH|tr ':' '\n'|awk '!x[$$0]++'|sed '/^$$/d'|head -1)
ifeq ($(DEFAULT_GOPATH),)
DEFAULT_GOPATH := ~/go
endif
DEFAULT_GOBIN = $(DEFAULT_GOPATH)/bin
export PATH := $(PATH):$(DEFAULT_GOBIN)

GOLANGCI_LINT = $(DEFAULT_GOBIN)/golangci-lint
RICH_GO = $(DEFAULT_GOBIN)/richgo
GODA = $(DEFAULT_GOBIN)/goda

DVCS_HOST = github.com
ORG = bluest-eel
DOCKER_ORG = bluesteelabm
PROJ = state
FQ_PROJ = $(DVCS_HOST)/$(ORG)/$(PROJ)

LD_VERSION = -X $(FQ_PROJ)/common.version=$(VERSION)
LD_BUILDDATE = -X $(FQ_PROJ)/common.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LD_GITCOMMIT = -X $(FQ_PROJ)/common.gitCommit=$(shell git rev-parse --short HEAD)
LD_GITBRANCH = -X $(FQ_PROJ)/common.gitBranch=$(shell git rev-parse --abbrev-ref HEAD)
LD_GITSUMMARY = -X $(FQ_PROJ)/common.gitSummary=$(shell git describe --tags --dirty --always)

LDFLAGS = -w -s $(LD_VERSION) $(LD_BUILDDATE) $(LD_GITBRANCH) $(LD_GITSUMMARY) $(LD_GITCOMMIT)

default: all

all-pre-test: clean lint
all-post-test: build
all: all-pre-test test all-post-test
all-cicd: all-pre-test test-nocolor all-post-test

#############################################################################
###   Source Code   #########################################################
#############################################################################
###
### Linting, building, testing, etc.
###

show-version:
	@echo $(VERSION)

deps:
	@echo '>> Downloading deps ...'
	@$(GO) get -v -d ./...

$(GOLANGCI_LINT):
	@echo ">> Couldn't find $(GOLANGCI_LINT); installing ..."
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | \
	sh -s -- -b $(DEFAULT_GOBIN) v1.15.0

show-linter:
	@echo $(GOLANGCI_LINT)

lint: $(GOLANGCI_LINT)
	@echo '>> Linting source code'
	@GL_DEBUG=linters_output GOPACKAGESPRINTGOLISTERRORS=1 $(GOLANGCI_LINT) \
	--enable=golint \
	--enable=gocritic \
	--enable=misspell \
	--enable=nakedret \
	--enable=unparam \
	--enable=lll \
	--enable=goconst \
	run ./...

$(RICH_GO):
	@echo ">> Couldn't find $(RICH_GO); installing ..."
	@GOPATH=$(DEFAULT_GOPATH) \
	GOBIN=$(DEFAULT_GOBIN) \
	GO111MODULE=on \
	$(GO) get -u github.com/kyoh86/richgo

test: $(RICH_GO)
	@echo '>> Running all tests'
	@$(RICH_GO) test ./... -v

test-nocolor:
	@echo '>> Running all tests'
	@$(GO) test ./... -v

bin:
	@mkdir ./bin

TOOL = tool
bin/$(TOOL): bin
	@echo '>> Building tool binary'
	@$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(TOOL) ./cmd/$(TOOL)

build-tool: | bin/$(TOOL)
build: build-tool
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(MAKE) build

clean:
	@echo '>> Removing project binaries ...'
	@rm -f bin/$(TOOL)

#############################################################################
###   Infrastructure   ######################################################
#############################################################################
###
### Docker, docker-compose, etc., for local dev work
###

up:
	@VERSION=$(VERSION) RELAYD_VERSION=$(WS_RELAY_VERSION) \
	docker-compose -f ./docker/dev/compose.yml up

rebuild:
	@VERSION=$(VERSION) RELAYD_VERSION=$(WS_RELAY_VERSION) \
	docker-compose -f ./docker/dev/compose.yml build

rebuild-up: rebuild up

down:
	@VERSION=$(VERSION) RELAYD_VERSION=$(WS_RELAY_VERSION) \
	docker-compose -f ./docker/dev/compose.yml down

sqlsh: NODE ?= db1
sqlsh:
	@echo '>> Connecting to db $(NODE) ...'
	@docker exec -it $(NODE) ./cockroach sql --insecure

bash: NODE ?= db1
bash:
	@docker exec -it $(NODE) bash

$(WS_RELAY_CODE_DIR):
	@cd $(WS_RELAY_DIR) && \
	git clone $(WS_RELAY_REPO) $(WS_RELAY_CODE_NAME) && \
	git checkout v$(WS_RELAY_VERSION)

$(WS_RELAY_RENAME): $(WS_RELAY_CODE_DIR)
	@docker build -t $(DOCKER_ORG)/$(WS_RELAY_RENAME):$(WS_RELAY_VERSION) $(WS_RELAY_DIR)

images: $(WS_RELAY_RENAME)

tags:
	@docker tag $(DOCKER_ORG)/$(WS_RELAY_RENAME):$(WS_RELAY_VERSION) \
	$(DOCKER_ORG)/$(WS_RELAY_RENAME):latest

dockerhub: tags
	@docker push $(DOCKER_ORG)/$(WS_RELAY_RENAME):$(WS_RELAY_VERSION)
	@docker push $(DOCKER_ORG)/$(WS_RELAY_RENAME):latest

clean-docker:
	@docker system prune -f

#############################################################################
###   Release Process   #####################################################
#############################################################################

tag:
	@echo "Tags:"
	@git tag|tail -5
	@git tag "v$(VERSION)"
	@echo "New tag list:"
	@git tag|tail -6

#############################################################################
###   Misc   ################################################################
#############################################################################

clean-cache:
	@echo '>> Purging Go mod cahce ...'
	# @$(GO) clean -cache
	@$(GO) clean -modcache

show-targets:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | \
	awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | \
	sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

check-modules:
	@echo '>> Checking modules ...'
	@GO111MODULE=on $(GO) mod tidy
	@#@GO111MODULE=on $(GO) mod verify

$(GODA):
	@echo ">> Couldn't find $(GODA); installing ..."
	@GOPATH=$(DEFAULT_GOPATH) \
	GOBIN=$(DEFAULT_GOBIN) \
	GO111MODULE=on \
	$(GO) get -u github.com/loov/goda

deps-tree: $(GODA)
	@GO111MODULE=on $(GODA) tree ./...

deps-graph: $(GODA)
	@GO111MODULE=on $(GODA) graph ./... | dot -Tsvg -o graph.svg

show-ldflags:
	@echo $(LDFLAGS)