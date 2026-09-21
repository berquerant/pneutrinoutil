#!/bin/bash

set -e
set -o pipefail

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly bind="${d}/../bin/tools"
mkdir -p "$bind"

log() {
    echo >&2 "$(basename "$0"): $*"
}

readonly name="$1"
if [[ -z "$name" ]] ; then
    log "name(\$1) required"
    exit 1
fi
shift

# In non-CI environments, extract the GitHub token from the local 'gh' CLI if
# available and valid, so that 'mise' can make authenticated GitHub API calls to
# avoid unauthenticated rate limiting (403 Forbidden).
if [[ -z "${GITHUB_TOKEN:-}" && -z "${GH_TOKEN:-}" && "${GITHUB_ACTIONS:-}" != "true" ]] && command -v gh >/dev/null 2>&1; then
    if gh auth status >/dev/null 2>&1; then
        _gh_token="$(gh auth token 2>/dev/null || true)"
        if [[ -n "$_gh_token" ]]; then
            export GITHUB_TOKEN="$_gh_token"
        fi
    fi
fi

if command -v mise >/dev/null 2>&1 ; then
    exec mise exec -- "$name" "$@"
elif [ -x "${HOME}/.local/bin/mise" ]; then
    exec "${HOME}/.local/bin/mise" exec -- "$name" "$@"
fi

exec "$name" "$@"
