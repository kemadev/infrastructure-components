package gwapi

import (
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DeployGatewayAPICRDs(ctx *pulumi.Context) (*yamlv2.ConfigFile, error) {
	crd, err := yamlv2.NewConfigFile(ctx, "gateway-api-crds", &yamlv2.ConfigFileArgs{
		// TODO add renovate tracking
		File: pulumi.String(
			"https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/experimental-install.yaml",
		),
	})
	if err != nil {
		return nil, err
	}
	return crd, err
}
