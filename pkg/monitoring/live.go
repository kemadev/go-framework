// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"context"
	"net/http"
	"time"

	"github.com/kemadev/go-framework/pkg/config"
	khttp "github.com/kemadev/go-framework/pkg/http"
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
		kclient := khttp.ClientInfo{
			Ctx:    context.Background(),
			Writer: w,
		}

		checks := CheckLiveness()
		status := GetLivenessStatus(checks)

		khttp.SendJSONResponse(
			kclient,
			status.HTTPCode(),
			LivenessResponse{
				Timestamp:   time.Now().UTC(),
				Started:     true,
				Status:      status.String(),
				Version:     conf.AppVersion,
				Environment: conf.Environment,
				Checks:      checks,
			},
		)
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
