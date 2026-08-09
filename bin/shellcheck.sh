#!/bin/bash

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"

find_by_shebang() {
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        git grep -l "^${1}"
    else
        find . -type d \( -name .git -o -name .venv -o -name node_modules -o -name tmp -o -name dist \) -prune -o -type f -exec grep -l "^${1}" {} + 2>/dev/null || true
    fi
}

find_by_extension() {
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        git ls-files | grep "\.${1}$" || true
    else
        find . -type d \( -name .git -o -name .venv -o -name node_modules -o -name tmp -o -name dist \) -prune -o -name "*.${1}" -print 2>/dev/null || true
    fi
}

find_by_interpreter() {
    local -r interpreter="$1"
    local -r extension="$2"
    {
        find_by_shebang "#\!/bin/${interpreter}"
        find_by_shebang "#\!/usr/bin/env ${interpreter}"
        find_by_extension "$extension"
    } | sort -u
}

ignore() {
    grep -v -E 'ui/app/api/client|charts' || true
}

files="$(find_by_interpreter bash sh | ignore)"
if [[ -n "$files" ]]; then
    echo "$files" | xargs -n 4 "${d}/../tools/run.sh" shellcheck -f gcc
fi
