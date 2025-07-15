package config

const (
	// HTTPLivenessCheckPath is the path for the liveness check over HTTP.
	// It is used by Kubernetes to check if the application is alive.
	HTTPLivenessCheckPath = "/healthz"
	// HTTPReadinessCheckPath is the path for the readiness check over HTTP.
	// It is used by Kubernetes to check if the application is ready to serve traffic, as
	// well as returning some metrics.
	HTTPReadinessCheckPath = "/readyz"
)
