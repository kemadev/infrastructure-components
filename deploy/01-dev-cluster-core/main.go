package main

import (
	"fmt"
	"net"

	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"vcs.kema.run/kema/infrastructure-components/pkg/k8s/cni"
	"vcs.kema.run/kema/infrastructure-components/pkg/k8s/gwapi"
	"vcs.kema.run/kema/infrastructure-components/pkg/k8s/kind"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Run for dev only, other clusters are created in ad-hoc repo
		if ctx.Stack() != "dev" {
			return nil
		}
		network, err := docker.LookupNetwork(ctx, &docker.LookupNetworkArgs{
			Name: "kind",
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to get docker network: %w", err)
		}
		var ipv4Subnet, ipv6Subnet string
		for _, ipam := range network.IpamConfigs {
			ip, subnet, err := net.ParseCIDR(*ipam.Subnet)
			if err != nil {
				return fmt.Errorf("failed to parse docker network subnet: %w", err)
			}
			if ip.To4() != nil {
				ipv4Subnet = subnet.String()
			} else if ip.To16() != nil {
				ipv6Subnet = subnet.String()
			} else {
				return fmt.Errorf("failed to parse docker network ip")
			}
		}
		clusterName, err := kind.GetClusterName(ctx, "../../config/kind/kind-config.yaml")
		if err != nil {
			return fmt.Errorf("failed to get cluster name: %w", err)
		}
		gwapiCrd, err := gwapi.DeployGatewayAPICRDs(ctx)
		if err != nil {
			return fmt.Errorf("failed to deploy gateway api crds: %w", err)
		}
		cni, err := cni.DeployCNI(
			ctx,
			gwapiCrd,
			clusterName,
			ipv4Subnet,
			ipv6Subnet,
		)
		_ = cni
		return nil
	})
}
