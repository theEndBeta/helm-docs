PROJECTNAME=$(shell basename "$(PWD)")

GOBASE=$(shell pwd)
GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN="$(GOBASE)/bin"


.PHONY build
build:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/$(PROJECTNAME) ./cmd/yaml-docs

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: dist
dist:
	goreleaser release --rm-dist --snapshot --skip-sign
