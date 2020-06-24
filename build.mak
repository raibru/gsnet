#
#
#

.SUFFIXES:
.SUFFIXES: .go

include $(_ROOT)/color.mak

SHELL := /bin/bash

GOCMD       := /opt/go/bin/go
GOBUILD     := $(GOCMD) build
GOCLEAN     := $(GOCMD) clean
GOTEST      := $(GOCMD) test
GOGET       := $(GOCMD) get
BUILD_DIR   := $(_OUTTOP)

$(_MODULE_NAME)_BINARY := $(_OUTTOP)/$(BINARY)
$(_MODULE_NAME)_LDFLAGS := $(GO_LDFLAGS)

#----------------------------------------

ifneq ($($(_MODULE_NAME)_DEFINED),T)

all: build-$(_MODULE_NAME) test-$(_MODULE_NAME)
 
# Build services
_BUILD := build-$(_MODULE_NAME)

.PHONY: build-$(_BUILD)
build: $(_BUILD)

$(_BUILD):
	@$(GOBUILD) $($(_MODULE_NAME)_LDFLAGS) -o $($(_MODULE_NAME)_BINARY) -v $(_MODULE_PATH)/$(GOMAIN)

# Test services
_TEST := test-$(_MODULE_NAME)

.PHONY: test-$(_TEST)
test: $(_TEST)

$(_TEST):
	@$(GOTEST) -v $(_MODULE_PATH)/... $(COLOR_OUTPUT)

#build-$(_MODULE_NAME): 
#	$(GOBUILD) -o $($(_MODULE_NAME)_BINARY) -v $(_MODULE_PATH)/$(GOMAIN)
#
#test-$(_MODULE_NAME): 
#	$(GOTEST) -v $(_MODULE_PATH)/...

# Smoke test check version
_CHECK := check-$(_MODULE_NAME)

.PHONY: check-$(_CHECK)
check: $(_CHECK)

$(_CHECK):
	@$($(_MODULE_NAME)_BINARY) --version

# Clean build binaries
_CLEAN := clean-$(_MODULE_NAME)

.PHONY: clean $(_CLEAN)
clean: $(_CLEAN)

$(_CLEAN):
	rm -rf $($(_MODULE_NAME)_BINARY)
#	$(GOCLEAN)

$(_MODULE_NAME)_DEFINED := T

endif

