// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
)

type LivenessResponse struct {
	Timestamp   time.Time         `json:"timestamp"`
	Started     bool              `json:"started"`
	Status      Status            `json:"status"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Checks      map[string]Status `json:"checks"`
}

// LivenessHandler returns the pattern that should handle liveness checks, as well as associated liveness checking function
// The function that is passed is used as status checker.
func LivenessHandler(
	livenessChecker func() CheckResults,
) (string, http.HandlerFunc) {
	return HTTPLivenessCheckPattern, func(w http.ResponseWriter, r *http.Request) {
		checks := livenessChecker()

		status := LivenessResponse{
			Timestamp: time.Now(),
			// If serving HTTP, we started
			Started: true,
			Status:  checks.Status(),
			Checks:  checks,
		}

		conf, err := config.NewManager().Load()
		if err != nil {
			status.Checks["config"] = StatusDown
		}

		status.Version = conf.Runtime.AppVersion
		status.Environment = conf.Runtime.Environment

		body, err := json.Marshal(status)
		if err != nil {
			status.Checks["jsonMarshal"] = StatusDown
		}

		w.Header().Set(headkey.ContentType, headval.AcceptJSON)
		w.WriteHeader(status.Status.HTTPCode())
		w.Write(body)
	}
}
