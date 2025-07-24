package monitoring

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/kemadev/go-framework/pkg/config"
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

type ReadinessResponse struct {
	Timestamp      time.Time              `json:"timestamp"`
	Ready          bool                   `json:"ready"`
	RuntimeMetrics RuntimeMetrics         `json:"runtimeMetrics"`
	Services       map[string]interface{} `json:"services"`
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

		status := GetLivenessStatus()

		khttp.SendJSONResponse(
			kclient,
			status.HTTPCode(),
			LivenessResponse{
				Timestamp:   time.Now().UTC(),
				Started:     true,
				Status:      status.String(),
				Version:     cfg.AppVersion,
				Environment: cfg.RuntimeEnv,
				Checks:      GetLivenessChecks(),
			},
		)
	}
}

// Liveness handles readiness checks
func Readiness(server *route.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := server.GetConfig()

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

		status := GetReadinessStatus(cfg)

		khttp.SendJSONResponse(
			kclient,
			status.HTTPCode(),
			ReadinessResponse{
				Timestamp: time.Now().UTC(),
				Ready:     status.IsReady(),
				Services:  CheckServicesStatus(cfg),
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

func GetLivenessStatus() Status {
	return StatusOK
}

func GetLivenessChecks() map[string]string {
	return map[string]string{
		"server": "ok",
	}
}

func GetReadinessStatus(cfg *config.Config) Status {
	return StatusOK
}

func CheckServicesStatus(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"database": "connected",
		"cache":    "connected",
	}
}
