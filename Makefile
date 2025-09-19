.PHONY: docs docs-components docs-test

docs docs-components:
	./scripts/render-docs.sh

docs-test:
	go test ./internal/docscomponents
