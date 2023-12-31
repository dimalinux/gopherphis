GOPATH ?= $(shell go env GOPATH)
LINTER_VERSION = "1.54.2"
LINTER = $(GOPATH)/bin/golangci-lint

.PHONY: all
all: format lint test

.PHONY: format
format:
	test -x "$(GOPATH)/bin/goimports" || go install golang.org/x/tools/cmd/goimports@latest
	"$(GOPATH)/bin/goimports" -local github.com/dimalinux/gopherphis -w .


# Install golangci-lint if it is not already installed. See here for details:
# https://golangci-lint.run/usage/install/#local-installation
$(LINTER):
	curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" \
		| sh -s -- -b "$(GOPATH)/bin" "v$(LINTER_VERSION)"

# If the version test below fails, delete the executable so that it can be
# reinstalled with the correct version.
.PHONY: lint
lint: $(LINTER)
	@echo test "$$("$(LINTER)" version --format=short)" = "$(LINTER_VERSION)" || \
		($$(@echo "Please delete the out-of-date $(LINTER) executable") && false)
	"$(LINTER)" run

.PHONY: test
test: 
	go test -coverpkg=./... -v -covermode=atomic -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html
