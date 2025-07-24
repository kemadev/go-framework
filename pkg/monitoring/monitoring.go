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

type LivenessResponse struct {
	Timestamp time.Time         `json:"timestamp"`
	Started   bool              `json:"started"`
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Checks    map[string]string `json:"checks"`
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

func Routes() []route.Route {
	var routes []route.Route

	liveness := route.Route{
		Pattern:     route.HTTPLivenessCheckPath,
		HandlerFunc: Liveness,
	}
	routes = append(routes, liveness)

	readiness := route.Route{
		Pattern:     route.HTTPReadinessCheckPath,
		HandlerFunc: Readiness,
	}
	routes = append(routes, readiness)

	return routes
}

func Liveness(w http.ResponseWriter, r *http.Request) {
	kclient := khttp.ClientInfo{
		Ctx:    context.Background(),
		Writer: w,
	}

	khttp.SendJSONResponse(
		kclient,
		200,
		LivenessResponse{
			Timestamp: time.Now().UTC(),
			Started:   true,
			Status:    GetLivenessStatus(),
			Version:   config.AppVersion(),
		},
	)
}

func Readiness(w http.ResponseWriter, r *http.Request) {
	kclient := khttp.ClientInfo{
		Ctx:    context.Background(),
		Writer: w,
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	khttp.SendJSONResponse(
		kclient,
		200,
		ReadinessResponse{
			Timestamp: time.Now().UTC(),
			Ready:     GetReadinessStatus(),
			Services:  CheckServicesStatus(),
			RuntimeMetrics: RuntimeMetrics{
				Memory: MemoryMetrics{
					UsedBytes:    float64(m.Alloc),
					MaxBytes:     float64(m.Sys),
					UsagePercent: float64(m.Alloc) * float64(m.Sys) / 100,
					GCRuns:       m.NumGC,
				},
				CPU: CPUMetrics{
					Goroutines: runtime.NumGoroutine(),
				},
			},
		},
	)
}
