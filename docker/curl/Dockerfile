# syntax=docker/dockerfile:1
# check=error=true
FROM debian:trixie-slim

RUN --mount=type=cache,target=/var/lib/apt,sharing=locked \
    --mount=type=cache,target=/var/cache/apt,sharing=locked \
    apt-get update && apt-get install -y --no-install-recommends curl ca-certificates
