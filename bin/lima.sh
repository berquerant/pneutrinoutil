#!/bin/bash

readonly d="$(cd "$(dirname "$0")" || exit 1; pwd)"
readonly name="pneutrinoutil"

readonly uv_version=0.11.3
readonly pnpm_version=10.33.0

readonly vm_repo_dir="pneutrinoutil"

readonly target_ref="${TARGET_REF}"

readonly cpus="${VM_CPUS:-4}"
readonly memory="${VM_MEMORY_GB:-4}"
readonly disk="${VM_DISK_GB:-50}"
readonly host_cache_dir="${HOST_CACHE_DIR:-${d}/../tmp/lima}"
readonly vm_cache_dir="${VM_CACHE_DIR:-/tmp/cache}"

readonly go_cache_dir="${vm_cache_dir}/go/cache"
readonly gomod_cache_dir="${vm_cache_dir}/go/modcache"
readonly docker_cache_dir="${vm_cache_dir}/docker"

start() {
    local __spec
    __spec="$(mktemp).yaml"
    cat <<EOS > "$__spec"
base:
  - template:docker
cpus: ${cpus}
memory: "${memory}GiB"
disk: "${disk}GiB"
EOS
    local __cmd="limactl start --name ${name} --yes ${__spec}"
    if [[ -n "$host_cache_dir" && -n "$vm_cache_dir" ]] ; then
        mkdir -p "$host_cache_dir"
        $__cmd --set=".mounts = [{\"location\": \"${host_cache_dir}\", \"mountPoint\": \"${vm_cache_dir}\", \"writable\": true}]"
    else
        $__cmd
    fi

    limactl copy "${d}/lima-setup.sh" "${name}:/tmp/"
    limactl shell "$name" /tmp/lima-setup.sh \
            "${vm_repo_dir}" \
            "${uv_version}" \
            "${pnpm_version}" \
            "${target_ref}"
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
cd "$vm_repo_dir"
export CACHEDIR="${vm_cache_dir}"
export GOCACHE="${go_cache_dir}"
export GOMODCACHE="${gomod_cache_dir}"
export DOCKERCACHE="${docker_cache_dir}"
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
