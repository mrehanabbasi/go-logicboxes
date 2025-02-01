LINTER_VERSION := v1.63.4
LINTER := golangci/golangci-lint:$(LINTER_VERSION)

tidy:
	go mod tidy

lint:
	docker run --rm -v $(pwd):/app -w /app $(LINTER) golangci-lint run -v --timeout 3m --max-same-issues 0
