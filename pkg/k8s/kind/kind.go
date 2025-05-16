package kind

import (
	"fmt"
	"os"

	clusterDef "github.com/kemadev/imds/pkg/hardware/cluster"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v2"
)

// GetClusterName reads the kind config file and returns the cluster name, and an error if any.
func GetClusterName(ctx *pulumi.Context, kindConfigPath string) (string, error) {
	if _, err := os.Stat(kindConfigPath); os.IsNotExist(err) {
		return "", fmt.Errorf("kind-config-path %s does not exist", kindConfigPath)
	}
	content, err := os.ReadFile(kindConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to read kind-config-path: %w", err)
	}
	var contentMap map[string]any
	err = yaml.Unmarshal(content, &contentMap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal kind-config-path: %w", err)
	}
	clusterName, ok := contentMap["name"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get cluster name from kind-config-path")
	}
	if clusterName == "" {
		return "", fmt.Errorf("cluster name is empty in kind-config-path")
	}
	return clusterName, nil
}

// addNodeLabels adds conventional labels to the nodes.
func addNodeLabels(ctx *pulumi.Context, clusterName string, cluster *local.Command) error {
	nodes := map[string]map[string]string{
		clusterName + "-control-plane": {
			clusterDef.NodeRoleControlPlaneLabelKey: clusterDef.NodeRoleControlPlaneLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-1",
		},
		clusterName + "-control-plane2": {
			clusterDef.NodeRoleControlPlaneLabelKey: clusterDef.NodeRoleControlPlaneLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-2",
		},
		clusterName + "-control-plane3": {
			clusterDef.NodeRoleControlPlaneLabelKey: clusterDef.NodeRoleControlPlaneLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-3",
		},

		clusterName + "-worker": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-1",
		},
		clusterName + "-worker2": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-1",
		},
		clusterName + "-worker3": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-1",
		},

		clusterName + "-worker4": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-2",
		},
		clusterName + "-worker5": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-2",
		},
		clusterName + "-worker6": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-2",
		},

		clusterName + "-worker5": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-3",
		},
		clusterName + "-worker6": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-3",
		},
		clusterName + "-worker7": {
			clusterDef.NodeRoleWorkerDefaultLabelKey: clusterDef.NodeRoleWorkerDefaultLabelValue,
			clusterDef.NodeRegionLabelKey:            "region-1",
			clusterDef.NodeZoneLabelKey:              "region-1-3",
		},
	}
	for name, node := range nodes {
		_, err := corev1.NewNodePatch(ctx, name+"-patch", &corev1.NodePatchArgs{
			Metadata: &metav1.ObjectMetaPatchArgs{
				Name: pulumi.String(name),
				Labels: func() pulumi.StringMap {
					labels := make(pulumi.StringMap)
					for k, v := range node {
						labels[k] = pulumi.String(v)
					}
					return labels
				}(),
			},
		}, pulumi.DependsOn([]pulumi.Resource{cluster}))
		if err != nil {
			return fmt.Errorf("failed to create node patch: %w", err)
		}
	}
	return nil
}

// CreateKindCluster creates a kind cluster using the provided kind config path, and returns a command object and an error if any.
func CreateKindCluster(ctx *pulumi.Context, kindConfigPath string) (*local.Command, error) {
	clusterName, err := GetClusterName(ctx, kindConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster name: %w", err)
	}
	cluster, err := local.NewCommand(ctx, "cluster", &local.CommandArgs{
		Create: pulumi.String("kind create cluster --config " + kindConfigPath),
		Delete: pulumi.String("kind delete cluster --name " + clusterName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}
	err = addNodeLabels(ctx, clusterName, cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to apply node labels: %w", err)
	}
	ctx.Export("clusterName", pulumi.String(clusterName))
	return cluster, nil
}
