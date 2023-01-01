.PHONY: test
test:
	go test -v -cover ./...

.PHONY: vendor
vendor:
	rm -rf vendor
	go mod tidy -v
	go mod download
	go mod vendor


.PHONY: validate-commit
validate-commit:
	pre-commit run

.PHONY: lint
lint:
	golangci-lint --version
	golangci-lint --verbose --color=always --timeout=5m run ./...

.PHONY: install-tools
install-tools: install-pre-commit-linter
	go install golang.org/x/tools/cmd/goimports@vlatest

install-pre-commit-linter:
	python3 -m pip install pre-commit
	pre-commit install --hook-type commit-msg

e2e-test:
	go generate tests/e2e/e2e_test.go
	go test --tags=e2e ./tests/...
