// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"net/http"
	"time"

	"github.com/kemadev/go-framework/pkg/config"
)

type LivenessResponse struct {
	Timestamp   time.Time         `json:"timestamp"`
	Started     bool              `json:"started"`
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Checks      map[string]Status `json:"checks"`
}

// LivenessHandler handles liveness checks.
func LivenessHandler(conf config.Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// GetLivenessStatus return liveness status.
func GetLivenessStatus(checks CheckResults) Status {
	// TODO Obviously implement this status check
	return StatusOK
}

// CheckLiveness performs liveness checks and returns a map of results
// It always returns [StatusOK] as responding to liveness probe via HTTP means that the app is alive.
func CheckLiveness() CheckResults {
	// TODO Obviously implement this status check
	return CheckResults{
		"http": StatusOK,
	}
}
