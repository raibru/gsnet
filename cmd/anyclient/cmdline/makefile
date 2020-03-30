#
# anyclient makefile
#

GOCMD       := go
GOBUILD     := $(GOCMD) build
GOCLEAN     := $(GOCMD) clean
GOTEST      := $(GOCMD) test
GOGET       := $(GOCMD) get
BINARY_NAME := anyclient
BUILD_DIR   := ../../../build/bin
GOMAIN      := main.go

all: build test
build: 
				$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(GOMAIN)
test: 
				$(GOTEST) -v ./...
clean: 
				$(GOCLEAN)
				rm -f $(BUILD_DIR)/$(BINARY_NAME)
run:
				$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(GOMAIN)
				$(BUILD_DIR)/$(BINARY_NAME) --version

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
