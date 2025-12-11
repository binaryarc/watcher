# Watcher 👁️

> *"I observe all that transpires here, but I do not interfere."*

A kubectl-style CLI tool for observing runtime versions and system information across your infrastructure. Built for DevOps and SRE professionals.

## Overview

Watcher (`wctl`) is a lightweight system observation tool inspired by Marvel's Watchers - cosmic beings who observe all events but never interfere. Like them, this tool operates in read-only mode, never modifying your systems.

**Core Philosophy**: Observe everything, interfere with nothing.

## Features

- 🔍 **Runtime Detection**: Automatically detect installed runtime versions
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

### From Source
```bash
# Clone the repository
git clone https://github.com/binaryarc/watcher.git
cd watcher

# Build
go build -o wctl ./cmd/wctl

# Optional: Install to GOPATH
go install ./cmd/wctl
```

### Using Go Install
```bash
go install github.com/binaryarc/watcher/cmd/wctl@latest
```

## Quick Start

### Get All Runtimes
```bash
# Table format (default)
wctl get runtimes

# JSON format
wctl get runtimes -o json

# YAML format
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

### Get Specific Runtime
```bash
# Query specific runtime
wctl get runtime java
wctl get runtime python
wctl get runtime docker

# With different output formats
wctl get runtime go -o json
wctl get runtime node -o yaml
```

**Example Output:**
```
👁️  Observing java runtime...

✅ java detected!

┌──────────┬─────────────────────────────────────────────┐
│ PROPERTY │                    VALUE                    │
├──────────┼─────────────────────────────────────────────┤
│ Name     │ java                                        │
│ Version  │ 17.0.16                                     │
│ Path     │ /usr/lib/jvm/java-17-openjdk-amd64/bin/java │
└──────────┴─────────────────────────────────────────────┘
```

## Usage

### Command Structure
```
wctl [command] [subcommand] [flags]
```

### Available Commands
```bash
wctl get runtimes           # Get all detected runtimes
wctl get runtime <name>     # Get specific runtime
wctl --help                 # Show help
wctl get --help             # Show get command help
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

### Output Formats

Use the `-o` or `--output` flag:
```bash
-o table    # ASCII table (default)
-o json     # JSON format
-o yaml     # YAML format
```

### Examples
```bash
# Check if Java is installed
wctl get runtime java

# Get all runtimes in JSON (useful for scripts)
wctl get runtimes -o json | jq '.[] | select(.Name=="docker")'

# Export to YAML file
wctl get runtimes -o yaml > runtimes.yaml

# Check multiple specific runtimes
wctl get runtime java
wctl get runtime python
wctl get runtime docker
```

## Project Structure
```
watcher/
├── cmd/
│   └── wctl/              # CLI entry point
├── pkg/
│   └── cmd/               # Command implementations
│       ├── root.go        # Root command
│       └── get/           # Get command group
├── internal/
│   ├── detector/          # Runtime detection logic
│   │   ├── java.go
│   │   ├── python.go
│   │   ├── nodejs.go
│   │   ├── golang.go
│   │   ├── docker.go
│   │   ├── mysql.go
│   │   ├── redis.go
│   │   └── nginx.go
│   └── output/            # Output formatters
│       ├── table.go
│       ├── json.go
│       └── yaml.go
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

# Build
go build -o wctl ./cmd/wctl

# Run tests
go test ./...
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

## Roadmap

### Current Phase: MVP - Local Detection ✅

- [x] Basic CLI structure
- [x] Runtime detection (Java, Python, Node.js, Go, Docker, MySQL, Redis, Nginx)
- [x] Multiple output formats (table, json, yaml)

### Next Phase: Remote Observation

- [ ] gRPC protocol definition
- [ ] Server implementation (watcher-server)
- [ ] Remote runtime detection
- [ ] Multi-server comparison
- [ ] TLS/mTLS support

### Future Enhancements

- [ ] Service detection (systemd, docker containers)
- [ ] Version history tracking
- [ ] Security vulnerability detection
- [ ] Configuration file support
- [ ] Web UI dashboard

## Why "Watcher"?

Inspired by Marvel's Watchers - cosmic beings who observe all events across the multiverse but never interfere. This philosophy perfectly matches what a monitoring tool should do:

- **Observe everything**: Monitor all your servers and runtimes
- **Never interfere**: Read-only operations, zero system modifications
- **Multiverse aware**: Designed for multi-server environments (coming soon)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Author

Built with ❤️ for DevOps and SRE professionals.

## Acknowledgments

- Inspired by kubectl's intuitive command structure
- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
- Table output powered by [tablewriter](https://github.com/olekukonko/tablewriter)

---

*"I am the Watcher. I am your guide through these vast new realities."*