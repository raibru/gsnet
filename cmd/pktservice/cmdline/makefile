#
# pktservice makefile
#
sp :=
sp +=

_walk = $(if $1,$(wildcard /$(subst $(sp),/,$1)/$2) $(call _walk,$(wordlist 2,$(words $1),x $1),$2))
_find = $(firstword $(call _walk,$(strip $(subst /, ,$1)),$2))
_ROOT := $(patsubst %/root.mak,%,$(call _find,$(CURDIR),root.mak))

include $(_ROOT)/root.mak
include $(_ROOT)/top.mak

SERVICE := pktservice
SERVICE_EXT := $(_PKTSERVICE_EXT)

GOMAIN      := main.go
BINARY      := $(SERVICE)

LDBUILD     := -X "github.com/raibru/gsnet/cmd/pktservice/cmdline/cli.Build"
LDVERSION   := -X "github.com/raibru/gsnet/cmd/pktservice/cmdline/cli.Version"

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
