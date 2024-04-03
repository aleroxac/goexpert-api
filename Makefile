## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@sudo rm -f *.db coverage.* 



## ---------- SETUP
.PHONY: install
install: ## install requirements
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go mod tidy



## ---------- MAIN
.PHONY: fmt
fmt: ## format the code
	@go fmt ./...

.PHONY: vet
vet: ## run static analysis
	@go vet ./...

.PHONY: docs
docs: ## generate/update swagger docs
	@swag init -g cmd/server/main.go || true
	@swag fmt



## ---------- TESTS
.PHONY: test
test: ## run unit-tests
	@go test -v ./... -coverprofile coverage.out
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: bench
bench: ## run benchmarks
	@go test -v -benchmem ./...



## ---------- MAIN
.PHONY: run
run: ## run the app
	@cd cmd/server && go run main.go

