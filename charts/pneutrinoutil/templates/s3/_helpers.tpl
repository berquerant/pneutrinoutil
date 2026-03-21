{{- define "pneutrinoutil.s3.bucketTxt" -}}
pneutrinoutil
test
{{- end }}

{{- define "pneutrinoutil.s3.setupSh" -}}
#!/bin/bash

readonly bucket_txt="$BUCKET_TXT"

client() {
  aws s3 "$@"
}

exist_bucket() {
  local -r bucket="$1"
  client ls | awk '{print $3}' | grep -q "${bucket}"
}

create_bucket() {
  local -r bucket="$1"
  client mb "s3://${bucket}"
}

set -e
set -o pipefail
while read -r b ; do
  echo >&2 "ensure ${b}"
  if ! exist_bucket "$b" ; then
    create_bucket "$b"
  fi
done < "$bucket_txt"
{{- end }}

{{- define "pneutrinoutil.s3.configSetup" -}}
bucket.txt: |
  {{- include "pneutrinoutil.s3.bucketTxt" . | nindent 4 }}
setup.sh: |
  {{- include "pneutrinoutil.s3.setupSh" . | nindent 4 }}
{{- end }}
