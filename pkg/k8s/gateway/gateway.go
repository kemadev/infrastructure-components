package gateway

import (
	"fmt"
	"net"

	"github.com/kemadev/infrastructure-components/pkg/k8s/label"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// SharedGatewayName is the name of the shared gateway.
	SharedGatewayName = "shared-gateway"
	// SharedGatewayNamespace is the namespace where the shared gateway resources are deployed.
	SharedGatewayNamespace = "shared-gateway"
)

// deployGatewayResources deploys the Gateway and LB-IPAM resources for all domains, creating setting
// up TLS termination and wildcard certificates for each domain.
func DeployGatewayResources(
	ctx *pulumi.Context,
	certIssuerName string,
	lbPoolCIDR net.IPNet,
	gatewayIPs []net.IP,
	domains []string,
) error {
	sharedLabels := label.DefaultLabels(
		pulumi.String("shared-gateway"),
		pulumi.String("shared-gateway"),
		pulumi.String("1"),
		pulumi.String("gateway"),
		pulumi.String("network"),
	)

	_, err := corev1.NewNamespace(ctx, SharedGatewayNamespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(SharedGatewayNamespace),
			Namespace: pulumi.String(SharedGatewayNamespace),
			Labels:    sharedLabels,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create namespace %s: %w", SharedGatewayName, err)
	}

	_, err = yamlv2.NewConfigGroup(ctx, "lb-pool-1", &yamlv2.ConfigGroupArgs{
		Objs: pulumi.Array{
			pulumi.Map{
				"apiVersion": pulumi.String("cilium.io/v2alpha1"),
				"kind":       pulumi.String("CiliumLoadBalancerIPPool"),
				"metadata": pulumi.Map{
					"name":      pulumi.String("lb-pool-1"),
					"namespace": pulumi.String(SharedGatewayNamespace),
					"labels":    sharedLabels,
				},
				"spec": pulumi.Map{
					"blocks": pulumi.Array{
						pulumi.Map{
							"cidr": pulumi.String(lbPoolCIDR.String()),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to deploy CiliumLoadBalancerIPPool: %w", err)
	}

	_, err = yamlv2.NewConfigGroup(ctx, "announcement-policy-1", &yamlv2.ConfigGroupArgs{
		Objs: pulumi.Array{
			pulumi.Map{
				"apiVersion": pulumi.String("cilium.io/v2alpha1"),
				"kind":       pulumi.String("CiliumL2AnnouncementPolicy"),
				"metadata": pulumi.Map{
					"name":      pulumi.String("announcement-policy-1"),
					"namespace": pulumi.String(SharedGatewayNamespace),
					"labels":    sharedLabels,
				},
				"spec": pulumi.Map{
					"externalIPs":     pulumi.Bool(true),
					"loadBalancerIPs": pulumi.Bool(true),
					"nodeSelector": pulumi.Map{
						"matchExpressions": pulumi.Array{
							pulumi.Map{
								"key":      pulumi.String("node-role.kubernetes.io/control-plane"),
								"operator": pulumi.String("DoesNotExist"),
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to deploy CiliumL2AnnouncementPolicy: %w", err)
	}

	_, err = yamlv2.NewConfigGroup(ctx, "Gateway", &yamlv2.ConfigGroupArgs{
		Objs: pulumi.Array{
			pulumi.Map{
				"apiVersion": pulumi.String("gateway.networking.k8s.io/v1"),
				"kind":       pulumi.String("Gateway"),
				"metadata": pulumi.Map{
					"name":      pulumi.String(SharedGatewayName),
					"namespace": pulumi.String(SharedGatewayNamespace),
					"labels":    sharedLabels,
					"annotations": pulumi.Map{
						// Integrate with cert-manager
						"cert-manager.io/issuer": pulumi.String(certIssuerName),
					},
				},
				"spec": pulumi.Map{
					"addresses": func() pulumi.ArrayInput {
						if len(gatewayIPs) == 0 {
							return nil
						}
						addrs := make(pulumi.Array, len(gatewayIPs))
						for i, ip := range gatewayIPs {
							addrs[i] = pulumi.Map{
								"type":  pulumi.String("IPAddress"),
								"value": pulumi.String(ip.String()),
							}
						}
						return addrs
					}(),
					"gatewayClassName": pulumi.String("cilium"),
					"listeners": func() pulumi.ArrayInput {
						var listeners pulumi.Array = nil
						if ctx.Stack() == "dev" {
							// Let HTTP traffic through for dev
							listeners = append(listeners, pulumi.Map{
								"name":     pulumi.String("http-dev"),
								"port":     pulumi.Int(80),
								"protocol": pulumi.String("HTTP"),
								"allowedRoutes": pulumi.Map{
									"namespaces": pulumi.Map{
										"from": pulumi.String("Selector"),
										"selector": pulumi.Map{
											"matchLabels": pulumi.Map{
												label.SharedGatewayAccessLabelKey: pulumi.String(
													label.SharedGatewayAccessLabelValue,
												),
											},
										},
									},
								},
							})
						}
						for _, domain := range domains {
							l := pulumi.Map{
								"name":     pulumi.String(domain + "-wildcard"),
								"port":     pulumi.Int(443),
								"protocol": pulumi.String("HTTPS"),
								"hostname": pulumi.String("*." + domain),
								"tls": pulumi.Map{
									"mode": pulumi.String("Terminate"),
									"certificateRefs": pulumi.Array{
										pulumi.Map{
											"kind": pulumi.String("Secret"),
											"name": pulumi.String("wildcard-cert-" + domain),
										},
									},
								},
								"allowedRoutes": pulumi.Map{
									"namespaces": pulumi.Map{
										"from": pulumi.String("Selector"),
										"selector": pulumi.Map{
											"matchLabels": pulumi.Map{
												label.SharedGatewayAccessLabelKey: pulumi.String(
													label.SharedGatewayAccessLabelValue,
												),
											},
										},
									},
								},
							}
							listeners = append(listeners, l)
						}
						return listeners
					}(),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to deploy Gateway: %w", err)
	}

	return nil
}
