package main

import (
	"fmt"

	"github.com/kemadev/infrastructure-components/pkg/k8s/cni"
	"github.com/kemadev/infrastructure-components/pkg/k8s/gwapicrds"
	"github.com/kemadev/infrastructure-components/pkg/k8s/kind"
	"github.com/kemadev/infrastructure-components/pkg/k8s/priorityclass"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Run for dev only, other clusters are created in ad-hoc repo
		if ctx.Stack() != "dev" {
			return nil
		}
		clusterName, err := kind.GetClusterName(ctx, "../../config/kind/kind-config.yaml")
		if err != nil {
			return fmt.Errorf("failed to get cluster name: %w", err)
		}

		gwapiCrd, err := gwapicrds.DeployGatewayAPICRDs(ctx)
		if err != nil {
			return fmt.Errorf("failed to deploy gateway api crds: %w", err)
		}

		_, err = cni.DeployCNI(
			ctx,
			gwapiCrd,
			clusterName,
		)
		if err != nil {
			return fmt.Errorf("failed to deploy cni: %w", err)
		}

		err = priorityclass.CreateDefaultPriorityClasses(ctx)
		if err != nil {
			return fmt.Errorf("failed to deploy priority classes: %w", err)
		}
		return nil
	})
}
