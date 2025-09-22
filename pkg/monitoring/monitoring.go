// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

// Package monitoring defines monitoring endpoints handlers as well as checker functions,
// used to determine the health status of the application.
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
	HTTPLivenessCheckPath  = "/healthz"
	HTTPReadinessCheckPath = "/readyz"
	// HTTPLivenessCheckPattern is the pattern for the liveness check over HTTP.
	// It is used by Kubernetes to check if the application is alive.
	HTTPLivenessCheckPattern = "GET " + HTTPLivenessCheckPath
	// HTTPReadinessCheckPattern is the pattern for the readiness check over HTTP.
	// It is used by Kubernetes to check if the application is ready to serve traffic, as
	// well as returning some metrics.
	HTTPReadinessCheckPattern = "GET " + HTTPReadinessCheckPath
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

// MarshalText satisfies [encoding.TextMarshaller].
func (s Status) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
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

func (checks CheckResults) Status() Status {
	status := StatusOK

	for _, val := range checks {
		if val > status {
			status = val
		}
	}

	return status
}
