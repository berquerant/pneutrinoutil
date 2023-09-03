#!/bin/bash

script_path="$0"
cd "$(dirname $0)"

default_config="./skeleton.json"

# Print usage
usage() {
    cat << EOF
usage: $(basename $script_path) [-n] [-c FILE]

INSTALLATION
  NEUTRINO/$(basename $script_path)

ENVIRONMENT VARIABLES
  DEBUG_FILE
    Dry run if set.
    Record the commands to be executed in the file.

OPTIONS
  -c
    Specify the config file.
    Default: ${default_config}

  -h
    Show this help.

  -n
    Write the config skeleton to stdout.

SKELETON
  Description
    Any, for memo.

  Play
    Bool, play the output wav file if true.

  Score
    String, input musicxml file path.

  Output
    String, output diretory path.

  RealtimeSynthesis
    A WORLD option.

  RandomSeed
    A NEUTRINO option.

See NEUTRINO readme for other keys.
EOF
}

# Print skeleton
skeleton() {
    cat << EOF
{
  "Description": {},
  "Play": true,
  "Score": "./score/musicxml/sample1.musicxml",
  "Output": "./output",
  "NumThreads": 4,
  "InferenceMode": 3,
  "ModelDir": "MERROW",
  "StyleShift": 0,
  "PitchShiftNsf": 0,
  "PitchShiftWorld": 0,
  "FormantShift": 1.0,
  "SmoothPitch": 0,
  "SmoothFormant": 0,
  "EnhanceBreathiness": 0,
  "RealtimeSynthesis": false,
  "RandomSeed": 1234
}
EOF
}

# parse cli arguments

config="$default_config"

while getopts 'hnc:' opt; do
    case "${opt}" in
        n)
            skeleton
            exit
            ;;
        c)
            config="${OPTARG}"
            ;;
        h)
            usage
            exit 1
            ;;
        *)
            usage
            exit 1
            ;;
    esac
done

# make binaries executable
chmod 755 ./bin/*
xattr -dr com.apple.quarantine ./bin

gen_salt() {
    date +"%Y%m%d%H%M%S"
}

# salt to uniquify basename
salt="$(gen_salt)"

is_debug() {
    [ -n "$DEBUG_FILE" ]
}

# ensure that empty debug file exists
if is_debug ; then
    touch "$DEBUG_FILE"
fi

# save executed commands
log_file=$(mktemp)

run_or_dry() {
    if is_debug ; then
        echo "$@" >> "$DEBUG_FILE"
    else
        echo "$@" >> "$log_file"
        "$@"
    fi
}

run_and_log() {
    if is_debug ; then
        echo "$@" >> "$DEBUG_FILE"
    fi
    echo "$@" >> "$log_file"
    "$@"
}

# save config to protect against rewriting during execution
temp_config=$(mktemp)
cp "$config" "$temp_config"

# get arg by key from config file
arg() {
    jq -r "$1" "$temp_config"
}

# echo score file name except extension
get_score_base() {
    x=$(arg ".Score")
    x=$(basename $x)
    echo "${x%%.*}"
}

get_basename() {
    echo "$(get_score_base)_${salt}"
}

save_result() {
    dst=$(arg ".Output")/$(get_basename)
    run_and_log mkdir -p "$dst"

    result="${dst}/$(get_basename).json"
    input_xml="$(get_basename).musicxml"
    output_wav="$(get_basename).wav"

    query='{"skeleton": ., "score": $i, "out": $o}'
    if is_debug ; then
        echo "jq --arg i \"$input_xml\" --arg o \"$output_wav\" \"$query\" \"$config\" > \"$result\"" >> "$DEBUG_FILE"
    else
        jq --arg i "$input_xml" --arg o "$output_wav" "$query" "$config" > "$result"
    fi
    run_or_dry cp "$(gen_name_score_with_salt)" "$dst"
    run_or_dry cp "$(gen_name_output_wav)" "$dst"
    run_or_dry cp "$script_path" "$dst"
    run_or_dry cp "$log_file" "$dst/script.log"
}

gen_name_score_with_salt() {
    echo "./score/musicxml/$(get_basename).musicxml"
}

copy_score_with_salt() {
    src=$(arg ".Score")
    dst="$(gen_name_score_with_salt)"
    run_or_dry cp -f "$src" "$dst"
}

gen_name_output_wav() {
    echo "./output/$(get_basename).wav"
}

play_output_wav() {
    run_or_dry open "$(gen_name_output_wav)"
}

play_after_run() {
    if arg ".Play" ; then
        play_output_wav
    fi
}

# retry MAX_ATTEMPTS INTERVAL_SECOND COMMAND...
retry() {
  retries=$1
  interval=$2
  shift 2

  for i in $(seq $retries); do
    if "$@"; then
      return 0
    fi

    sleep "$interval"
    echo "RETRY:${i}"
  done

  echo "RETRY:exhausted!"
  return 1
}

num_threads() {
    arg ".NumThreads"
}
inference_mode() {
    arg ".NumThreads"
}
model_dir() {
    arg ".ModelDir"
}
style_shift() {
    arg ".StyleShift"
}
pitch_shift_nsf() {
    arg ".PitchShiftNsf"
}
pitch_shift_world() {
    arg ".PitchShiftWorld"
}
formant_shift() {
    arg ".FormantShift"
}
smooth_pitch() {
    arg ".SmoothPitch"
}
smooth_formant() {
    arg ".SmoothFormant"
}
enhance_breathiness() {
    arg ".EnhanceBreathiness"
}
realtime_synthesis() {
    arg ".RealtimeSynthesis"
}
random_seed() {
    arg ".RandomSeed"
}

# Project settings
BASENAME="$(get_basename)"
NumThreads="$(num_threads)"
InferenceMode="$(inference_mode)"

# musicXML_to_label
SUFFIX=musicxml

# NEUTRINO
ModelDir="$(model_dir)"
StyleShift="$(style_shift)"
RandomSeed="$(random_seed)"

# NSF
PitchShiftNsf="$(pitch_shift_nsf)"

# WORLD
PitchShiftWorld="$(pitch_shift_world)"
FormantShift="$(formant_shift)"
SmoothPitch="$(smooth_pitch)"
SmoothFormant="$(smooth_formant)"
EnhanceBreathiness="$(enhance_breathiness)"
RealtimeSynthesis=""
if realtime_synthesis ; then
    RealtimeSynthesis="-r"
fi

if [ ${InferenceMode} -eq 4 ]; then
    NsfModel=va
    SamplingFreq=48
elif [ ${InferenceMode} -eq 3 ]; then
    NsfModel=vs
    SamplingFreq=48
elif [ ${InferenceMode} -eq 2 ]; then
    NsfModel=ve
    SamplingFreq=24
fi

run_musicxmltolabel() {
    echo "`date +"%M:%S"` : start MusicXMLtoLabel"
    run_or_dry ./bin/musicXMLtoLabel score/musicxml/${BASENAME}.${SUFFIX} score/label/full/${BASENAME}.lab score/label/mono/${BASENAME}.lab
}

run_neutrino() {
    echo "`date +"%M:%S"` : start NEUTRINO"
    run_or_dry ./bin/NEUTRINO score/label/full/${BASENAME}.lab score/label/timing/${BASENAME}.lab ./output/${BASENAME}.f0 ./output/${BASENAME}.melspec ./model/${ModelDir}/ -w ./output/${BASENAME}.mgc ./output/${BASENAME}.bap -n 1 -o ${NumThreads} -k ${StyleShift} -d ${InferenceMode} -t -r "${RandomSeed}"
}

run_nsf() {
    echo "`date +"%M:%S"` : start NSF"
    run_or_dry ./bin/NSF output/${BASENAME}.f0 output/${BASENAME}.melspec ./model/${ModelDir}/${NsfModel}.bin output/${BASENAME}.wav -l score/label/timing/${BASENAME}.lab -n 1 -p ${NumThreads} -s ${SamplingFreq} -f ${PitchShiftNsf} -t
}

run_world() {
    echo "`date +"%M:%S"` : start WORLD"
    run_or_dry ./bin/WORLD output/${BASENAME}.f0 output/${BASENAME}.mgc output/${BASENAME}.bap output/${BASENAME}_world.wav -f ${PitchShiftWorld} -m ${FormantShift} -p ${SmoothPitch} -c ${SmoothFormant} -b ${EnhanceBreathiness} -n ${NumThreads} -t "${RealtimeSynthesis}"
}

on_error() {
    echo "SCRIPT LOG: ${log_file}" >&2
    exit 1
}

main() {
    trap on_error ERR
    run_or_dry set -ex
    # PATH to current library
    run_or_dry export DYLD_LIBRARY_PATH=$PWD/bin:$DYLD_LIBRARY_PATH
    copy_score_with_salt
    set +e
    retry 10 1 run_musicxmltolabel || exit 1
    retry 10 1 run_neutrino || exit 1
    retry 10 1 run_nsf || exit 1
    # retry to avoid segfault
    retry 10 1 run_world || exit 1
    set -e
    echo "`date +"%M:%S"` : END"
    play_after_run
    save_result
}

main
