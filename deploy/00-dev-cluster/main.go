package main

import (
	"github.com/kemadev/infrastructure-components/pkg/k8s/kind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Run for dev only, other clusters are created in ad-hoc repo
		if ctx.Stack() != "dev" {
			return nil
		}
		cluster, err := kind.CreateKindCluster(ctx, "../../config/kind/", false)
		if err != nil {
			return err
		}
		_ = cluster
		return nil
	})
}
