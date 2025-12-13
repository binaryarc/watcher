# watcher 👀

**Compare runtime environments across servers — fast, read-only, and CI-friendly.**

Watcher is a lightweight CLI that detects and compares runtimes (Java, Python, Node.js, Go, Docker, etc.) across local and remote servers over gRPC. It gives operators and SREs reliable answers to:

> “Why does this server behave differently from that one?”  
> “Are all environments actually running the same versions?”

⭐ If this tool helps you inspect or debug infrastructure, consider giving it a star — it really helps.

---

## Why Watcher?

- Read-only, safe to run on production boxes
- Fast comparisons across many hosts
- Extensible detector registry for new runtimes
- Friendly to CI pipelines and automation
- API-key authentication for remote access

Watcher observes; it never mutates target machines. Run it on prod nodes, CI runners, or anywhere you need version truth.

---

## What it does

- Detects installed runtimes and versions:
  - Java, Python, Node.js, Go, Docker
  - MySQL/MariaDB, Redis, Nginx
- Collects data locally or remotely via gRPC
- Compares versions across multiple servers
- Outputs results as tables, JSON, or YAML

---

## Example

### Multi-server comparison

```
┌─────────┬───────────┬───────────┬─────────┐
│ RUNTIME │ SERVER-1  │ SERVER-2  │ STATUS  │
├─────────┼───────────┼───────────┼─────────┤
│ python  │ 3.10.12   │ 3.10.12   │ SAME    │
│ go      │ 1.25.5    │ 1.25.5    │ SAME    │
│ node    │ x         │ 20.18.0   │ PARTIAL │
│ docker  │ x         │ 27.5.1    │ PARTIAL │
│ java    │ x         │ 17.0.17   │ PARTIAL │
└─────────┴───────────┴───────────┴─────────┘
```

Environment drift becomes instantly visible.

---

## Installation

### Quick install

```bash
go install github.com/binaryarc/watcher/cmd/wctl@latest
go install github.com/binaryarc/watcher/cmd/wsctl@latest
```

Ensure `$GOPATH/bin` is on your `PATH`.

### Build from source

```bash
git clone https://github.com/binaryarc/watcher.git
cd watcher
make build
```

This produces:

- `wctl` — client CLI
- `wsctl` — server CLI

---

## Usage

### Local runtime check

```bash
wctl get runtimes
wctl get runtime java
```

### Remote runtime check

On each target server:

```bash
wsctl run --port 9090
```

From your client:

```bash
wctl get runtimes --host 192.168.1.100:9090
wctl get runtime java --host server.example.com:9090 -o json
```

### Compare multiple servers

```bash
wctl compare runtimes --hosts server1:9090,server2:9090,server3:9090
```

---

## Authentication & API keys

Watcher uses shared API keys for every remote RPC.

### Client side

Generate and store a key:

```bash
wctl key generate
```

Print the currently saved key:

```bash
wctl get key
```

Keys load automatically in this order:

1. `--api-key` flag
2. `WATCHER_API_KEY` environment variable
3. File at `~/.watcher/keys/default`

### Server side

Register keys and start the server:

```bash
wsctl add key <api-key> "CI pipeline token"
wsctl run --port 9090
```

Manage keys:

```bash
wsctl get keys
wsctl delete key --name <api-key>
wsctl clear keys
```

For quick tests you can disable auth:

```bash
wsctl run --disable-auth
```

(Not recommended for production.)

---

## How it works

```
wctl (client) --- gRPC ---> wsctl (server)
                           |
                     runtime detectors
```

---

## Supported runtimes

- Java
- Python
- Node.js
- Go
- Docker
- MySQL / MariaDB
- Redis
- Nginx

Detection relies on standard version commands (`java -version`, `python3 --version`, ...).

---

## Project structure

```
cmd/
  wctl/           CLI client
  wsctl/          gRPC server CLI
internal/
  detector/       runtime detection logic
  grpcclient/     client wrapper
  grpcserver/     server implementation
proto/            gRPC definitions
```

---

## Roadmap

- Binary releases (Linux/macOS) — first tag: `v0.1.0`
- One-line install script
- Shell auto-completion (bash/zsh/fish)
- Baseline comparison (`--baseline`)
- CI mode with exit codes (`--ci`)
- Snapshot & diff support

Ideas, bug reports, and feature requests are welcome.

---

## Contributing

Watcher is evolving. Issues and pull requests are encouraged.

If this fits your workflow, a ⭐ on the repository makes a big difference.

---

## License

Watcher is distributed under the MIT License. See `LICENSE` for details.
