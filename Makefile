.PHONY: help
help:
	@echo "lint             run lint"
	@echo "release-all        compile for all platforms "

PROJECT=zkcli
VERSION=$(shell cat main.go |grep 'version = "[0-9]\+.[0-9]\+.[0-9]\+"' | awk -F '"' '{print $$2}')

GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))

ARCNAME=$(PROJECT)-$(VERSION)-$(GOOS)-$(GOARCH)
RELDIR=$(ARCNAME)

DISTDIR=dist

.PHONY: release
release:
	rm -rf $(DISTDIR)/$(RELDIR)
	mkdir -p $(DISTDIR)/$(RELDIR)
	go clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build
	cp $(PROJECT)$(SUFFIX_EXE) $(DISTDIR)/$(RELDIR)/
	tar czf $(DISTDIR)/$(ARCNAME).tar.gz -C $(DISTDIR) $(RELDIR)
	go clean

.PHONY: release-all
release-all:
	@$(MAKE) release GOOS=windows GOARCH=amd64 SUFFIX_EXE=.exe
	@$(MAKE) release GOOS=windows GOARCH=386   SUFFIX_EXE=.exe
	@$(MAKE) release GOOS=linux   GOARCH=amd64
	@$(MAKE) release GOOS=linux   GOARCH=386
	@$(MAKE) release GOOS=darwin  GOARCH=amd64

.PHONY: lint
lint:
	gofmt -s -w .
	golint .
	go vet
