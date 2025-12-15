.PHONY: proto build build-verbose clean clean-all run-server run-server-noauth run-server-custom test-local test-remote test-remote-noauth test-compare test-compare-json test-json test-yaml test-auth test-auth-setup test-auth-teardown test-key-gen test-key-get test-server-keys-get test-completion completions generate-completions test help

BIN_DIR := $(CURDIR)/bin
COMPLETIONS_DIR := $(CURDIR)/dist/completions
export PATH := $(BIN_DIR):$(PATH)
# 기본 타겟
all: build

# Proto 파일 생성
proto:
	@echo "Generating proto files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/watcher.proto
	@echo "Proto files generated!"

# 빌드
build:
	@echo "Building binaries..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/wctl ./cmd/wctl
	@go build -o $(BIN_DIR)/wsctl ./cmd/wsctl
	@echo "Binaries ready:"
	@echo "  $(BIN_DIR)/wctl"
	@echo "  $(BIN_DIR)/wsctl"
	@$(MAKE) --no-print-directory generate-completions

# 빌드 (verbose)
build-verbose:
	@echo "Building binaries (verbose)..."
	@mkdir -p $(BIN_DIR)
	go build -v -o $(BIN_DIR)/wctl ./cmd/wctl
	go build -v -o $(BIN_DIR)/wsctl ./cmd/wsctl
	@$(MAKE) --no-print-directory generate-completions

# 클린
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@echo "Clean complete!"

# 클린 (키 파일 포함)
clean-all: clean
	@echo "Cleaning keys..."
	@rm -rf ~/.watcher/keys ~/.watcher/server
	@echo "All clean!"

# 서버 실행
run-server:
	@echo "Starting Watcher server..."
	wsctl run

# 서버 실행 (인증 비활성화 - 테스트용)
run-server-noauth:
	@echo "Starting Watcher server (auth disabled)..."
	wsctl run --disable-auth

# 서버 실행 (커스텀 포트)
run-server-custom:
	@echo "Starting Watcher server on port 8080..."
	wsctl run --port 8080

# === 클라이언트 키 관리 테스트 ===
test-key-gen:
	@echo "Generating API key..."
	wctl key gen

test-key-get:
	@echo "Getting API key..."
	wctl get key

# === 서버 키 관리 테스트 ===
test-server-keys-get:
	@echo "Getting server API keys..."
	wsctl get keys

# === 인증 플로우 테스트 ===
test-auth-setup:
	@echo "Setting up authentication test..."
	@echo "1. Generating client key..."
	@wctl key gen > /tmp/watcher_key.txt
	@echo ""
	@echo "2. Extracting key..."
	@API_KEY=$$(grep "watcher_" /tmp/watcher_key.txt | head -1 | xargs); \
	echo "   Key: $$API_KEY"; \
	echo ""; \
	echo "3. Adding key to server..."; \
	wsctl add key "$$API_KEY" "Test key"
	@echo ""
	@echo "Auth setup complete!"
	@echo "   Now start server with: make run-server"
	@echo "   Then test with: make test-remote"

test-auth-teardown:
	@echo "Cleaning up auth test..."
	wsctl clear keys <<< "yes"
	@rm -f /tmp/watcher_key.txt
	@echo "Auth teardown complete!"

# === 기존 테스트 ===
test-local:
	@echo "Testing local observation..."
	@wctl get runtimes
	@echo ""
	@wctl get runtime java

test-remote:
	@echo "Testing remote observation (with auth)..."
	@wctl get runtimes --host localhost:9090
	@echo ""
	@wctl get runtime java --host localhost:9090

test-remote-noauth:
	@echo "Testing remote observation (no auth)..."
	@wctl get runtimes --host localhost:9090 --api-key ""

test-compare:
	@echo "Testing multi-server comparison..."
	@wctl compare runtimes --hosts localhost:9090,localhost:9091

test-compare-json:
	@echo "Testing comparison with JSON output..."
	@wctl compare runtimes --hosts localhost:9090,localhost:9091 -o json

test-json:
	@echo "Testing JSON output..."
	@wctl get runtimes -o json

test-yaml:
	@echo "Testing YAML output..."
	@wctl get runtimes -o yaml

test-auth: build test-auth-setup
	@echo ""
	@echo "Starting server in background..."
	@wsctl run > /dev/null 2>&1 & echo $$! > /tmp/watcher_server.pid
	@sleep 2
	@echo ""
	@echo "Running authenticated request..."
	@wctl get runtimes --host localhost:9090 || true
	@echo ""
	@echo "Stopping server..."
	@kill $$(cat /tmp/watcher_server.pid) 2>/dev/null || true
	@rm -f /tmp/watcher_server.pid
	@echo ""
	@$(MAKE) test-auth-teardown

generate-completions:
	@echo "Generating completion scripts..."
	@rm -rf $(COMPLETIONS_DIR)
	@mkdir -p $(COMPLETIONS_DIR)
	@for shell in bash zsh fish powershell; do \
		case $$shell in \
			powershell) ext=ps1 ;; \
			*) ext=$$shell ;; \
		esac; \
		echo "  wctl ($$shell) -> $(COMPLETIONS_DIR)/wctl.$${ext}"; \
		$(BIN_DIR)/wctl completion $$shell > $(COMPLETIONS_DIR)/wctl.$${ext}; \
		echo "  wsctl ($$shell) -> $(COMPLETIONS_DIR)/wsctl.$${ext}"; \
		$(BIN_DIR)/wsctl completion $$shell > $(COMPLETIONS_DIR)/wsctl.$${ext}; \
	done
	@echo "Completion scripts generated in $(COMPLETIONS_DIR)"

completions:
	@if [ ! -x $(BIN_DIR)/wctl ] || [ ! -x $(BIN_DIR)/wsctl ]; then \
		$(MAKE) build; \
	else \
		$(MAKE) --no-print-directory generate-completions; \
	fi

test-completion:
	@echo "Testing shell completions..."
	@go test ./pkg/cmd/wctl -run TestCompletionCommandGeneratesBashScript -count=1
	@go test ./pkg/cmd/wsctl -run TestCompletionCommandGeneratesBashScript -count=1
	@echo "Completion tests passed!"

test:
	@echo "Running full Go test suite..."
	@go test ./...

# 도움말
help:
	@echo "Watcher Makefile Commands:"
	@echo ""
	@echo "Build & Clean:"
	@echo "  make build              - Build wctl and wsctl binaries"
	@echo "  make build-verbose      - Build with verbose output"
	@echo "  make proto              - Generate proto files"
	@echo "  make clean              - Remove built binaries"
	@echo "  make clean-all          - Remove binaries and key files"
	@echo ""
	@echo "Server:"
	@echo "  make run-server         - Start watcher server on :9090"
	@echo "  make run-server-noauth  - Start server without authentication"
	@echo "  make run-server-custom  - Start server on custom port"
	@echo ""
	@echo "Key Management (Client):"
	@echo "  make test-key-gen       - Generate API key"
	@echo "  make test-key-get       - Get current API key"
	@echo ""
	@echo "Key Management (Server):"
	@echo "  make test-server-keys-get - Get all server keys"
	@echo ""
	@echo "Testing:"
	@echo "  make test-local         - Test local runtime observation"
	@echo "  make test-remote        - Test remote observation (with auth)"
	@echo "  make test-remote-noauth - Test remote observation (no auth)"
	@echo "  make test-compare       - Test multi-server comparison"
	@echo "  make test-json          - Test JSON output format"
	@echo "  make test-yaml          - Test YAML output format"
	@echo "  make test-completion    - Run shell completion tests"
	@echo "  make test               - Run go test ./..."
	@echo ""
	@echo "Authentication:"
	@echo "  make test-auth-setup    - Setup authentication test"
	@echo "  make test-auth-teardown - Cleanup authentication test"
	@echo "  make test-auth          - Full authentication flow test"
	@echo ""
	@echo "  make help               - Show this help message"
