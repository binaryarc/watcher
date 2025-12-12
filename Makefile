.PHONY: proto

proto:
	@echo "🔄 Generating proto files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/watcher.proto
	@echo "✅ Proto files generated!"