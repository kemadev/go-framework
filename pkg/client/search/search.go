package search

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/monitoring"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

func NewClient(conf config.SearchConfig, runtime config.Runtime) (*opensearchapi.Client, error) {
	clientAddresses := []string{}
	for _, addr := range conf.ClientAddress {
		clientAddresses = append(clientAddresses, addr.String())
	}

	baseTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	instrumentedTransport := otelhttp.NewTransport(
		baseTransport,
		otelhttp.WithSpanOptions(trace.WithAttributes(
			semconv.DBSystemNameOpenSearch,
		)),
	)

	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Transport:     instrumentedTransport,
			Addresses:     clientAddresses,
			Username:      conf.Username,
			Password:      conf.Password,
			EnableMetrics: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating search client: %w", err)
	}

	return client, nil
}

func Check(c *opensearchapi.Client) monitoring.StatusCheck {
	resp, err := c.Ping(context.Background(), &opensearchapi.PingReq{})
	if err != nil || resp.IsError() {
		return monitoring.StatusCheck{
			Status:  monitoring.StatusDown,
			Message: "ping failed",
		}
	}

	return monitoring.StatusCheck{
		Status:  monitoring.StatusOK,
		Message: monitoring.StatusOK.String(),
	}
}
