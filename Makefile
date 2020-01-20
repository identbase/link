# Commands
GO_CMD = `which go`
LINT_CMD = $(GOPATH)/bin/golint
DOCKER_CMD = `which docker`

# Directories
PACKAGE = github.com/identbase/identserv
SRC = $(GOPATH)/src/$(PACKAGE)
DIST = $(SRC)/dist
BINARY = identsrv

default: lint test clean compile

build: lint test clean compile

run: lint test go

dependencies:
	cat $(SRC)/requirements.txt | xargs -I \\# go get -u github.com/\\#

lint:
	$(LINT_CMD) ./...

test:
	cd $(SRC)
	$(GO_CMD) test -v ./...
	# $(GOPATH)/bin/goveralls -coverprofile=coverage.out -service=travis-ci # -repotoken $(COVERALLS_TOKEN)

# coverage:
# 	cd $(SRC)
# 	$(GOPATH)/bin/overalls -project=$(PACKAGE) -covermode=count
# 	$(GOPATH)/bin/goveralls -coverprofile=overalls.coverprofile -service=travis-ci

clean:
	rm -rf $(DIST)/$(BINARY)

compile:
	mkdir -p $(DIST)
	cd $(SRC)
	$(GO_CMD) build -o $(DIST)/daemon

go:
	cd $(SRC)
	$(GO_CMD) run main.go

