#!/bin/bash

set -e
set -o pipefail
set -x

d="$(cd "$(dirname "$0")"/.. || exit 1; pwd)"
cd "$d"
rm -rf app/api/client/docs
npm run swag

__sed='sed'
if which gsed >/dev/null 2>&1 ; then
    __sed='gsed'
fi
"$__sed" -i 's|import type|import|g' app/api/client/api.ts
