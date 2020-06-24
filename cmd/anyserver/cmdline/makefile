#
# anyserver makefile
#
sp :=
sp +=

_walk = $(if $1,$(wildcard /$(subst $(sp),/,$1)/$2) $(call _walk,$(wordlist 2,$(words $1),x $1),$2))
_find = $(firstword $(call _walk,$(strip $(subst /, ,$1)),$2))
_ROOT := $(patsubst %/root.mak,%,$(call _find,$(CURDIR),root.mak))

include $(_ROOT)/root.mak
include $(_ROOT)/top.mak

_GIT := $(_ROOT)/.git

SERVICE := anyserver
SERVICE_EXT := $(_ANYSERVER_EXT)

BUILD_HASH   := $(shell git rev-parse --short HEAD)
BUILD_NUM   := $(shell expr `cat $(_GIT)/buildnum-$(SERVICE) 2>/dev/null` + 1 >$(_GIT)/buildnum-$(SERVICE) && cat $(_GIT)/buildnum-$(SERVICE))
BUILD_DATE  := $(shell date +'%Y-%m-%d.%H:%M:%S')
BUILD_TAG  := $(BUILD_HASH).$(BUILD_NUM)
BUILD_VERSION := $(BUILD_HASH).$(BUILD_NUM) ($(BUILD_DATE))

GOMAIN      := main.go
BINARY      := $(SERVICE)
LDBUILD     := -X github.com/raibru/gsnet/cmd/anyserver/cmdline/cli.buildTag=$(BUILD_TAG)
LDDATE      := -X github.com/raibru/gsnet/cmd/anyserver/cmdline/cli.buildDate=$(BUILD_DATE)
GO_LDFLAGS  := -ldflags="-X github.com/raibru/gsnet/cmd/anyserver/cmdline/cli.buildTag=$(BUILD_TAG) -X github.com/raibru/gsnet/cmd/anyserver/cmdline/cli.buildDate=$(BUILD_DATE)"

#include $(_ROOT)/bottom.mak
include $(_ROOT)/build.mak

#deps:
#				$(GOGET) github.com/markbates/goth
#				$(GOGET) github.com/markbates/pop


# Cross compilation
#build-linux:
#				CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
#
#
#docker-build:
#				docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
#
#compile:
#	# 32-Bit Systems
#	# FreeBDS
#	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
#	# MacOS
#	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 main.go
#	# Linux
#	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
#	# Windows
#	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
#	# 64-Bit Systems
#	# FreeBDS
#	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd-amd64 main.go
#	# MacOS
#	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
#	# Linux
#	GOOS=linux GOARCH=amd64 go build -o bin/main-linux-amd64 main.go
#	# Windows
#	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 main.go
# EOF
