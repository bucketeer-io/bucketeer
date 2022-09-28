{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "kafka-cluster.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kafka-cluster.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion value to use for the kafka-cluster managed k8s resources
*/}}
{{- define "kafka-cluster.apiVersion" -}}
{{- if lt .Values.image.tag "v0.12.0" }}
{{- printf "%s" "monitoring.coreos.com/v1alpha1" -}}
{{- else -}}
{{- printf "%s" "monitoring.coreos.com/v1" -}}
{{- end -}}
{{- end -}}
