.PHONY: all
all: build run-build

.PHONY: run
run: flight-check
	@echo "Running the app on port $${PORT:-8080}..."
	@PORT=$${PORT:-8080} go run main.go

.PHONY: build
build:
	@echo "Building the app..."
	@mkdir -p bin
	@go build -o bin/app

.PHONY: run-build
run-build:
	@echo "Running ./bin/app on port $${PORT:-8080}..."
	@PORT=$${PORT:-8080} ./bin/app

.PHONY: test
test: flight-check
	@echo "Running tests..."
	@go test -race -cover ./...

.PHONY: lint
lint: flight-check
	@echo "Running golangci-lint..."
	@./bin/golangci-lint run

.PHONY: clean
clean: flight-check
	@echo "Cleaning up..."
	@rm -rf bin/app

.PHONY: pre-flight
# Clean install; not setting Makefile build targets for the below
pre-flight:
	@echo "Setting up Git hooks..."
	@mkdir -p bin
	@echo "Installing golangci-lint to ./bin..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.63.4
	@cp scripts/git/pre-commit .git/hooks/pre-commit
	@cp scripts/git/pre-push .git/hooks/pre-push
	@chmod +x .git/hooks/pre-commit .git/hooks/pre-push
	@git config core.hooksPath .git/hooks
	@echo "Git hooks set up successfully!"

.PHONY: flight-check
flight-check:
	@echo "Performing flight check..."
	@if [ ! -f ".git/hooks/pre-commit" ] || \
		[ ! -f ".git/hooks/pre-push" ] || \
		[ ! -x "./bin/golangci-lint" ]; then \
			echo "Git hooks or golangci-lint are missing. Run 'make pre-flight'"; \
			exit 1; \
	fi
