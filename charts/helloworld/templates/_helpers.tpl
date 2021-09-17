{{/*
Expand the name of the chart.
*/}}
{{- define "helloworld.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "helloworld.fullname" -}}
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
{{- define "helloworld.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "helloworld.labels" -}}
helm.sh/chart: {{ include "helloworld.chart" . }}
{{ include "helloworld.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "helloworld.selectorLabels" -}}
app.kubernetes.io/name: {{ include "helloworld.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "helloworld.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "helloworld.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the domain
*/}}
{{- define "helloworld.domain" -}}
{{- if .Values.helloworld.domain -}}
{{- .Values.helloworld.domain -}}
{{- else -}}
{{- printf "%s.%s.svc.cluster.local" (include "helloworld.fullname" .) .Release.Namespace -}}
{{- end -}}
{{- end -}}

{{/*
Create the certs volume
*/}}
{{- define "helloworld.configMapCerts" -}}
{{- if .Values.helloworld.root_ca -}}
{{- printf "%s-certs" (include "helloworld.fullname" .) -}}
{{- else -}}
{{- .Values.helloworld.configMapCerts -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the certificate issuer
*/}}
{{- define "helloworld.certIssuerName" -}}
{{- if .Values.certIssuer.name -}}
{{ .Values.certIssuer.name }}
{{- else -}}
{{- printf "%s-issuer" (include "helloworld.fullname" .) -}}
{{- end -}}
{{- end -}}

