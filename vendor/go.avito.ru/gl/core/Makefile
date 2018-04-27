IMPORT_PATH := go.avito.ru/gl/core
IGNORED_PACKAGES := /vendor/

all: test

.PHONY: help test vet fmt vendor

help:
	@echo "test - run tests"
	@echo "vet  - run vet"
	@echo "fmt   - format application sources"

vendor: .GOPATH/.ok
	$(shell ( cd $(CURDIR)/.GOPATH/src/$(IMPORT_PATH) && GOPATH=$(CURDIR)/.GOPATH govendor sync))

test: fmt vet vendor .GOPATH/.ok
	go test -v $(allpackages)

vet: .GOPATH/.ok
	go vet $(allpackages)

fmt: .GOPATH/.ok
	go fmt $(allpackages)

clean:
	rm -rf bin .GOPATH

# cd into the GOPATH to workaround ./... not following symlinks
_allpackages = $(shell ( cd $(CURDIR)/.GOPATH/src/$(IMPORT_PATH) && \
    GOPATH=$(CURDIR)/.GOPATH go list ./... 2>&1 1>&3 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) 1>&2 ) 3>&1 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)))

# memoize allpackages, so that it's executed only once and only if used
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)

export GOPATH := $(CURDIR)/.GOPATH
unexport GOBIN

.GOPATH/.ok:
	@mkdir -p "$(dir .GOPATH/src/$(IMPORT_PATH))"
	@ln -s ../../../.. ".GOPATH/src/$(IMPORT_PATH)"
	@mkdir -p .GOPATH/test .GOPATH/cover
	@mkdir -p bin
	@ln -s ../bin .GOPATH/bin
	@touch $@
