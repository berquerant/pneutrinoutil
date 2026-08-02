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

## 2. Development Workflow & Verification

AI assistants **MUST** follow a strict **"Edit → Verify (Lint/Test/Build)"** loop. Editing files alone is not considered task completion.

### Task Runner Commands (`Taskfile`)

Use the `./task` runner for common development commands:

| Task | Command | Description |
| :--- | :--- | :--- |
| **Run Unit Tests** | `./task test:unit` | Runs all Go unit tests. |
| **Lint Code** | `./task lint` | Runs `golangci-lint` and `eslint`. |
| **Build Binaries** | `./task build` | Builds all Go binaries (`cli`, `server`, `worker`, etc.). |
| **Tidy Modules** | `./task tidy` | Runs `go mod tidy` for root and tool modules. |
| **Full Verification** | `./task` | Runs linting, tests, and builds in sequence. |

---

## 3. Component-Specific Guidelines

### A. Go Backend (`server/`, `worker/`, `pkg/`)
1. **Maintain Layered Architecture:**
   * Business logic and processing pipelines belong in [`pkg/domain`](pkg/domain) and [`pkg/task`](pkg/task).
   * Database queries and object storage operations belong in [`pkg/repo`](pkg/repo) and [`pkg/infra`](pkg/infra).
2. **API Changes & Swagger:**
   * When modifying handlers in [`server/handler`](server/handler), update the Swagger annotations accordingly and verify client compatibility in `ui/`.
3. **Error Handling:**
   * Do not suppress errors or return silent fallbacks. Log errors with sufficient context and propagate them up the call chain.

### B. CLI Tool (`cli/`)
* When adding or updating CLI flags/parameters, update [`cli/README.md`](cli/README.md) with updated flag descriptions and usage examples.
* Ensure backward compatibility for NEUTRINO v2 and v3 parameters.

### C. Web UI (`ui/`)
1. **Package Management:** Always use `pnpm`.
2. **Type Safety:** Ensure strong TypeScript typing for API request/response structures. Avoid using `any`.
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

- [ ] `./task lint` passes without errors or warnings.
- [ ] `./task test:unit` passes completely.
- [ ] `./task build` succeeds for all target binaries.
- [ ] New or modified code is accompanied by unit tests where appropriate.
- [ ] Documentation (`README.md`, Swagger annotations, etc.) has been updated.
