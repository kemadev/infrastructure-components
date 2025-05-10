package kind

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v2"
)

func CreateKindCluster(ctx *pulumi.Context, kindConfigPath string) (*local.Command, error) {
	if _, err := os.Stat(kindConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("kind-config-path does not exist: %s", kindConfigPath)
	}
	content, err := os.ReadFile(kindConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read kind-config-path: %s", err)
	}
	var contentMap map[string]any
	err = yaml.Unmarshal(content, &contentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal kind-config-path: %s", err)
	}
	clusterName, ok := contentMap["name"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to get cluster name from kind-config-path")
	}
	if clusterName == "" {
		return nil, fmt.Errorf("cluster name is empty in kind-config-path")
	}
	cluster, err := local.NewCommand(ctx, "cluster", &local.CommandArgs{
		Create: pulumi.String("kind create cluster --config " + kindConfigPath),
		Delete: pulumi.String("kind delete cluster --name " + clusterName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %s", err)
	}
	ctx.Export("clusterName", pulumi.String(clusterName))
	return cluster, nil
}
