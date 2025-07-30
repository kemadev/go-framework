// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"net/http"
	"time"
)

type RuntimeMetrics struct {
	Memory MemoryMetrics `json:"memory"`
	CPU    CPUMetrics    `json:"cpu"`
}

type MemoryMetrics struct {
	UsedBytes    float64 `json:"usedBytes"`
	MaxBytes     float64 `json:"maxBytes"`
	UsagePercent float64 `json:"usagePercent"`
	GCRuns       uint32  `json:"gcRuns"`
}

type CPUMetrics struct {
	Goroutines int `json:"goroutines"`
}

type ReadinessResponse struct {
	Timestamp      time.Time         `json:"timestamp"`
	Ready          bool              `json:"ready"`
	RuntimeMetrics RuntimeMetrics    `json:"runtimeMetrics"`
	Checks         map[string]Status `json:"checks"`
}

// ReadinessHandler handles readiness checks.
func ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// GetLivenessStatus return readiness status.
func GetReadinessStatus(checks CheckResults) Status {
	return StatusOK
}

// CheckReadiness performs services checks and returns a map of results.
func CheckReadiness() CheckResults {
	return CheckResults{}
}
