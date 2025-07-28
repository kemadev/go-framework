package monitoring

import (
	"context"
	"net/http"
	"runtime"
	"time"

	khttp "github.com/kemadev/go-framework/pkg/http"
)

type RuntimeMetrics struct {
	Memory MemoryMetrics `json:"memory"`
	CPU    CPUMetrics    `json:"cpu"`
}

type MemoryMetrics struct {
	UsedBytes    float64 `json:"used_bytes"`
	MaxBytes     float64 `json:"max_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	GCRuns       uint32  `json:"gc_runs"`
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

// ReadinessHandler handles readiness checks
func ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kclient := khttp.ClientInfo{
			Ctx:    context.Background(),
			Writer: w,
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		usagePercent := float64(0)
		if m.Sys > 0 {
			usagePercent = (float64(m.Alloc) / float64(m.Sys)) * 100
		}

		checks := CheckReadiness()
		status := GetReadinessStatus(checks)

		khttp.SendJSONResponse(
			kclient,
			status.HTTPCode(),
			ReadinessResponse{
				Timestamp: time.Now().UTC(),
				Ready:     status.IsReady(),
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
			},
		)
	}
}

// GetLivenessStatus return readiness status
func GetReadinessStatus(checks CheckResults) Status {
	return StatusOK
}

// CheckReadiness performs services checks and returns a map of results
func CheckReadiness() CheckResults {
	return CheckResults{}
}
