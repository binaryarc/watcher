.PHONY: proto build clean run-server test-local test-remote test-compare test-compare-json test-json test-yaml test-auth help

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
	@go build -o wctl ./cmd/wctl
	@go build -o wsctl ./cmd/wsctl
	@echo "Build complete!"
	@echo "  wctl binary created"
	@echo "  wsctl binary created"

# 빌드 (verbose)
build-verbose:
	@echo "Building binaries (verbose)..."
	go build -v -o wctl ./cmd/wctl
	go build -v -o wsctl ./cmd/wsctl
	@echo "Build complete!"

# 클린
clean:
	@echo "Cleaning up..."
	@rm -f wctl wsctl
	@echo "Clean complete!"

# 클린 (키 파일 포함)
clean-all: clean
	@echo "Cleaning keys..."
	@rm -rf ~/.watcher/keys ~/.watcher/server
	@echo "All clean!"

# 서버 실행
run-server:
	@echo "Starting Watcher server..."
	./wsctl run

# 서버 실행 (인증 비활성화 - 테스트용)
run-server-noauth:
	@echo "Starting Watcher server (auth disabled)..."
	./wsctl run --disable-auth

# 서버 실행 (커스텀 포트)
run-server-custom:
	@echo "Starting Watcher server on port 8080..."
	./wsctl run --port 8080

# === 클라이언트 키 관리 테스트 ===
test-key-gen:
	@echo "Generating API key..."
	./wctl key generate

test-key-get:
	@echo "Getting API key..."
	./wctl get key

# === 서버 키 관리 테스트 ===
test-server-keys-get:
	@echo "Getting server API keys..."
	./wsctl get keys

# === 인증 플로우 테스트 ===
test-auth-setup:
	@echo "Setting up authentication test..."
	@echo "1. Generating client key..."
	@./wctl key generate > /tmp/watcher_key.txt
	@echo ""
	@echo "2. Extracting key..."
	@API_KEY=$$(grep "watcher_" /tmp/watcher_key.txt | head -1 | xargs); \
	echo "   Key: $$API_KEY"; \
	echo ""; \
	echo "3. Adding key to server..."; \
	./wsctl add key "$$API_KEY" "Test key"
	@echo ""
	@echo "Auth setup complete!"
	@echo "   Now start server with: make run-server"
	@echo "   Then test with: make test-remote"

test-auth-teardown:
	@echo "Cleaning up auth test..."
	./wsctl clear keys <<< "yes"
	@rm -f /tmp/watcher_key.txt
	@echo "Auth teardown complete!"

# === 기존 테스트 ===
test-local:
	@echo "Testing local observation..."
	@./wctl get runtimes
	@echo ""
	@./wctl get runtime java

test-remote:
	@echo "Testing remote observation (with auth)..."
	@./wctl get runtimes --host localhost:9090
	@echo ""
	@./wctl get runtime java --host localhost:9090

test-remote-noauth:
	@echo "Testing remote observation (no auth)..."
	@./wctl get runtimes --host localhost:9090 --api-key ""

test-compare:
	@echo "Testing multi-server comparison..."
	@./wctl compare runtimes --hosts localhost:9090,localhost:9091

test-compare-json:
	@echo "Testing comparison with JSON output..."
	@./wctl compare runtimes --hosts localhost:9090,localhost:9091 -o json

test-json:
	@echo "Testing JSON output..."
	@./wctl get runtimes -o json

test-yaml:
	@echo "Testing YAML output..."
	@./wctl get runtimes -o yaml

test-auth: build test-auth-setup
	@echo ""
	@echo "Starting server in background..."
	@./wsctl run > /dev/null 2>&1 & echo $$! > /tmp/watcher_server.pid
	@sleep 2
	@echo ""
	@echo "Running authenticated request..."
	@./wctl get runtimes --host localhost:9090 || true
	@echo ""
	@echo "Stopping server..."
	@kill $$(cat /tmp/watcher_server.pid) 2>/dev/null || true
	@rm -f /tmp/watcher_server.pid
	@echo ""
	@$(MAKE) test-auth-teardown

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
	@echo ""
	@echo "Authentication:"
	@echo "  make test-auth-setup    - Setup authentication test"
	@echo "  make test-auth-teardown - Cleanup authentication test"
	@echo "  make test-auth          - Full authentication flow test"
	@echo ""
	@echo "  make help               - Show this help message"