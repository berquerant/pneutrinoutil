# Gemini Code Understanding

This document provides a comprehensive overview of the `pneutrinoutil` project, designed to be used as a context for AI-powered development with Gemini.

## Project Overview

`pneutrinoutil` is a set of utilities for [NEUTRINO](https://studio-neutrino.com/), a singing voice synthesizer. It provides a command-line interface (CLI) for processing MusicXML files and a web-based user interface for interacting with the NEUTRINO engine.

The project is a monorepo containing several components:

*   **`server`**: A Go-based backend server that exposes a REST API for synthesizing singing voices.
*   **`ui`**: A React-based frontend that provides a user interface for the server.
*   **`cli`**: A Go-based command-line interface for batch processing of MusicXML files.
*   **`worker`**: A Go-based worker that processes background jobs.
*   **`pkg`**: A collection of shared Go packages used by the other components.

## Architecture

The project follows a microservices-like architecture, with separate components for the backend server, frontend UI, and background worker. These components communicate with each other through a REST API and a Redis-based message broker.

The infrastructure is managed using Docker and Docker Compose. It consists of the following services:

*   **`server`**: The Go backend server.
*   **`ui`**: The React frontend.
*   **`mysql`**: A MySQL database for storing application data.
*   **`redis`**: A Redis instance for caching and message broking.
*   **`minio`**: An S3-compatible object storage service for storing generated audio files.

## Building and Running

The project uses a `Taskfile.yml` to define common development tasks. The following are some of the most important commands:

*   **`./task build`**: Build all the project's binaries.
*   **`./task start`**: Start all the services (server, UI, and worker).
*   **`./task stop`**: Stop all the services.
*   **`./task test`**: Run the unit tests.
*   **`./task e2e`**: Run the end-to-end tests.
*   **`./task lint`**: Run the linters.

The web UI is available at `http://localhost:9201/` and the Swagger API documentation is at `http://localhost:9101/v1/swagger/index.html`.

## Development Conventions

*   **Go**: The Go code follows standard Go conventions.
*   **React**: The React code is written in TypeScript and uses `pnpm` for package management.
*   **Linting**: The project uses `golangci-lint` for Go and `eslint` for TypeScript to enforce code style and catch potential errors.
*   **Testing**: The project has a suite of unit tests and end-to-end tests. New code should be accompanied by tests.
