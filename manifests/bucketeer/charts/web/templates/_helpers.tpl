{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "web.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "web.fullname" -}}
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
{{- define "web.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "service-cert-secret" -}}
{{- if .Values.tls.service.secret }}
{{- printf "%s" .Values.tls.service.secret -}}
{{- else -}}
{{ template "web.fullname" . }}-service-cert
{{- end -}}
{{- end -}}

{{- define "oauth-key-secret" -}}
{{- if .Values.oauth.key.secret }}
{{- printf "%s" .Values.oauth.key.secret -}}
{{- else -}}
{{ template "web.fullname" . }}-oauth-key
{{- end -}}
{{- end -}}

{{- define "service-token-secret" -}}
{{- if .Values.serviceToken.secret }}
{{- printf "%s" .Values.serviceToken.secret -}}
{{- else -}}
{{ template "web.fullname" . }}-service-token
{{- end -}}
{{- end -}}

{{- define "web.oauth-config-secret" -}}
{{- if .Values.oauth.existingSecret }}
{{- printf "%s" .Values.oauth.existingSecret -}}
{{- else -}}
{{ template "web.fullname" . }}-oauth-config
{{- end -}}
{{- end -}}

{{- define "web.email-config-secret" -}}
{{- $email := .Values.global.email | default dict -}}
{{- if $email.existingSecret }}
{{- printf "%s" $email.existingSecret -}}
{{- else -}}
{{ template "web.fullname" . }}-email-config
{{- end -}}
{{- end -}}

{{- define "web.db-credentials-secret" -}}
{{ template "web.fullname" . }}-db-credentials
{{- end -}}

{{/*
Name of the existing Secret that holds the password of the given database
engine ("mysql" or "postgres"), or an empty string when the chart renders the
password from values into its own Secret.
Usage: include "web.db-password-existing-secret" (dict "root" . "engine" "mysql")
*/}}
{{- define "web.db-password-existing-secret" -}}
{{- $local := index .root.Values.env.operationalDatabase .engine -}}
{{- $global := index .root.Values.global.operationalDatabase .engine -}}
{{- $local.existingSecret | default $global.existingSecret | default "" -}}
{{- end -}}

{{/*
secretKeyRef fields (name and key) of the password env var of the given
database engine ("mysql" or "postgres").
Usage: include "web.db-password-secret-key-ref" (dict "root" . "engine" "mysql")
*/}}
{{- define "web.db-password-secret-key-ref" -}}
{{- $local := index .root.Values.env.operationalDatabase .engine -}}
{{- $global := index .root.Values.global.operationalDatabase .engine -}}
{{- $defaultKey := printf "%s-pass" .engine -}}
{{- $existing := include "web.db-password-existing-secret" . -}}
{{- if $existing -}}
name: {{ $existing | quote }}
key: {{ $local.existingSecretKey | default $global.existingSecretKey | default $defaultKey | quote }}
{{- else -}}
name: {{ template "web.db-credentials-secret" .root }}
key: {{ $defaultKey }}
{{- end -}}
{{- end -}}
