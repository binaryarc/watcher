# Watcher 👁️

> *"I observe all that transpires here, but I do not interfere."*

A kubectl-style CLI tool for observing runtime versions and system information across your infrastructure. Built for DevOps and SRE professionals.

## Overview

Watcher (`wctl`) is a lightweight system observation tool inspired by Marvel's Watchers - cosmic beings who observe all events but never interfere. Like them, this tool operates in read-only mode, never modifying your systems.

**Core Philosophy**: Observe everything, interfere with nothing.

## Features

- 🔍 **Runtime Detection**: Automatically detect installed runtime versions
- 🌐 **Remote Observation**: Query runtime information from remote servers via gRPC
- 🔄 **Multi-Server Comparison**: Compare runtime versions across multiple servers  # 👈 추가
- 📊 **Multiple Output Formats**: Table (default), JSON, and YAML
- 🚀 **Fast & Lightweight**: Single binary, no dependencies
- 🎯 **kubectl-style UX**: Familiar commands for DevOps users
- 🔒 **Read-Only**: Never modifies your system
- 🐧 **Cross-Platform**: Works on Linux, macOS, and Windows

## Supported Runtimes

### Currently Supported (Tier 1)

| Runtime | Detection Method | Example Version |
|---------|------------------|-----------------|
| Java | `java -version` | 17.0.16, 11.0.19, 8.x |
| Python | `python3 --version` | 3.10.12, 3.9.16 |
| Node.js | `node --version` | 20.18.0, 18.16.0 |
| Go | `go version` | 1.21.5, 1.22.0 |
| Docker | `docker --version` | 24.0.5, 27.5.1 |
| MySQL/MariaDB | `mysql --version` | 8.0.34, 10.11.4 |
| Redis | `redis-server --version` | 7.0.12 |
| Nginx | `nginx -v` | 1.24.0 |

## Installation

### Prerequisites

- Go 1.21 or higher
- Protocol Buffers compiler (for development only)

### From Source
```bash
# Clone the repository
git clone https://github.com/binaryarc/watcher.git
cd watcher

# Build client and server
make build

# Or build individually
go build -o wctl ./cmd/wctl
go build -o watcher-server ./cmd/watcher-server
```

### Using Go Install
```bash
# Install client
go install github.com/binaryarc/watcher/cmd/wctl@latest

# Install server
go install github.com/binaryarc/watcher/cmd/watcher-server@latest
```

## Quick Start

### Local Observation
```bash
# Get all runtimes on local machine
wctl get runtimes

# Get specific runtime
wctl get runtime java

# Different output formats
wctl get runtimes -o json
wctl get runtimes -o yaml
```

**Example Output:**
```
👁️  Observing all runtimes...

┌─────────┬─────────┬──────────────────────────────────────────────────┐
│ RUNTIME │ VERSION │                       PATH                       │
├─────────┼─────────┼──────────────────────────────────────────────────┤
│ java    │ 17.0.16 │ /usr/lib/jvm/java-17-openjdk-amd64/bin/java      │
│ python  │ 3.10.12 │ /usr/bin/python3                                 │
│ node    │ 20.18.0 │ /home/user/.nvm/versions/node/v20.18.0/bin/node  │
│ go      │ 1.21.5  │ /usr/local/go/bin/go                             │
│ docker  │ 24.0.5  │ /usr/bin/docker                                  │
└─────────┴─────────┴──────────────────────────────────────────────────┘

📊 Total: 5 runtime(s) detected
```

### Remote Observation

#### Step 1: Start Watcher Server
```bash
# On the remote server (or locally for testing)
watcher-server serve

# Custom port
watcher-server serve --port 8080

# Custom host and port
watcher-server serve --host 0.0.0.0 --port 9090
```

**Server Output:**
```
🚀 Watcher server listening on 0.0.0.0:9090
📡 Ready to accept observation requests...
```

#### Step 2: Query Remote Server
```bash
# Get all runtimes from remote server
wctl get runtimes --host server.example.com:9090

# Get specific runtime from remote server
wctl get runtime java --host 192.168.1.100:9090

# JSON output from remote
wctl get runtimes --host localhost:9090 -o json
```

**Example Remote Output:**
```
🌐 Connecting to remote server: server.example.com:9090...

┌─────────┬─────────┬──────────────────────────────────────────────────┐
│ RUNTIME │ VERSION │                       PATH                       │
├─────────┼─────────┼──────────────────────────────────────────────────┤
│ java    │ 11.0.19 │ /usr/bin/java                                    │
│ python  │ 3.9.16  │ /usr/bin/python3                                 │
│ docker  │ 24.0.7  │ /usr/bin/docker                                  │
└─────────┴─────────┴──────────────────────────────────────────────────┘

📊 Total: 3 runtime(s) detected
```

### Multi-Server Comparison

Compare runtime versions across multiple servers to identify inconsistencies:
```bash
# Compare runtimes across multiple servers
wctl compare runtimes --hosts server1:9090,server2:9090,server3:9090

# JSON output
wctl compare runtimes --hosts localhost:9090,prod-server:9090 -o json

# YAML output
wctl compare runtimes --hosts server1:9090,server2:9090 -o yaml
```

**Example Output:**
```
🌐 Comparing runtimes across 3 server(s)...

┌─────────┬───────────┬───────────┬───────────┬─────────────┐
│ RUNTIME │ SERVER-1  │ SERVER-2  │ SERVER-3  │   STATUS    │
├─────────┼───────────┼───────────┼───────────┼─────────────┤
│ java    │ 11.0.19   │ 11.0.19   │ 17.0.8    │ ⚠️  DIFF     │
│ python  │ 3.9.16    │ 3.9.16    │ 3.9.16    │ ✅ SAME     │
│ docker  │ 24.0.5    │ 24.0.7    │ -         │ ⚠️  PARTIAL │
└─────────┴───────────┴───────────┴───────────┴─────────────┘

📊 Summary:
   • 3 server(s) compared
   • 2 runtime(s) with differences
   • 1 runtime(s) partially installed
   • 0 runtime(s) consistent
```

## Usage

### Command Structure
```
wctl [command] [subcommand] [flags]
watcher-server [command] [flags]
```

### Client Commands (wctl)

#### Get Commands
```bash
# Local observation
wctl get runtimes                    # Get all runtimes
wctl get runtime <name>              # Get specific runtime

# Remote observation
wctl get runtimes --host <address>   # Get all from remote server
wctl get runtime <name> --host <address>  # Get specific from remote

# Output formats
wctl get runtimes -o table           # Table format (default)
wctl get runtimes -o json            # JSON format
wctl get runtimes -o yaml            # YAML format
```

#### Help Commands
```bash
wctl --help                 # Show help
wctl get --help             # Show get command help
```

### Server Commands (watcher-server)
```bash
# Start server
watcher-server serve                    # Start on default :9090
watcher-server serve --port 8080        # Custom port
watcher-server serve --host 0.0.0.0     # Custom host

# Help
watcher-server --help
watcher-server serve --help
```

#### Compare Commands
```bash
# Compare runtimes across multiple servers
wctl compare runtimes --hosts ,,

# With output formats
wctl compare runtimes --hosts server1:9090,server2:9090 -o json
wctl compare runtimes --hosts server1:9090,server2:9090 -o yaml
```

### Supported Runtime Names

- `java` - Java/OpenJDK
- `python` - Python 2/3
- `node`, `nodejs` - Node.js
- `go`, `golang` - Go language
- `docker` - Docker Engine
- `mysql`, `mariadb` - MySQL/MariaDB
- `redis` - Redis
- `nginx` - Nginx web server

### Real-World Examples
```bash
# Check Java version across multiple servers
wctl get runtime java --host prod-web-01:9090
wctl get runtime java --host prod-web-02:9090
wctl get runtime java --host prod-api-01:9090

# Export inventory to JSON
wctl get runtimes --host prod-web-01:9090 -o json > prod-web-01.json

# Quick health check script
for server in web-{01..05}; do
  echo "=== $server ==="
  wctl get runtimes --host $server.prod.local:9090 -o table
done

# Compare local vs production
diff <(wctl get runtimes -o json) \
     <(wctl get runtimes --host prod-web-01:9090 -o json)
```

## Architecture
```
┌─────────────┐                ┌──────────────────┐
│   wctl      │                │ watcher-server   │
│  (Client)   │────gRPC────────│   (Server)       │
│             │    :9090       │                  │
└─────────────┘                └──────────────────┘
      │                                 │
      │ Local Detection                 │ Local Detection
      │                                 │
      ▼                                 ▼
┌─────────────┐                ┌──────────────────┐
│   Local     │                │   Remote         │
│   System    │                │   System         │
└─────────────┘                └──────────────────┘
```

## Project Structure
```
watcher/
├── cmd/
│   ├── wctl/                  # CLI client entry point
│   └── watcher-server/        # gRPC server entry point
├── pkg/
│   └── cmd/                   # Command implementations
│       ├── root.go            # Root command
│       ├── get/               # Get command group
│       └── serve/             # Server command group
├── internal/
│   ├── detector/              # Runtime detection logic
│   │   ├── detector.go        # Detector interface
│   │   ├── java.go
│   │   ├── python.go
│   │   ├── nodejs.go
│   │   ├── golang.go
│   │   ├── docker.go
│   │   ├── mysql.go
│   │   ├── redis.go
│   │   ├── nginx.go
│   │   └── registry.go        # Detector registry
│   ├── output/                # Output formatters
│   │   ├── table.go
│   │   ├── json.go
│   │   └── yaml.go
│   ├── grpcserver/            # gRPC server implementation
│   │   └── server.go
│   └── grpcclient/            # gRPC client wrapper
│       └── client.go
├── proto/
│   ├── watcher.proto          # Protocol Buffers definition
│   ├── watcher.pb.go          # Generated code
│   └── watcher_grpc.pb.go     # Generated gRPC code
├── Makefile
└── README.md
```

## Development

### Building from Source
```bash
# Clone the repository
git clone https://github.com/binaryarc/watcher.git
cd watcher

# Install dependencies
go mod download

# Build both binaries
make build

# Build individually
go build -o wctl ./cmd/wctl
go build -o watcher-server ./cmd/watcher-server
```

### Development Commands
```bash
# Build
make build              # Build both binaries
make build-verbose      # Build with verbose output

# Generate proto files (after modifying .proto)
make proto

# Testing
make test-local         # Test local observation
make run-server         # Start server (terminal 1)
make test-remote        # Test remote observation (terminal 2)

# Clean
make clean              # Remove built binaries

# Help
make help               # Show all available commands
```

### Testing Workflow

**Terminal 1 (Server):**
```bash
make run-server
```

**Terminal 2 (Client):**
```bash
# Test local
make test-local

# Test remote
make test-remote

# Or manual testing
./wctl get runtimes --host localhost:9090
./wctl get runtime java --host localhost:9090 -o json
```

### Adding New Runtime Detector

1. Create new detector file in `internal/detector/`:
```go
package detector

type NewRuntimeDetector struct{}

func (d *NewRuntimeDetector) Name() string {
    return "newruntime"
}

func (d *NewRuntimeDetector) Detect() (*Runtime, error) {
    // Implementation
}
```

2. Add to `registry.go`:
```go
func GetAllDetectors() []Detector {
    return []Detector{
        // ... existing detectors
        &NewRuntimeDetector{},
    }
}
```

3. Update `runtime.go` switch case

## Protocol

### gRPC Service Definition
```protobuf
service WatcherService {
  rpc ObserveRuntimes(ObserveRequest) returns (ObserveResponse);
}

message Runtime {
  string name = 1;
  string version = 2;
  string path = 3;
  bool found = 4;
}

message ObserveResponse {
  repeated Runtime runtimes = 1;
  SystemInfo system_info = 2;
  int64 timestamp = 3;
}
```

## Roadmap

### Phase 1: Local Detection ✅ (Completed)

- [x] Basic CLI structure
- [x] Runtime detection (Java, Python, Node.js, Go, Docker, MySQL, Redis, Nginx)
- [x] Multiple output formats (table, json, yaml)

### Phase 2: Remote Observation ✅ (Completed)

- [x] gRPC protocol definition
- [x] Server implementation (watcher-server)
- [x] Client implementation with --host flag
- [x] Remote runtime detection

### Phase 3: Multi-Server Comparison ✅ (Completed)

- [x] `wctl compare runtimes --hosts server1:9090,server2:9090`
- [x] Side-by-side version comparison
- [x] Detect version mismatches across infrastructure

### Phase 4: Security & Production Ready

- [ ] TLS/mTLS support
- [ ] Authentication and authorization
- [ ] Connection pooling and retry logic
- [ ] Rate limiting

### Future Enhancements

- [ ] Service detection (systemd, docker containers)
- [ ] Version history tracking
- [ ] Security vulnerability detection
- [ ] Configuration file support (~/.watcher/config.yaml)
- [ ] Web UI dashboard
- [ ] Prometheus metrics export

## Why "Watcher"?

Inspired by Marvel's Watchers - cosmic beings who observe all events across the multiverse but never interfere. This philosophy perfectly matches what a monitoring tool should do:

- **Observe everything**: Monitor all your servers and runtimes
- **Never interfere**: Read-only operations, zero system modifications
- **Multiverse aware**: Designed for multi-server environments

## Use Cases

### Version Consistency Check
```bash
# Compare versions across production fleet
wctl compare runtimes --hosts \
  prod-web-01:9090,\
  prod-web-02:9090,\
  prod-web-03:9090 \
  -o table
```

### Infrastructure Audit
```bash
# Quick inventory of all servers
for server in prod-{01..10}; do
  wctl get runtimes --host $server:9090 -o json >> inventory.jsonl
done
```

### Pre-Deployment Verification
```bash
# Ensure Java 17 before deploying Spring Boot 3.x
wctl get runtime java --host prod-web-01:9090
```

### Compliance Checking
```bash
# Find servers running outdated Python versions
wctl get runtime python --host server:9090 -o json | \
  jq 'select(.Version | startswith("2."))'
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Contribution Areas

- New runtime detectors
- Bug fixes and improvements
- Documentation enhancements
- Test coverage
- Performance optimizations

## Security

- All communication uses gRPC (TLS support coming in Phase 4)
- Read-only operations only
- No sensitive data collection
- No system modifications

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Built with ❤️ for DevOps and SRE professionals.

## Acknowledgments

- Inspired by kubectl's intuitive command structure
- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
- Table output powered by [tablewriter](https://github.com/olekukonko/tablewriter)
- gRPC communication via [grpc-go](https://github.com/grpc/grpc-go)

---

*"I am the Watcher. I am your guide through these vast new realities."*