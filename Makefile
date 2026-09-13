# Portable across GNU make on Linux and macOS (make 3.81+).
GO            ?= go
PKG           ?= ./...
FUZZTIME      ?= 30s
COVERPROFILE  ?= coverage.out
GOLANGCI_LINT ?= golangci-lint
GOVULNCHECK   ?= $(GO) run golang.org/x/vuln/cmd/govulncheck@latest

.DEFAULT_GOAL := check
.PHONY: check fmt fmt-check vet lint test cover fuzz bench tidy vuln clean docs docs-check help

check: fmt-check vet lint test ## Run everything CI runs (except fuzz and vuln)

fmt: ## Format sources
	gofmt -w .

fmt-check: ## Fail if any file needs gofmt
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

vet: ## go vet
	$(GO) vet $(PKG)

lint: ## golangci-lint
	$(GOLANGCI_LINT) run $(PKG)

test: ## Tests with race detector and coverage
	$(GO) test -race -cover $(PKG)

cover: ## Coverage HTML report and 100% gate
	$(GO) test -covermode=atomic -coverprofile=$(COVERPROFILE) $(PKG)
	$(GO) tool cover -html=$(COVERPROFILE) -o coverage.html
	@total="$$($(GO) tool cover -func=$(COVERPROFILE) | awk '/^total:/ {sub(/%/, "", $$NF); print $$NF}')"; \
	echo "total coverage: $$total%"; \
	awk -v t="$$total" 'BEGIN { exit (t + 0 < 100.0) ? 1 : 0 }' || \
	{ echo "coverage below 100%:"; $(GO) tool cover -func=$(COVERPROFILE) | grep -v '100.0%'; exit 1; }

fuzz: ## Run every fuzz target for FUZZTIME each
	@set -e; for target in $$(grep -ho '^func Fuzz[A-Za-z0-9_]*' *_test.go | cut -d' ' -f2); do \
		echo "== $$target"; \
		$(GO) test -run='^$$' -fuzz="^$$target\$$" -fuzztime=$(FUZZTIME) .; \
	done

bench: ## Benchmarks with allocations
	$(GO) test -run='^$$' -bench=. -benchmem $(PKG)

tidy: ## go mod tidy
	$(GO) mod tidy

vuln: ## govulncheck
	$(GOVULNCHECK) $(PKG)

clean: ## Remove generated coverage files
	rm -f $(COVERPROFILE) coverage.html

docs: ## Regenerate API pages in docs/ from GoDoc
	python3 scripts/gendocs.py

docs-check: docs ## Fail if docs/ is out of date with GoDoc
	git diff --exit-code -- docs

help: ## List targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "%-10s %s\n", $$1, $$2}'
