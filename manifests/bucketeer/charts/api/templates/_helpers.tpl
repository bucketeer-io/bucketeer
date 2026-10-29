{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "api.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "api.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "api.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "service-cert-secret" -}}
{{- if .Values.tls.service.secret }}
{{- printf "%s" .Values.tls.service.secret -}}
{{- else -}}
{{ template "api.fullname" . }}-service-cert
{{- end -}}
{{- end -}}

{{- define "service-token-secret" -}}
{{- if .Values.serviceToken.secret }}
{{- printf "%s" .Values.serviceToken.secret -}}
{{- else -}}
{{ template "api.fullname" . }}-service-token
{{- end -}}
{{- end -}}

{{- define "api.db-credentials-secret" -}}
{{ template "api.fullname" . }}-db-credentials
{{- end -}}

{{/*
Name of the existing Secret that holds the password of the given database
engine ("mysql" or "postgres"), or an empty string when the chart renders the
password from values into its own Secret.
Usage: include "api.db-password-existing-secret" (dict "root" . "engine" "mysql")
*/}}
{{- define "api.db-password-existing-secret" -}}
{{- $local := index .root.Values.env.operationalDatabase .engine -}}
{{- $global := index .root.Values.global.operationalDatabase .engine -}}
{{- $local.existingSecret | default $global.existingSecret | default "" -}}
{{- end -}}

{{/*
secretKeyRef fields (name and key) of the password env var of the given
database engine ("mysql" or "postgres").
Usage: include "api.db-password-secret-key-ref" (dict "root" . "engine" "mysql")
*/}}
{{- define "api.db-password-secret-key-ref" -}}
{{- $local := index .root.Values.env.operationalDatabase .engine -}}
{{- $global := index .root.Values.global.operationalDatabase .engine -}}
{{- $defaultKey := printf "%s-pass" .engine -}}
{{- $existing := include "api.db-password-existing-secret" . -}}
{{- if $existing -}}
name: {{ $existing | quote }}
key: {{ $local.existingSecretKey | default $global.existingSecretKey | default $defaultKey | quote }}
{{- else -}}
name: {{ template "api.db-credentials-secret" .root }}
key: {{ $defaultKey }}
{{- end -}}
{{- end -}}
