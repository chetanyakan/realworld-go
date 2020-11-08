GO ?= $(shell command -v go 2> /dev/null)
GOARCH := amd64
GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOPATH ?= $(shell go env GOPATH)
GO_TEST_FLAGS ?= -race
GO_BUILD_FLAGS ?= -mod=vendor
DLV_DEBUG_PORT := 2346

TMPFILEGOLINT=golint.tmp

export GO111MODULE=on

# You can include assets this directory into the bundle. This can be e.g. used to include profile pictures.
ASSETS_DIR ?= assets


## Define the default target (make all)
.PHONY: default
default: all

## Checks the code style, tests, builds and bundles the plugin.
.PHONY: all
all: check-style test build


## Checks for style guide compliance.
.PHONY: check-style
check-style: gofmt govet golint


## Runs gofmt against all packages.
.PHONY: gofmt
gofmt:
	@echo Running GOFMT
	@for package in $$(go list ./...); do \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failure\n"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success\n"


## Runs govet against all packages.
.PHONY: govet
govet:
	@echo Running GOVET
	@# Workaround because you can't install binaries without adding them to go.mod
	@env GO111MODULE=off $(GO) get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	@$(GO) vet ./...
	@$(GO) vet -vettool=$(GOPATH)/bin/shadow ./...
	@echo "govet success\n"


## Runs go lint against all packages.
.PHONY: golint
golint:
	@echo Running GOLINT
	env GO111MODULE=off $(GO) get golang.org/x/lint/golint
	$(eval PKGS := $(shell go list ./... | grep -v /vendor/))
	@touch $(TMPFILEGOLINT)
	-@for pkg in $(PKGS) ; do \
		echo `$(GOPATH)/bin/golint $$pkg | grep -v "have comment" | grep -v "comment on exported" | grep -v "lint suggestions"` >> $(TMPFILEGOLINT) ; \
	done
	@grep -Ev "^$$" $(TMPFILEGOLINT) || true
	@if [ "$$(grep -Ev "^$$" $(TMPFILEGOLINT) | wc -l)" -gt "0" ]; then \
		rm -f $(TMPFILEGOLINT); echo "golint failure\n"; exit 1; else \
		rm -f $(TMPFILEGOLINT); echo "golint success\n"; \
	fi


## Runs golangci-lint.
.PHONY: golangci-lint
golangci-lint:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi; \
	echo Running golangci-lint
	golangci-lint run ./...


## Builds the server.
.PHONY: build
build:
	mkdir -p dist;
	go build $(GO_BUILD_FLAGS) -o ./dist/realworld-go .


## Runs the server.
.PHONY: run
run: build
	./dist/realworld-go


## Builds the server in DEBUG mode.
.PHONY: server-debug
server-debug:
	$(info DEBUG mode is on)
	mkdir -p dist;
	go build $(GO_BUILD_FLAGS) -gcflags "all=-N -l" -o ./dist/realworld-go .


## Runs any unit tests, if they exist.
.PHONY: test
test:
	$(GO) test -v $(GO_TEST_FLAGS) ./...


## Creates a coverage report for the code.
.PHONY: coverage
coverage:
	$(GO) test $(GO_TEST_FLAGS) -coverprofile=coverage.txt ./...
	$(GO) tool cover -html=server/coverage.txt


## Removes all dependencies and build-artifacts.
.PHONY: clean
clean:
	rm -fr vendor
	rm -fr dist/
	rm -fr coverage.txt


# Help documentation à la https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@cat Makefile build/*.mk | grep -v '\.PHONY' |  grep -v '\help:' | grep -B1 -E '^[a-zA-Z0-9_.-]+:.*' | sed -e "s/:.*//" | sed -e "s/^## //" |  grep -v '\-\-' | sed '1!G;h;$$!d' | awk 'NR%2{printf "\033[36m%-30s\033[0m",$$0;next;}1' | sort
