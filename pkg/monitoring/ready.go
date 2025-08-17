// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
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
	Ready          Status            `json:"ready"`
	RuntimeMetrics RuntimeMetrics    `json:"runtimeMetrics"`
	Checks         map[string]Status `json:"checks"`
}

// CheckReadiness performs services checks and returns a map of results.
func CheckReadiness() CheckResults {
	return CheckResults{}
}

// ReadinessHandler returns the pattern that should handle readiness checks, as well as associated readiness checking function.
// The function that is passed is used as status checker.
func ReadinessHandler(
	func() CheckResults,
) (string, http.HandlerFunc) {
	return HTTPReadinessCheckPattern, func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		usagePercent := float64(0)
		if m.Sys > 0 {
			usagePercent = (float64(m.Alloc) / float64(m.Sys)) * 100
		}

		checks := CheckReadiness()

		status := ReadinessResponse{
			Timestamp: time.Now().UTC(),
			Ready:     checks.Status(),
			Checks:    checks,
			RuntimeMetrics: RuntimeMetrics{
				Memory: MemoryMetrics{
					UsedBytes:    float64(m.Alloc),
					MaxBytes:     float64(m.Sys),
					UsagePercent: usagePercent,
					GCRuns:       m.NumGC,
				},
				CPU: CPUMetrics{
					Goroutines: runtime.NumGoroutine(),
				},
			},
		}

		body, err := json.Marshal(status)
		if err != nil {
			status.Checks["jsonMarshal"] = StatusDown
		}

		w.Header().Set(headkey.ContentType, headval.MIMEApplicationJSONCharsetUTF8)
		w.WriteHeader(status.Ready.HTTPCode())
		w.Write(body)
	}
}
