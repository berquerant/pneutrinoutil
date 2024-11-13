#!/bin/bash

target="$1"
shift

bin="./dist/NEUTRINO/bin"

if [ -z "$target" ] ; then
    echo "Require target:"
    ls -1 "$bin" | grep -v "dylib"
    exit 1
fi

export DYLD_LIBRARY_PATH=$PWD/dist/NEUTRINO/bin
"${bin}/${target}" "$@"
