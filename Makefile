format: ## Format go code with goimports
	@go install github.com/rinchsan/gosimports/cmd/gosimports@latest
	@find . -name \*.go -exec gosimports -local github.com/accrete-capital/accrete-backend/ -w {} \;

format-check: ## Check if the code is formatted
	@echo $($GOPATH) $($GOROOT)
	@go install golang.org/x/tools/cmd/goimports@latest
	@for i in $$(goimports -l .); do echo "Code is not formatted run 'make format'" && exit 1; done

check: format format-check ## Linting and static analysis
	@if grep -r --include='*.go' --exclude-dir='vendor' -E "[^\/\/ ]+(fmt.Print|spew.Dump)"  *; then \
		echo "code contains fmt.Print* or spew.Dump function"; \
		exit 1; \
	fi

	@if test ! -e ./bin/golangci-lint; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.50.0; \
	fi
	@./bin/golangci-lint run --timeout 180s

indexer:
	go run console/indexer/main.go

rest:
	go run console/rest/main.go