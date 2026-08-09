# pneutrinoutil

Utilities for [NEUTRINO](https://studio-neutrino.com/) AI singing voice synthesizer.

Provides a Command Line Interface (CLI), REST API server, background job worker, and Web UI.

---

## Requirements

- macOS (or Linux)
- [mise](https://github.com/jdx/mise) (polyglot dev tool & environment manager)
- Docker
- AWS CLI

*(Or use [lima](https://github.com/lima-vm/lima) environment)*

---

## Setup & Provisioning

### 1. Initial Environment & Tools Setup
Install all required runtime versions (Go, Node.js, Python, uv, pnpm) and development CLI tools (`golangci-lint`, `kubectl`, `helm`, `kind`, `task`, `yq` etc.) defined in [`mise.toml`](./mise.toml):
```shell
mise install
```

### 2. Download NEUTRINO & Singer Voice Models
Download and install NEUTRINO binaries and singer libraries via Ansible:
```shell
./task ansible
```

---

## Usage

### CLI

Batch generate `.wav` audio from a `.musicxml` file:
```shell
./task build:cli
./dist/pneutrinoutil --score /path/to/some.musicxml
```

### HTTP Server & Web UI (Local Kubernetes Deployment)

#### Start Services
Deploy local Kind Kubernetes cluster with Server, UI, Worker, MySQL, Redis, and MinIO:
```shell
./task k8s
```

Once running, access:
- **Web UI (Kind):** [http://localhost:3000/](http://localhost:3000/)
- **Swagger API Docs:** [http://localhost:9101/v1/swagger/index.html](http://localhost:9101/v1/swagger/index.html)

#### Stop Services
Stop and tear down local Kind cluster and background worker:
```shell
./task k8s:stop
```

---

## Development & Component-specific Execution

### Run Web UI Dev Server
```shell
./task ui-dev
```

### Run API Server Directly
```shell
./task build:server
./dist/pneutrinoutil-server
```

### Reload K8s Worker Process
```shell
./task run:reload-k8s-worker
```

### Development with Lima
```shell
./bin/lima.sh start
./bin/lima.sh run ./task test:unit
```
