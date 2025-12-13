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
	@go build -o watcher-server ./cmd/watcher-server
	@echo "Build complete!"
	@echo "  wctl binary created"
	@echo "  watcher-server binary created"

# 빌드 (verbose)
build-verbose:
	@echo "Building binaries (verbose)..."
	go build -v -o wctl ./cmd/wctl
	go build -v -o watcher-server ./cmd/watcher-server
	@echo "Build complete!"

# 클린
clean:
	@echo "Cleaning up..."
	@rm -f wctl watcher-server
	@echo "Clean complete!"

# 클린 (키 파일 포함)
clean-all: clean
	@echo "Cleaning keys..."
	@rm -rf ~/.watcher/keys ~/.watcher/server
	@echo "All clean!"

# 서버 실행
run-server:
	@echo "Starting Watcher server..."
	./watcher-server run

# 서버 실행 (인증 비활성화 - 테스트용)
run-server-noauth:
	@echo "Starting Watcher server (auth disabled)..."
	./watcher-server run --disable-auth

# 서버 실행 (커스텀 포트)
run-server-custom:
	@echo "Starting Watcher server on port 8080..."
	./watcher-server run --port 8080

# === 클라이언트 키 관리 테스트 ===
test-key-gen:
	@echo "Generating API key..."
	./wctl key generate

test-key-list:
	@echo "Listing API keys..."
	./wctl key list

test-key-show:
	@echo "Showing default API key..."
	./wctl key show

# === 서버 키 관리 테스트 ===
test-server-key-list:
	@echo "Listing server API keys..."
	./watcher-server key list

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
	./watcher-server key add "$$API_KEY" "Test key"
	@echo ""
	@echo "Auth setup complete!"
	@echo "   Now start server with: make run-server"
	@echo "   Then test with: make test-remote"

test-auth-teardown:
	@echo "Cleaning up auth test..."
	./watcher-server key clear
	@rm -f /tmp/watcher_key.txt
	@echo "Auth teardown complete!"

# === 기존 테스트 (인증 포함) ===
# 로컬 테스트
test-local:
	@echo "Testing local observation..."
	@./wctl get runtimes
	@echo ""
	@./wctl get runtime java

# 원격 테스트 (인증 사용)
test-remote:
	@echo "Testing remote observation (with auth)..."
	@./wctl get runtimes --host localhost:9090
	@echo ""
	@./wctl get runtime java --host localhost:9090

# 원격 테스트 (인증 없이)
test-remote-noauth:
	@echo "Testing remote observation (no auth)..."
	@./wctl get runtimes --host localhost:9090 --api-key ""

# 멀티 서버 비교 테스트
test-compare:
	@echo "Testing multi-server comparison..."
	@./wctl compare runtimes --hosts localhost:9090,localhost:9091

# 멀티 서버 비교 - JSON 출력
test-compare-json:
	@echo "Testing comparison with JSON output..."
	@./wctl compare runtimes --hosts localhost:9090,localhost:9091 -o json

# JSON 출력 테스트
test-json:
	@echo "Testing JSON output..."
	@./wctl get runtimes -o json

# YAML 출력 테스트
test-yaml:
	@echo "Testing YAML output..."
	@./wctl get runtimes -o yaml

# 전체 인증 플로우 테스트
test-auth: build test-auth-setup
	@echo ""
	@echo "Starting server in background..."
	@./watcher-server run > /dev/null 2>&1 & echo $$! > /tmp/watcher_server.pid
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
	@echo "  make build              - Build wctl and watcher-server binaries"
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
	@echo "  make test-key-gen       - Generate a new API key"
	@echo "  make test-key-list      - List all client keys"
	@echo "  make test-key-show      - Show default key"
	@echo ""
	@echo "Key Management (Server):"
	@echo "  make test-server-key-list - List all server keys"
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