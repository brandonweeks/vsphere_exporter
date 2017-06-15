GO            ?= GO15VENDOREXPERIMENT=1 go
GOPATH        := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))

PROMU         ?= $(GOPATH)/bin/promu
GOMETALINTER  ?= $(GOPATH)/bin/gometalinter
pkgs           = $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX        ?= $(shell pwd)
BIN_DIR       ?= $(shell pwd)
TARGET         = "vsphere_exporter"

DOCKER_IMAGE_NAME       ?= brandonweeks/vsphere_exporter
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

all: format vet build test

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

gometalinter: $(GOMETALINTER)
	@echo ">> linting code"
	$(GOMETALINTER) --vendor

build: $(PROMU)
	@echo ">> building binaries"
	@CGO_ENABLED=0 $(PROMU) build --prefix $(PREFIX)

test:
	@echo ">> running tests"
	@$(GO) test -short $(pkgs)

docker:
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

tarball: $(PROMU)
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

clean:
	@echo ">> Cleaning up"
	@$(RM) $(TARGET) *.tar.gz *~

$(GOPATH)/bin/promu promu:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
		GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
		$(GO) get -u github.com/prometheus/promu

$(GOPATH)/bin/gometalinter lint:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
		GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
		$(GO) get -u github.com/alecthomas/gometalinter
	@$(GOMETALINTER) --install >/dev/null


.PHONY: all clean style format build test vet tarball docker $(GOPATH)/bin/promu promu $(GOPATH)/bin/gometalinter lint
