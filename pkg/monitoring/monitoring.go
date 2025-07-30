// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"net/http"
)

// Status represents the health status of a service or component.
type Status int

const (
	// StatusOK means that the app is running fine.
	StatusOK Status = iota
	// StatusDegraded means that the app is running, yet has some minor, non-critical issues.
	StatusDegraded
	// StatusDown means that the app in unable to process requests.
	StatusDown
	// StatusUnknown means that it is impossible to determine the status of the application.
	StatusUnknown
)

const (
	// HTTPLivenessCheckPath is the path for the liveness check over HTTP.
	// It is used by Kubernetes to check if the application is alive.
	HTTPLivenessCheckPath = "/healthz"
	// HTTPReadinessCheckPath is the path for the readiness check over HTTP.
	// It is used by Kubernetes to check if the application is ready to serve traffic, as
	// well as returning some metrics.
	HTTPReadinessCheckPath = "/readyz"
)

// String returns the string representation of the status.
func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusDegraded:
		return "degraded"
	case StatusDown:
		return "down"
	case StatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

// HTTPCode returns the appropriate HTTP status code for the status.
func (s Status) HTTPCode() int {
	switch s {
	case StatusOK:
		return http.StatusOK
	case StatusDegraded:
		return http.StatusOK
	case StatusDown:
		return http.StatusServiceUnavailable
	case StatusUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsHealthy returns true if the status indicates the service is healthy.
// Status degraded is still considered as healthy.
func (s Status) IsHealthy() bool {
	return s == StatusOK || s == StatusDegraded
}

// IsReady returns true if the status indicates the service is ready
// Status degraded is still considered as ready.
func (s Status) IsReady() bool {
	return s == StatusOK || s == StatusDegraded
}

type CheckResults map[string]Status
