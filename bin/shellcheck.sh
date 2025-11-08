#!/bin/bash

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"

find_by_shebang() {
    git grep -l "^${1}"
}

find_by_extension() {
    git ls-files | grep "\.${1}$"
}

find_by_interpreter() {
    local -r interpreter="$1"
    local -r extension="$2"
    {
        find_by_shebang "#\!/bin/${interpreter}"
        find_by_extension "$extension"
    } | sort -u
}


find_by_interpreter bash sh | xargs -n 4 "${d}/../tools/run.sh" shellcheck -f gcc
