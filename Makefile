.PHONY: all run test lint build run-build clean pre-flight flight-check

all: flight-check build run-build

run: flight-check
	@echo "Running the app on port $${PORT:-8080}..."
	@PORT=$${PORT:-8080} go run main.go

build: flight-check
	@echo "Building the app..."
	@if [ ! -d "bin" ]; then mkdir -f bin/; fi
	@go build -o bin/app

run-build: flight-check
	@echo "Running ./bin/app on port $${PORT:-8080}..."
	@PORT=$${PORT:-8080} ./bin/app

test: flight-check
	@echo "Running tests..."
	@go test -race -cover ./...

lint: flight-check
	@echo "Running golangci-lint..."
	@./bin/golangci-lint run

clean: flight-check
	@echo "Removing /bin/app..."
	@rm -rf bin/app

pre-flight:
	@echo "Setting up Git hooks..."
	@if [ ! -d "bin" ]; then mkdir -f bin/; fi
	@echo "Installing golangci-lint to ./bin..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.63.4
	@cp scripts/git/pre-commit .git/hooks/pre-commit
	@cp scripts/git/pre-push .git/hooks/pre-push
	@chmod +x .git/hooks/pre-commit .git/hooks/pre-push
	@git config core.hooksPath .git/hooks
	@echo "Git hooks set up successfully!"

flight-check:
	@if [ ! -f ".git/hooks/pre-commit" ] || \
		[ ! -f ".git/hooks/pre-push" ] || \
		! command -v ./bin/golangci-lint >/dev/null 2>&1; then \
			echo "Git hooks or golangci-lint are missing. Run 'make pre-flight'"; \
			exit 1; \
	fi
