package monitoring

import (
	"context"
	"net/http"
	"runtime"
	"time"

	khttp "github.com/kemadev/go-framework/pkg/http"
	"github.com/kemadev/go-framework/pkg/route"
)

// Status represents the health status of a service or component
type Status int

const (
	StatusOK Status = iota
	StatusDegraded
	StatusDown
	StatusUnknown
)

// String returns the string representation of the status
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

// HTTPCode returns the appropriate HTTP status code for the status
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

// IsHealthy returns true if the status indicates the service is healthy
func (s Status) IsHealthy() bool {
	return s == StatusOK || s == StatusDegraded
}

// IsReady returns true if the status indicates the service is ready
func (s Status) IsReady() bool {
	return s == StatusOK || s == StatusDegraded
}

type LivenessResponse struct {
	Timestamp   time.Time         `json:"timestamp"`
	Started     bool              `json:"started"`
	Status      string            `json:"status"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Checks      map[string]string `json:"checks"`
}

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

type CheckResults map[string]Status

func (results *CheckResults) Pretty() map[string]string {
	res := make(map[string]string, len(*results))

	for key, status := range *results {
		res[key] = status.String()
	}

	return res
}

type ReadinessResponse struct {
	Timestamp      time.Time         `json:"timestamp"`
	Ready          bool              `json:"ready"`
	RuntimeMetrics RuntimeMetrics    `json:"runtimeMetrics"`
	Checks         map[string]string `json:"checks"`
}

// Routes returns monitoring routes with dependency injection
func Routes() route.RoutesWithDependencies {
	return route.RoutesWithDependencies{
		route.CreateRoute(route.HTTPLivenessCheckPath, Liveness),
		route.CreateRoute(route.HTTPReadinessCheckPath, Readiness),
	}
}

// Liveness handles liveness checks
func Liveness(server *route.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := server.GetConfig()

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
				Version:     cfg.AppVersion,
				Environment: cfg.RuntimeEnv,
				Checks:      checks.Pretty(),
			},
		)
	}
}

// Readiness handles readiness checks
func Readiness(server *route.Server) http.HandlerFunc {
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
				Checks:    checks.Pretty(),
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

// GetLivenessStatus return liveness status
func GetLivenessStatus(checks CheckResults) Status {
	return StatusOK
}

// CheckLiveness performs liveness checks and returns a map of results
// It always returns [StatusOK] as responding to liveness probe via HTTP means that the app is alive
func CheckLiveness() CheckResults {
	return CheckResults{
		"http": StatusOK,
	}
}

// GetLivenessStatus return readiness status
func GetReadinessStatus(checks CheckResults) Status {
	return StatusOK
}

// CheckReadiness performs services checks and returns a map of results
func CheckReadiness() CheckResults {
	// TODO integrate with db / object storage clients
	return CheckResults{
		"database": StatusOK,
		"cache":    StatusOK,
	}
}
