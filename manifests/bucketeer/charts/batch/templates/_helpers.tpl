{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "batch-server.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "batch-server.fullname" -}}
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
{{- define "batch-server.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "service-cert-secret" -}}
{{- if .Values.tls.service.secret }}
{{- printf "%s" .Values.tls.service.secret -}}
{{- else -}}
{{ template "batch-server.fullname" . }}-service-cert
{{- end -}}
{{- end -}}

{{- define "service-token-secret" -}}
{{- if .Values.serviceToken.secret }}
{{- printf "%s" .Values.serviceToken.secret -}}
{{- else -}}
{{ template "batch-server.fullname" . }}-service-token
{{- end -}}
{{- end -}}

{{- define "oauth-key-secret" -}}
{{- if .Values.oauth.key.secret }}
{{- printf "%s" .Values.oauth.key.secret -}}
{{- else -}}
{{ template "backend.fullname" . }}-oauth-key
{{- end -}}
{{- end -}}

{{- define "issuer-cert-secret" -}}
{{- if .Values.tls.issuer.secret }}
{{- printf "%s" .Values.tls.issuer.secret -}}
{{- else -}}
{{ template "batch-server.fullname" . }}-issuer-cert
{{- end -}}
{{- end -}}

{{- define "batch.deploy" -}}
{{- if not .Values.gcpMultiCluster.enabled }}
  {{- true }}
{{- else if and .Values.gcpMultiCluster.enabled .Values.gcpMultiCluster.configCluster }}
  {{- true }}
{{- else }}
  {{- false }}
{{- end }}
{{- end }}
