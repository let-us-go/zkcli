.PHONY: help
help:
	@echo "lint             run lint"
	@echo "release-all      compile for all platforms "

PROJECT=zkcli
VERSION=$(shell cat main.go |grep 'version = "[0-9]\+.[0-9]\+.[0-9]\+"' | awk -F '"' '{print $$2}')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILT_TIME=$(shell date -u '+%FT%T%z')

GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
LDFLAGS="-X main.gitCommit=${GIT_COMMIT} -X main.built=${BUILT_TIME}"

ARCNAME=$(PROJECT)-$(VERSION)-$(GOOS)-$(GOARCH)
RELDIR=$(ARCNAME)

DISTDIR=dist
export GO111MODULE=on

.PHONY: release
release:
	rm -rf $(DISTDIR)/$(RELDIR)
	mkdir -p $(DISTDIR)/$(RELDIR)
	go clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags ${LDFLAGS}
	cp $(PROJECT)$(SUFFIX_EXE) $(DISTDIR)/$(RELDIR)/
	tar czf $(DISTDIR)/$(ARCNAME).tar.gz -C $(DISTDIR) $(RELDIR)
	go clean

.PHONY: release-all
release-all:
	@$(MAKE) release GOOS=linux   GOARCH=amd64
	@$(MAKE) release GOOS=linux   GOARCH=386
	@$(MAKE) release GOOS=darwin  GOARCH=amd64

.PHONY: lint
lint:
	gofmt -s -w .
	golint .
	golint core
	go vet
