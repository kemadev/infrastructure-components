package cni

import (
	"fmt"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DeployCNI(
	ctx *pulumi.Context,
	gwapiCrd *yamlv2.ConfigFile,
	clusterName string,
) (*helm.Release, error) {
	const cniName = "cilium"

	const cniNsName = cniName
	ns, err := corev1.NewNamespace(ctx, cniNsName, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Namespace: pulumi.String(cniNsName),
			Labels: pulumi.StringMap{
				"app": pulumi.String(cniNsName),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace %s: %w", cniName, err)
	}

	release, err := helm.NewRelease(ctx, cniName, &helm.ReleaseArgs{
		Name:        pulumi.String(cniName),
		Description: pulumi.String("Pretty much all the networking stuff"),
		Namespace:   ns.Metadata.Name(),
		Timeout:     pulumi.Int(120),
		RepositoryOpts: &helm.RepositoryOptsArgs{
			Repo: pulumi.String("https://helm.cilium.io/"),
		},
		Chart: pulumi.String(cniName),
		// TODO add renovate tracking
		Version: pulumi.String("1.17.2"),
		Values: pulumi.Map{
			"image": pulumi.Map{
				"pullPolicy": pulumi.String("IfNotPresent"),
			},
			// Use a nillable value to avoid including the Helm value when it is nil
			"k8sServiceHost": func() pulumi.StringOutput {
				if ctx.Stack() != "dev" {
					return pulumi.StringOutput{}
				}
				return pulumi.String(clusterName + "-control-plane").ToStringOutput()
			}(),
			// Use a nillable value to avoid including the Helm value when it is nil
			"k8sServicePort": *func() *pulumi.Int {
				if ctx.Stack() != "dev" {
					return nil
				}
				res := pulumi.Int(6443)
				return &res
			}(),
			"kubeProxyReplacement": pulumi.Bool(true),
			"l7Proxy":              pulumi.Bool(true),
			"encryption": pulumi.Map{
				"enabled":        pulumi.Bool(true),
				"type":           pulumi.String("wireguard"),
				"nodeEncryption": pulumi.Bool(true),
			},
			"gatewayAPI": pulumi.Map{
				"enabled": pulumi.Bool(true),
				"gatewayClass": pulumi.Map{
					"create": pulumi.String("true"),
				},
				"enableAlpn":        pulumi.Bool(true),
				"enableAppProtocol": pulumi.Bool(true),
				// TODO
				// "hostNetwork": pulumi.Map{
				// 	"enabled": pulumi.Bool(true),
				// },
			},
			"hubble": pulumi.Map{
				"enabled": pulumi.Bool(true),
				"relay": pulumi.Map{
					"enabled":     pulumi.Bool(true),
					"rollOutPods": pulumi.Bool(true),
				},
				"ui": pulumi.Map{
					"enabled":     pulumi.Bool(true),
					"rollOutPods": pulumi.Bool(true),
				},
				"metrics": pulumi.Map{
					"enabled": pulumi.Array{
						pulumi.String("tcp"),
						pulumi.String("flow"),
						pulumi.String("port-distribution"),
						pulumi.String("icmp"),
						pulumi.String("dns:labelsContext=source_namespace,destination_namespace"),
						pulumi.String("drop:labelsContext=source_namespace,destination_namespace"),
						pulumi.String(
							"httpV2:exemplars=true;sourceContext=workload-name|pod-name|reserved-identity;destinationContext=workload-name|pod-name|reserved-identity;labelsContext=source_namespace,destination_namespace,traffic_direction",
						),
					},
					"enableOpenMetrics": pulumi.Bool(true),
				},
			},
			"prometheus": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"operator": pulumi.Map{
				"prometheus": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
				"rollOutPods": pulumi.Bool(true),
			},
			"ipam": pulumi.Map{
				"mode": pulumi.String("cluster-pool"),
			},
			"rollOutCiliumPods": pulumi.Bool(true),
			"envoyConfig": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"envoy": pulumi.Map{
				"rollOutPods": pulumi.Bool(true),
			},
			"loadBalancer": pulumi.Map{
				"l7": pulumi.Map{
					"backend": pulumi.String("envoy"),
				},
				"acceleration": pulumi.String("best-effort"),
			},
			"maglev": pulumi.Map{
				"tableSize": pulumi.Int(16381),
			},
			"authentication": pulumi.Map{
				"mutual": pulumi.Map{
					"spire": pulumi.Map{
						"enabled": pulumi.Bool(true),
						"install": pulumi.Map{
							"enabled": pulumi.Bool(true),
						},
					},
				},
			},
			"l2announcements": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"ipv6": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			// TODO enable those values and look for new ones
			// "routingMode": pulumi.String("native"),
			// "endpointRoutes": pulumi.Map{
			// 	"enabled": pulumi.Bool(true),
			// },
			// "bpf": pulumi.Map{
			// 	"masquerade": pulumi.Bool(true),
			// 	"dataPathMode": pulumi.String("netkit"),
			// 	"preallocateMaps": pulumi.Bool(true),
			// 	"tproxy": pulumi.Bool(true),
			// },
			// "bandwidthManager": pulumi.Map{
			// 	"enabled": pulumi.Bool(true),
			// 	"bbr": pulumi.Bool(true),
			// },
			// "autoDirectNodeRoutes": pulumi.Bool(true),
			// "localRedirectPolicy": pulumi.Bool(true),
		},
	}, pulumi.DependsOn([]pulumi.Resource{gwapiCrd}))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy cni: %w", err)
	}

	return release, nil
}
