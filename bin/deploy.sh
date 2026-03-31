#!/bin/bash

set -ex

readonly d="$(cd "$(dirname "$0")/.." || exit 1; pwd)"

"${d}/tools/run.sh" kubectl get pod -o wide --watch &
"${d}/tools/run.sh" stern --tail 10 pneutrinoutil &

cleanup() {
    echo "CLEANUP"
    pkill "stern" || true
    pkill "kubectl" || true
}
trap cleanup EXIT

echo "START deploy"
timeout 180s "${d}/tools/run.sh" helm upgrade pneutrinoutil ./charts/pneutrinoutil --install --debug --wait-for-jobs "$@"
echo "END deploy"
