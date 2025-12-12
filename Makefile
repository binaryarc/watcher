.PHONY: proto build clean run-server test-local test-remote test-compare test-compare-json test-json test-yaml help

# 기본 타겟
all: build

# Proto 파일 생성
proto:
	@echo "🔄 Generating proto files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/watcher.proto
	@echo "✅ Proto files generated!"

# 빌드
build:
	@echo "🔨 Building binaries..."
	@go build -o wctl ./cmd/wctl
	@go build -o watcher-server ./cmd/watcher-server
	@echo "✅ Build complete!"
	@echo "   📦 wctl binary created"
	@echo "   📦 watcher-server binary created"

# 빌드 (verbose)
build-verbose:
	@echo "🔨 Building binaries (verbose)..."
	go build -v -o wctl ./cmd/wctl
	go build -v -o watcher-server ./cmd/watcher-server
	@echo "✅ Build complete!"

# 클린
clean:
	@echo "🧹 Cleaning up..."
	@rm -f wctl watcher-server
	@echo "✅ Clean complete!"

# 서버 실행
run-server:
	@echo "🚀 Starting Watcher server..."
	./watcher-server serve

# 서버 실행 (커스텀 포트)
run-server-custom:
	@echo "🚀 Starting Watcher server on port 8080..."
	./watcher-server serve --port 8080

# 로컬 테스트
test-local:
	@echo "👁️  Testing local observation..."
	@./wctl get runtimes
	@echo ""
	@./wctl get runtime java

# 원격 테스트 (서버가 실행중이어야 함)
test-remote:
	@echo "🌐 Testing remote observation..."
	@./wctl get runtimes --host localhost:9090
	@echo ""
	@./wctl get runtime java --host localhost:9090

# 멀티 서버 비교 테스트 (서버들이 실행중이어야 함)
test-compare:
	@echo "🔍 Testing multi-server comparison..."
	@./wctl compare runtimes --hosts localhost:9090,localhost:9091

# 멀티 서버 비교 - JSON 출력
test-compare-json:
	@echo "🔍 Testing comparison with JSON output..."
	@./wctl compare runtimes --hosts localhost:9090,localhost:9091 -o json

# JSON 출력 테스트
test-json:
	@echo "📄 Testing JSON output..."
	@./wctl get runtimes -o json

# YAML 출력 테스트
test-yaml:
	@echo "📄 Testing YAML output..."
	@./wctl get runtimes -o yaml

# 도움말
help:
	@echo "Watcher Makefile Commands:"
	@echo ""
	@echo "  make build           - Build wctl and watcher-server binaries"
	@echo "  make proto           - Generate proto files"
	@echo "  make clean           - Remove built binaries"
	@echo "  make run-server      - Start watcher server on :9090"
	@echo "  make test-local      - Test local runtime observation"
	@echo "  make test-remote     - Test remote runtime observation (needs server)"
	@echo "  make test-compare    - Test multi-server comparison (needs servers)"  # 👈 추가
	@echo "  make test-json       - Test JSON output format"
	@echo "  make test-yaml       - Test YAML output format"
	@echo "  make help            - Show this help message"