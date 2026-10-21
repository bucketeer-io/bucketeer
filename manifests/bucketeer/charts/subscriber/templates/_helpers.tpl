{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "subscriber.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "subscriber.fullname" -}}
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
{{- define "subscriber.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "service-cert-secret" -}}
{{- if .Values.tls.service.secret }}
{{- printf "%s" .Values.tls.service.secret -}}
{{- else -}}
{{ template "subscriber.fullname" . }}-service-cert
{{- end -}}
{{- end -}}

{{- define "service-token-secret" -}}
{{- if .Values.serviceToken.secret }}
{{- printf "%s" .Values.serviceToken.secret -}}
{{- else -}}
{{ template "subscriber.fullname" . }}-service-token
{{- end -}}
{{- end -}}


{{- define "issuer-cert-secret" -}}
{{- if .Values.tls.issuer.secret }}
{{- printf "%s" .Values.tls.issuer.secret -}}
{{- else -}}
{{ template "subscriber.fullname" . }}-issuer-cert
{{- end -}}
{{- end -}}

{{- define "subscriber.deploy" -}}
{{- if not .Values.gcpMultiCluster.enabled }}
  {{- true }}
{{- else if and .Values.gcpMultiCluster.enabled .Values.gcpMultiCluster.configCluster }}
  {{- true }}
{{- else }}
  {{- false }}
{{- end }}
{{- end -}}

{{- /* Overwrites a processor's dataWarehouse settings with the global ones */}}
{{- define "subscriber.injectDataWarehouse" -}}
{{- $g := .global.dataWarehouse }}
{{- $_ := set .config "dataWarehouse" (dict
  "type" $g.type
  "batchSize" $g.batchSize
  "timezone" $g.timezone
  "mysql" (dict
    "useMainConnection" $g.mysql.useMainConnection
    "host" $g.mysql.host
    "port" $g.mysql.port
    "user" $g.mysql.user
    "password" $g.mysql.password
    "database" $g.mysql.database)
  "postgres" (dict
    "host" $g.postgres.host
    "port" $g.postgres.port
    "user" $g.postgres.user
    "password" $g.postgres.password
    "database" $g.postgres.database
    "sslMode" (default "" $g.postgres.sslMode)
    "sslRootCert" (default "" $g.postgres.sslRootCert)
    "sslCert" (default "" $g.postgres.sslCert)
    "sslKey" (default "" $g.postgres.sslKey))
  "bigquery" (dict
    "project" $g.bigquery.project
    "dataset" $g.bigquery.dataset
    "location" $g.bigquery.location)) }}
{{- end -}}
