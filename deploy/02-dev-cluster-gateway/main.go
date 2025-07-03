package main

import (
	"net"

	"github.com/kemadev/infrastructure-components/pkg/k8s/gateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Run for dev only, other clusters are created in ad-hoc repo
		if ctx.Stack() != "dev" {
			return nil
		}
		err := gateway.DeployGatewayResources(
			ctx,
			"",
			net.IPNet{
				IP:   net.IPv4(172, 18, 250, 0),
				Mask: net.CIDRMask(24, 32),
			},
			[]net.IP{
				net.IPv4(172, 18, 250, 10),
			},
			[]string{},
		)
		if err != nil {
			return err
		}
		return nil
	})
}
