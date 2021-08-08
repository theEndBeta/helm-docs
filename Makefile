PROJECTNAME=$(shell basename "$(PWD)")

GOBASE=$(shell pwd)
GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN="$(GOBASE)/bin"


.PHONY: build
build:
	mkdir -p $(GOBIN)
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o "$(GOBIN)/$(PROJECTNAME)" ./cmd/yaml-docs

.PHONY: vendor
vendor:
	mkdir -p "$(GOBASE)/vendor"
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod vendor

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf $(GOBIN)/*
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean -i

.PHONY: dist
dist:
	goreleaser release --rm-dist --snapshot --skip-sign
