# watcher

CLI tool for checking runtime versions across multiple servers. Useful when you need to quickly verify what's installed where.

## What it does

- Detects runtime versions on local/remote machines (Java, Python, Node, Go, Docker, MySQL, Redis, Nginx)
- Compare versions across multiple servers via gRPC
- Output in table, JSON, or YAML

## Install
```bash
git clone https://github.com/binaryarc/watcher.git
cd watcher
make build
```

Or:
```bash
go install github.com/binaryarc/watcher/cmd/wctl@latest
go install github.com/binaryarc/watcher/cmd/watcher-server@latest
```

## Usage

Local check:
```bash
wctl get runtimes
wctl get runtime java
```

Remote check (need to run `watcher-server serve` on target):
```bash
wctl get runtimes --host 192.168.1.100:9090
wctl get runtime java --host server.example.com:9090 -o json
```

Compare across servers:
```bash
wctl compare runtimes --hosts server1:9090,server2:9090,server3:9090
```

## Example output
```
RUNTIME  VERSION   PATH
java     17.0.16   /usr/lib/jvm/java-17-openjdk/bin/java
python   3.10.12   /usr/bin/python3
docker   24.0.5    /usr/bin/docker
```

## How it works
```
wctl (client) ---gRPC:9090---> watcher-server
                                     |
                              Detects runtimes
                              (reads version commands)
```

## Supported runtimes

Java, Python, Node.js, Go, Docker, MySQL/MariaDB, Redis, Nginx

Detection via version commands: `java -version`, `python3 --version`, etc.

## Project structure
```
cmd/
  wctl/           - CLI client
  watcher-server/ - gRPC server
internal/
  detector/       - Runtime detection logic
  grpcserver/     - Server implementation
  grpcclient/     - Client wrapper
proto/            - gRPC definitions
```

## Development
```bash
make build        # build both binaries
make proto        # regenerate proto files
make test-local   # test local detection
make run-server   # start server for testing
```

Adding new runtime detector:
1. Create `internal/detector/newruntime.go`
2. Implement `Detector` interface
3. Register in `registry.go`

## TODO

- [ ] TLS/mTLS support
- [ ] Auth
- [ ] Better error handling for network failures
- [ ] Tests
- [ ] Config file support

## Why "watcher"?

Read-only observation tool. Doesn't modify anything on target systems.

## License

MIT