{{- define "upi.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "upi.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- printf "%s" $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{- define "upi.labels" -}}
app.kubernetes.io/name: {{ include "upi.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end }}

{{- define "upi.selectorLabels" -}}
app.kubernetes.io/name: {{ include "upi.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "upi.postgresDsn" -}}
postgres://{{ .Values.global.postgres.user }}:{{ .Values.global.postgres.password }}@{{ .Values.global.postgres.host }}:{{ .Values.global.postgres.port }}/{{ .Values.global.postgres.database }}?sslmode=disable
{{- end }}

{{- define "upi.serviceEnv" -}}
- name: KAFKA_BROKERS
  value: {{ .Values.global.kafka.brokers | quote }}
- name: REDIS_ADDR
  value: {{ .Values.global.redis.addr | quote }}
- name: OTEL_EXPORTER_OTLP_ENDPOINT
  value: {{ .Values.global.otel.endpoint | quote }}
- name: POSTGRES_DSN
  valueFrom:
    secretKeyRef:
      name: upi-secrets
      key: postgres-dsn
{{- end }}
