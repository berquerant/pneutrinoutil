#!/bin/bash

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly name="pneutrinoutil"

readonly uv_version="0.11.3"
readonly pnpm_version="10.33.0"

readonly target_ref="${TARGET_REF}"

start() {
    limactl start \
            --name "$name" \
            --yes \
            --mount-none \
            --cpus=4 \
            --memory=4 \
            --disk=50 \
            template:docker
    limactl copy "${d}/lima-setup.sh" "${name}:/tmp/"
    limactl shell "$name" /tmp/lima-setup.sh "${uv_version}" "${pnpm_version}" "${target_ref}"
}

stop() {
    limactl stop "$name"
}

reload() {
    stop || true
    limactl delete "$name" || true
    start
}

ssh() {
    exec limactl shell "$name" "$@"
}

run() {
    local __script
    __script="$(mktemp)"
    cat <<EOS > "$__script"
#!/bin/bash
set -ex
source "\${HOME}/.bashrc"
cd pneutrinoutil
./task init
direnv allow
direnv exec . $@
EOS
    chmod +x "$__script"
    limactl copy "$__script" "${name}:/tmp/run.sh"
    exec limactl shell "$name" /tmp/run.sh
}

set -ex
readonly cmd="$1"
case "$cmd" in
    "ssh")
        shift
        ssh "$@"
        ;;
    "start") start ;;
    "stop") stop ;;
    "reload") reload ;;
    "run")
        shift
        run "$@"
        ;;
    *)
        set +x
        cat <<EOS
Usage:
$0 start
  start VM
$0 stop
  stop VM
$0 reload
  delete and start new VM
$0 ssh [COMMAND...]
  ssh to VM
$0 run ARG...
  run command in VM
EOS
        exit 1
        ;;
esac
