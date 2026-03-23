{{/*
Expand the name of the chart.
*/}}
{{- define "pneutrinoutil.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "pneutrinoutil.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "pneutrinoutil.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "pneutrinoutil.labels" -}}
helm.sh/chart: {{ include "pneutrinoutil.chart" . }}
{{ include "pneutrinoutil.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "pneutrinoutil.selectorLabels" -}}
app.kubernetes.io/name: {{ include "pneutrinoutil.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "pneutrinoutil.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "pneutrinoutil.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "pneutrinoutil.s3Common" -}}
AWS_S3_DISABLE_HTTPS: "true"
AWS_USE_PATH_STYLE_ENDPOINT: "true"
AWS_DEFAULT_REGION: "us-east-1"
AWS_ENDPOINT_URL: http://{{ include "pneutrinoutil.fullname" . }}-s3:9000
AWS_ACCESS_KEY_ID: {{ .Values.s3.user }}
AWS_SECRET_ACCESS_KEY: {{ .Values.s3.password }}
{{- end }}

{{- define "pneutrinoutil.serverCommon" -}}
STORAGES3: "true"
STORAGEBUCKET: {{ .Values.s3.bucket }}
REDISDSN: redis://{{ include "pneutrinoutil.fullname" . }}-redis:6379/{{ .Values.redis.db }}
MYSQLDSN: {{ .Values.mysql.user }}:{{ .Values.mysql.password }}@tcp({{ include "pneutrinoutil.fullname" . }}-mysql:3306)/{{ .Values.mysql.database }}?parseTime=true&loc=Asia%2FTokyo
{{- end }}

{{- define "pneutrinoutil.retrySh" -}}
#!/bin/bash

limit=100
interval=3

while [[ $# -gt 0 ]] ; do
    case "$1" in
        "--limit")
            if [[ -z "$2" ]] ; then
                echo >&2 "$0: --limit requires an argument"
                exit 1
            fi
            limit="$2"
            shift 2
            ;;
        "--interval")
            if [[ -z "$2" ]] ; then
                echo >&2 "$0: --interval requires an argument"
                exit 1
            fi
            interval="$2"
            shift 2
            ;;
        "--")
            shift
            break
            ;;
        *)
            echo >&2 "$0: unknown arguments!"
            exit 1
            ;;
    esac
done

for c in in $(seq "$limit") ; do
    echo >&2 "retry[$1][${c}]..."
    if ((c > 1)) ; then
        sleep "$interval"
    fi
    if "$@" ; then
        exit
    fi
done

echo >&2 "retry[$1] failed!"
exit 1
{{- end }}

{{- define "pneutrinoutil.util" -}}
retry.sh: |
  {{- include "pneutrinoutil.retrySh" . | nindent 4 }}
{{- end }}
