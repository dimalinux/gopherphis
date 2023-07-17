GOPATH ?= $(shell go env GOPATH)
LINTER_VERSION = "v1.53.3"
LINTER = $(GOPATH)/bin/golangci-lint

.PHONY: all
all: format lint test

.PHONY: format
format:
	test -x $(GOPATH)/bin/goimports || go install golang.org/x/tools/cmd/goimports@latest
	$(GOPATH)/bin/goimports -local github.com/dimalinux/gopherphis -w .


# See: https://golangci-lint.run/usage/install/#local-installation
# Installs the linter in the GOPATH, if an only if it does not already
# exist or is on the wrong version, before running it.
.PHONY: lint
lint: 
	! test -x "$${LINTER}" || \
		test "v$$("$${LINTER}" version --format=short)" != "$${VERSION}" || \
		curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b "$${GOPATH}/bin" "$${LINTER_VERSION}"
	${GOPATH}/bin/golangci-lint run

.PHONY: test
test: 
	go test -v ./...
