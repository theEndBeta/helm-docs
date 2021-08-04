yaml-docs:
	go build github.com/theEndBeta/yaml-doc/cmd/helm-docs

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -f yaml-docs

.PHONY: dist
dist:
	goreleaser release --rm-dist --snapshot --skip-sign
