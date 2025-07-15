package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kemadev/infrastructure-components/pkg/k8s/basichttpapp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := basichttpapp.DeployBasicHTTPApp(ctx, basichttpapp.AppParms{
			AppNamespace:        "changeme",
			AppComponent:        "changeme",
			BusinessUnitId:      "changeme",
			CustomerId:          "changeme",
			CostCenter:          "changeme",
			CostAllocationOwner: "changeme",
			OperationsOwner:     "changeme",
			Rpo:                 0 * time.Second,
			MonitoringUrl: url.URL{
				Scheme: "https",
				Host:   "changeme",
				Path:   "changeme",
			},
		})
		if err != nil {
			return fmt.Errorf("failed to deploy basic HTTP app: %w", err)
		}
		return nil
	})
}
