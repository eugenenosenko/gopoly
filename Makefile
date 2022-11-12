.PHONY: test all lint vendor

vendor:
	rm -rf vendor
	go mod tidy -v
	go mod download
	go mod vendor

validate-commit:
	pre-commit run

lint:
	golangci-lint --version
	golangci-lint --verbose --color=always --timeout=5m run ./...

install-tools: install-pre-commit-linter
	go install golang.org/x/tools/cmd/goimports@v0.2.0

install-pre-commit-linter:
	python3 -m pip install pre-commit
	pre-commit install --hook-type commit-msg
