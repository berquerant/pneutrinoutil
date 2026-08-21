# 🤖 AI Assistant Development Guide (`pneutrinoutil`)

This document provides guidelines and best practices for AI assistants (and pair-programming engineers) developing, maintaining, and adding features to the `pneutrinoutil` project.

---

## 1. Project Overview & Architecture

`pneutrinoutil` is a monorepo providing utilities for [NEUTRINO](https://studio-neutrino.com/), an AI singing voice synthesizer.

### Component Structure

* **`cli/`**: Go-based Command Line Interface for batch rendering `.musicxml` files to `.wav`.
* **`server/`**: Go-based HTTP REST API server (Echo framework) with Swagger/OpenAPI support.
* **`worker/`**: Go-based background task worker processing asynchronous synthesis jobs via Redis (`asynq`).
* **`ui/`**: React (React Router v7 + TypeScript) web frontend managed with `pnpm`.
* **`pkg/`**: Shared Go packages including domain models (`pkg/domain`), infrastructure (`pkg/infra`), repositories (`pkg/repo`), and task execution pipelines (`pkg/task`).

---

## 2. Development Workflow & Task Execution Matrix

AI assistants **MUST** follow a strict **"Edit → Verify (Lint/Test/Build)"** loop. Editing files alone is not considered task completion.

Use the `./task` runner (which executes `mise run` under the hood) for running standardized commands based on the scenario:

### 📋 Scenario-Based Task Matrix

| Situation / Development Trigger | Command | Description |
| :--- | :--- | :--- |
| **Routine Full Verification** | `./task` | Runs linters, unit/integration tests, and builds all binaries in sequence. |
| **Go Code Modifications** | `./task test:unit` | Runs all Go unit tests with coverage (`go test -cover`). |
| **API Handler / Swagger Annotation Changes** | `./task gen:swag` | Generates Swagger docs in `server/docs` and updates TypeScript API client for UI. |
| **Frontend UI (TypeScript/React) Editing** | `./task ui-lint` | Runs TypeScript type checking (`typecheck`). |
| **End-to-End System Testing** | `./task test:e2e` | Runs E2E tests, verifying integration between Server, Worker, Redis, MySQL, and S3. |
| **Building Individual Binaries** | `./task build:cli`<br>`./task build:server`<br>`./task build:worker` | Builds specific Go binaries to `dist/`. |
| **Building All Artifacts & Docker Images** | `./task build` | Builds all Go binaries and Docker images via `docker buildx bake`. |
| **Updating Go Module Dependencies** | `./task tidy` | Runs `go mod tidy` for Go module dependencies. |
| **Initial Project Setup / Environment Config** | `./task init` | Initializes local environment variables via `mise`. |
| **Deploying Local Kubernetes (Kind)** | `./task k8s` | Reloads Kind cluster, loads Docker images, and deploys via Helm. |
| **Stopping Local Kubernetes & Worker** | `./task k8s:stop` | Tears down local Kind cluster and stops background worker processes. |
| **Reloading K8s Worker Process** | `./task run:reload-k8s-worker` | Rebuilds CLI/Worker and restarts background K8s worker process. |
| **Provisioning NEUTRINO Engine & Singers** | `./task ansible` | Downloads and installs NEUTRINO binaries and singer voice models via Ansible. |
| **Running Unit Tests inside Lima VM** | `./task lima:unit` | Executes unit tests inside isolated Lima VM environment (`./bin/lima.sh run ./task test:unit`). |
| **Running E2E Tests inside Lima VM** | `./task lima:e2e` | Executes E2E tests inside isolated Lima VM environment (`./bin/lima.sh run ./task test:e2e`). |
| **Managing Lima VM Lifecycle** | `./task lima:start`<br>`./task lima:stop`<br>`./task lima:reload` | Starts, stops, or recreates Lima VM environment. |
| **Cleaning Generated Files / Tools** | `./task gen:clean` | Removes generated Go files (`*_generated.go`) and binary tool caches. |

---

## 3. Component-Specific Guidelines

### A. Go Backend (`server/`, `worker/`, `pkg/`)
1. **Maintain Layered Architecture:**
   * Business logic and processing pipelines belong in [`pkg/domain`](pkg/domain) and [`pkg/task`](pkg/task).
   * Database queries and object storage operations belong in [`pkg/repo`](pkg/repo) and [`pkg/infra`](pkg/infra).
2. **API Changes & Swagger:**
   * When modifying handlers in [`server/handler`](server/handler), update the Swagger annotations accordingly, run `./task gen:swag` to update Swagger docs and client interfaces, and verify UI compatibility in `ui/`.
3. **Error Handling:**
   * Do not suppress errors or return silent fallbacks. Log errors with sufficient context and propagate them up the call chain.

### B. CLI Tool (`cli/`)
* When adding or updating CLI flags/parameters, update [`cli/README.md`](cli/README.md) with updated flag descriptions and usage examples.
* Ensure backward compatibility for NEUTRINO v2 and v3 parameters.

### C. Web UI (`ui/`)
1. **Package Management:** Always use `pnpm`.
2. **Type Safety:** Ensure strong TypeScript typing for API request/response structures. Avoid using `any`. Run `./task ui-lint` after modifying UI code.
3. **UI/UX Excellence:** Maintain a modern, responsive interface using rich aesthetics (clean color palettes, micro-animations, clear status indicators for job processing and audio playback).

---

## 4. Golden Rules for AI Assistants

1. **Never Guess Schemas or Definitions:** Always inspect authoritative source code for structs, interfaces, or DB schemas before writing code that consumes them.
2. **Log-First Diagnosis:** When encountering runtime or test failures, inspect the raw error logs first rather than forming blind hypotheses or masking errors.
3. **Preserve Contracts:** Ensure API endpoints, database schemas, and background job message payloads retain compatibility across components (`server`, `worker`, `ui`).
4. **Always Run Verification:** Never declare success without executing empirical verification (e.g., `./task test:unit`, `./task lint`, `./task build`).
5. **Keep Documentation Synced:** Update [`GEMINI.md`](GEMINI.md), `README.md`, or this guide whenever architectural changes or new configuration parameters are introduced.

---

## 5. Definition of Done Checklist

Before submitting changes, verify the following:

- [ ] `./task lint` (or `./task ui-lint` for UI changes) passes without errors or warnings.
- [ ] `./task test:unit` passes completely.
- [ ] `./task build` succeeds for all target binaries.
- [ ] `./task gen:swag` was executed if API endpoints/annotations were modified.
- [ ] New or modified code is accompanied by unit tests where appropriate.
- [ ] Documentation (`README.md`, Swagger annotations, etc.) has been updated.
