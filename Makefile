.PHONY: docs docs-components docs-test fmt fmt-check lint test hooks release-docs

docs docs-components release-docs:
	./scripts/render-docs.sh

docs-test:
	go test ./internal/docscomponents

fmt:
	find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' -print0 | xargs -0 gofmt -w

fmt-check:
	@output=$$(find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' -print0 | xargs -0 gofmt -l); \
	if [ -n "$$output" ]; then \
		echo "Files need gofmt:"; \
		echo "$$output"; \
		exit 1; \
	fi

lint:
	golangci-lint run ./...

test:
	go test ./cmd/... ./internal/docscomponents ./internal/release

hooks:
	bash ./scripts/install-git-hooks.sh

release:
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION env var required (e.g., make release VERSION=0.1.0)" >&2; \
		exit 1; \
	fi
	go run ./cmd/release --version $(VERSION) $(if $(NOTES),--notes $(NOTES))

release-dry:
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION env var required (e.g., make release-dry VERSION=0.1.0)" >&2; \
		exit 1; \
	fi
	go run ./cmd/release --version $(VERSION) --dry-run $(if $(NOTES),--notes $(NOTES))
