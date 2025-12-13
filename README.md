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
make build        # builds wctl (client) and wsctl (server)
```

Or:
```bash
go install github.com/binaryarc/watcher/cmd/wctl@latest
go install github.com/binaryarc/watcher/cmd/wsctl@latest
```

## Usage

Local check:
```bash
wctl get runtimes
wctl get runtime java
```

Remote check (start the server on the target with `wsctl run`):
```bash
wctl get runtimes --host 192.168.1.100:9090
wctl get runtime java --host server.example.com:9090 -o json
```

Compare across servers:
```bash
wctl compare runtimes --hosts server1:9090,server2:9090,server3:9090
```

Client-side key management:
```bash
wctl key generate   # create & store a new API key under ~/.watcher/keys
wctl get key        # print the currently saved key
```

## Example output
```
RUNTIME  VERSION   PATH
java     17.0.16   /usr/lib/jvm/java-17-openjdk/bin/java
python   3.10.12   /usr/bin/python3
docker   24.0.5    /usr/bin/docker
```

Multi-server comparison (table output):
```
┌─────────┬───────────┬───────────┬─────────┐
│ RUNTIME │ SERVER-1  │ SERVER-2  │ STATUS  │
├─────────┼───────────┼───────────┼─────────┤
│ python  │ 3.10.12   │ 3.10.12   │ SAME    │
│ go      │ 1.25.5    │ 1.25.5    │ SAME    │
│ docker  │ x         │ 27.5.1    │ PARTIAL │
│ java    │ 17.0.17   │ x         │ PARTIAL │
│ node    │ 20.18.0   │ x         │ PARTIAL │
└─────────┴───────────┴───────────┴─────────┘
```

## How it works
```
wctl (client) ---gRPC:9090---> wsctl (server)
                                    |
                             Detects runtimes
                             (reads version commands)
```

## Authentication & API keys

Watcher uses shared API keys for every remote RPC call.

- Generate a client key with `wctl key generate`. The CLI saves it at `~/.watcher/keys/default` and automatically loads it in this order: `--api-key` flag > `WATCHER_API_KEY` env var > saved file.
- Start the server and register the client key:
  ```bash
  wsctl add key <api-key> "CI pipeline token"
  wsctl run --port 9090        # uses ~/.watcher/server/keys.json by default
  ```
- Inspect or clean keys on the server whenever needed:
  ```bash
  wsctl get keys
  wsctl delete key --name <api-key>
  wsctl clear keys             # remove all keys (asks for confirmation)
  ```
- For quick lab testing you can start `wsctl run --disable-auth`, but production deployments should always have keys registered.

Once a key is registered on both sides, any `wctl` command will include it automatically. You can still override it per command:
```bash
WATCHER_API_KEY=$(wctl get key)
wctl get runtimes --host host.example.com --api-key "$WATCHER_API_KEY"
```

## wsctl (server CLI)

`wsctl` controls the Watcher server process:

- `wsctl run [--port 9090 --host 0.0.0.0 --keystore /path/to/keys.json]` starts the gRPC server and enforces API keys unless `--disable-auth` is passed.
- `wsctl add key <api-key> "<description>"` registers a client key so remote calls will be accepted.
- `wsctl get keys` lists every registered key with timestamps; use `wsctl delete key --name <api-key>` or `wsctl clear keys` to revoke them.
- By default the server stores its key database under `~/.watcher/server/keys.json`; create a shared location with `--keystore` when running inside containers/VMs.

## Supported runtimes

Java, Python, Node.js, Go, Docker, MySQL/MariaDB, Redis, Nginx

Detection via version commands: `java -version`, `python3 --version`, etc.

## Project structure
```
cmd/
  wctl/           - CLI client
  wsctl/          - gRPC server CLI
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

- [ ] Better error handling for network failures
- [ ] Tests
- [ ] Config file support
- [ ] Command auto-completion (bash/zsh/fish)

## Why "watcher"?

Read-only observation tool. Doesn't modify anything on target systems.

## License

MIT
