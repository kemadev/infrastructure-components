package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"vcs.kema.run/kema/infrastructure-components/pkg/k8s/kind"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Run for dev only, other clusters are created in ad-hoc repo
		if ctx.Stack() != "dev" {
			return nil
		}
		cluster, err := kind.CreateKindCluster(ctx, "../../config/kind/kind-config.yaml")
		if err != nil {
			return err
		}
		_ = cluster
		return nil
	})
}
