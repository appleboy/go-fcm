GO ?= go
GOFILES := $(shell find . -name "*.go" -type f)
GOFMT ?= gofmt "-s"

all: build

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test: fmt-check
	@$(GO) test -v -cover -coverprofile coverage.out ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

clean:
	go clean -x -i ./...
