// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package config

const (
	// EnvVarKeyRuntimeEnv is the environment variable key used to set runtime environment.
	EnvVarKeyRuntimeEnv = "RUNTIME_ENV"
	// EnvVarKeyAppName is the environment variable key used to set application name.
	EnvVarKeyAppName = "APP_NAME"
	// EnvVarKeyAppVersion is the environment variable key used to set application version.
	EnvVarKeyAppVersion = "APP_VERSION"
	// EnvVarKeyAppNamespace is the environment variable key used to set application namespace.
	EnvVarKeyAppNamespace = "APP_NAMESPACE"
	// EnvVarKeyOtelEndpointURL is the environment variable key used to set OpenTelemetry endpoint URL.
	EnvVarKeyOtelEndpointURL = "OTEL_ENDPOINT_URL"
	// EnvVarKeyOtelExporterCompression is the environment variable key used to set OpenTelemetry exporter compression method.
	EnvVarKeyOtelExporterCompression = "OTEL_EXPORTER_COMPRESSION"
	// EnvVarKeyHTTPServePort is the environment variable key used to set HTTP serving port.
	EnvVarKeyHTTPServePort = "HTTP_SERVE_PORT"
	// EnvVarKeyHTTPReadTimeout is the environment variable key used to set HTTP read timeout.
	EnvVarKeyHTTPReadTimeout = "HTTP_READ_TIMEOUT"
	// EnvVarKeyHTTPWriteTimeout is the environment variable key used to set write timeout.
	EnvVarKeyHTTPWriteTimeout = "HTTP_WRITE_TIMEOUT"
	// EnvVarKeyHTTPIdleTimeout is the environment variable key used to set idle timeout.
	EnvVarKeyHTTPIdleTimeout = "HTTP_IDLE_TIMEOUT"
	// EnvVarKeyMetricsExportInterval is the environment variable key used to set metrics export interval.
	// A negative value in development mode will disable metrics export.
	EnvVarKeyMetricsExportInterval = "METRICS_EXPORT_INTERVAL"
	// EnvVarKeyTracesSampleRatio is the environment variable key used to set traces sample ratio.
	EnvVarKeyTracesSampleRatio = "TRACES_SAMPLE_RATIO"
	// EnvVarKeyBusinessUnitID is the environment variable key used to identify business unit id.
	EnvVarKeyBusinessUnitID = "BUSINESS_UNIT_ID"
	// EnvVarKeyCustomerID is the environment variable key used to identify customer id.
	EnvVarKeyCustomerID = "CUSTOMER_ID"
	// EnvVarKeyCostCenter is the environment variable key used to identify cost center.
	EnvVarKeyCostCenter = "COST_CENTER"
	// EnvVarKeyCostAllocationOwner is the environment variable key used to identify cost allocation owner.
	EnvVarKeyCostAllocationOwner = "COST_ALLOCATION_OWNER"
	// EnvVarKeyOperationsOwner is the environment variable key used to identify operations owner.
	EnvVarKeyOperationsOwner = "OPERATIONS_OWNER"
	// EnvVarKeyRpo is the environment variable key used to identify recovery point objective.
	EnvVarKeyRpo = "RPO"
	// EnvVarKeyDataClassification is the environment variable key used to identify data classification.
	EnvVarKeyDataClassification = "DATA_CLASSIFICATION"
	// EnvVarKeyComplianceFramework is the environment variable key used to identify compliance framework.
	EnvVarKeyComplianceFramework = "COMPLIANCE_FRAMEWORK"
	// EnvVarKeyExpiration is the environment variable key used to identify application expiration.
	EnvVarKeyExpiration = "EXPIRATION"
	// EnvVarKeyProjectURL is the environment variable key used to identify application project url.
	EnvVarKeyProjectURL = "PROJECT_URL"
	// EnvVarKeyMonitoringURL is the environment variable key used to identify application monitoring url.
	EnvVarKeyMonitoringURL = "MONITORING_URL"
)
