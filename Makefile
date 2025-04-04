#!/usr/bin/make -f

# Available commands:
# make proto                  - Generate Go protobuf files only
# make proto PYTHON_GEN=true  - Generate Go and Python protobuf files
# make proto PULSAR_GEN=true  - Generate Go and Pulsar protobuf files
# make proto PYTHON_GEN=true PULSAR_GEN=true  - Generate all protobuf files (Go, Python, and Pulsar)
#
# Note: If you hit rate limits with buf.build, wait a few minutes before retrying.
# To avoid rate limits, you can set up a local cache by adding dependencies to buf.yaml
# and running 'buf mod update' before generating files.

PACKAGES_NOSIMULATION=$(shell go list ./... | grep -Ev 'vendor|importer|wasm|simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || echo $(shell git describe --tags git rev-list --tags="v*" --max-count=1) | sed 's/^v//')
VERSION_DATE ?= $(shell git show -s --format=%ci)
TMVERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
NODE_BINARY = zenrockd
NODE_DIR = zenrock
BUILDDIR ?= $(CURDIR)/build
SIMAPP = ./app
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.0.0-rc8
PYTHON_GEN ?= false
PULSAR_GEN ?= false



export GO111MODULE = on
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X "github.com/cosmos/cosmos-sdk/version.Name=zenrock" \
		  -X "github.com/cosmos/cosmos-sdk/version.AppName=$(NODE_BINARY)" \
		  -X "github.com/cosmos/cosmos-sdk/version.Version=$(VERSION)" \
		  -X "github.com/cosmos/cosmos-sdk/version.Commit=$(commit_hash)" \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X "github.com/tendermint/tendermint/version.TMCoreSemVer=$(TMVERSION)" \
	      -X "github.com/zenrocklabs/zenrock/zrchain/version.linkedDate=$(build_date)" \
		  -X "github.com/zenrocklabs/zenrock/zrchain/version.linkedSemVer=$(version)" \
		  -X "github.com/zenrocklabs/zenrock/zrchain/version.linkedCommit=$(commit_hash)"

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  CGO_ENABLED=1
  BUILD_TAGS += rocksdb
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  BUILD_TAGS += boltdb
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/
build-linux:
	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/


$(MOCKS_DIR):
	mkdir -p $(MOCKS_DIR)

distclean: clean tools-clean

clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/

all: build

build-all: tools build lint test

.PHONY: distclean clean build-all

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=latest
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --user root --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

# ------
# NOTE: If you are experiencing problems running these commands, try deleting
#       the docker images and execute the desired command again.
#
proto-all: proto-format proto-lint proto-gen

# Add new variables for Python and Pulsar generation

proto proto-gen:
	@echo "Generating Protobuf files"
	go run ./protogen $(if $(filter true,$(PYTHON_GEN)),--python) $(if $(filter true,$(PULSAR_GEN)),--pulsar)
ifeq ($(PYTHON_GEN),true)
	@echo "Generating Python files"
	./scripts/protocgen.sh --python
endif
ifeq ($(PULSAR_GEN),true)
	@echo "Generating Pulsar files"
	./scripts/pulsargen.sh --pulsar
endif

proto-lint:
	@echo "Linting Protobuf files"
	@$(protoImage) buf lint --error-format=json

proto-check-breaking:
	@echo "Checking Protobuf files for breaking changes"
	@$(protoImage) buf breaking --against $(HTTPS_GIT)#branch=main

# TODO: Rethink API docs generation
proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

.PHONY: proto-all proto proto-gen proto-format proto-lint proto-check-breaking


###############################################################################
###                                Web                                      ###
###############################################################################

web-gen:
	@echo "Generating web client and hooks from Protobuf files"
	sh ./scripts/webgen.sh

.PHONY: web-gen

###############################################################################
###                              Sidecar                                    ###
###############################################################################

build-sidecar:
	go build -o sidecar-new ./sidecar
	rm -f sidecar/sidecar
	mv sidecar-new sidecar/sidecar
	chmod +x sidecar/sidecar

run-sidecar:
	cd sidecar && ./sidecar

run-alt-sidecar:
	cd sidecar && ./sidecar --port 9393 --cache-file cache_alt.json --neutrino-port 54321 --neutrino-path /neutrino_alt_

sidecar: build-sidecar run-sidecar

alt-sidecar: build-sidecar run-alt-sidecar

.PHONY: build-sidecar run-sidecar sidecar alt-sidecar run-alt-sidecar
