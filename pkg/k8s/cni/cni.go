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
	nativeIPv4CIDR string,
	nativeIPv6CIDR string,
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
				// Don't pull if image already present
				"pullPolicy": pulumi.String("IfNotPresent"),
			},
			// kind specific, permit initial operator deployment, use a nillable value to avoid including the Helm value when it is nil
			"k8sServiceHost": func() pulumi.StringOutput {
				if ctx.Stack() != "dev" {
					return pulumi.StringOutput{}
				}
				return pulumi.String(clusterName + "-control-plane").ToStringOutput()
			}(),
			// kind specific, permit initial operator deployment, use a nillable value to avoid including the Helm value when it is nil
			"k8sServicePort": *func() *pulumi.Int {
				if ctx.Stack() != "dev" {
					return nil
				}
				res := pulumi.Int(6443)
				return &res
			}(),
			// Replace kube-proxy
			"kubeProxyReplacement": pulumi.Bool(true),
			// Enable L7 Gateway API capabilities
			"l7Proxy": pulumi.Bool(true),
			"encryption": pulumi.Map{
				// Enable transparent pod-to-pod encryption
				"enabled": pulumi.Bool(true),
				// Use WireGuard as encryption method
				"type": pulumi.String("wireguard"),
				// Encrypt pure node-to-node traffic
				"nodeEncryption": pulumi.Bool(true),
				// Force pod-to-pod encrpytion in all case, see https://docs.cilium.io/en/stable/security/network/encryption/#egress-traffic-to-not-yet-discovered-remote-endpoints-may-be-unencrypted
				// "strictMode":     pulumi.String("enabled"),
			},
			"gatewayAPI": pulumi.Map{
				// Enable cilium Gateway API
				"enabled": pulumi.Bool(true),
				"gatewayClass": pulumi.Map{
					// Create Cilium's GatewayClass
					"create": pulumi.String("true"),
				},
				// Enable ALPN
				"enableAlpn": pulumi.Bool(true),
				// Enable appProtocol, see https://kubernetes.io/docs/concepts/services-networking/service/#application-protocol
				"enableAppProtocol": pulumi.Bool(true),
				"hostNetwork": pulumi.Map{
					// Allow exposing Gateway API Gateway on host network
					"enabled": pulumi.Bool(true),
				},
			},
			"hubble": pulumi.Map{
				// Enable Hubble
				"enabled": pulumi.Bool(true),
				"relay": pulumi.Map{
					// Enable Hubble relay
					"enabled": pulumi.Bool(true),
					// Rollout pods on ConfigMap change
					"rollOutPods": pulumi.Bool(true),
				},
				"ui": pulumi.Map{
					// Enable Hubble UI
					"enabled": pulumi.Bool(true),
					// Rollout pods on ConfigMap change
					"rollOutPods": pulumi.Bool(true),
				},
				"metrics": pulumi.Map{
					// Expose Hubble metrics
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
					// Also expose as OpenMetrics format
					"enableOpenMetrics": pulumi.Bool(true),
				},
			},
			"prometheus": pulumi.Map{
				// Expose cilium-envoy metrics
				"enabled": pulumi.Bool(true),
			},
			"operator": pulumi.Map{
				"prometheus": pulumi.Map{
					// Expose cilium-operator metrics
					"enabled": pulumi.Bool(true),
				},
				// Rollout pods on ConfigMap change
				"rollOutPods": pulumi.Bool(true),
			},
			// Rollout pods on ConfigMap change
			"rollOutCiliumPods": pulumi.Bool(true),
			"envoyConfig": pulumi.Map{
				// Enable CiliumEnvoyConfig CRD
				"enabled": pulumi.Bool(true),
			},
			"envoy": pulumi.Map{
				// Rollout pods on ConfigMap change
				"rollOutPods": pulumi.Bool(true),
				"prometheus": pulumi.Map{
					// Expose envoy metrics
					"enabled": pulumi.Bool(true),
				},
				"log": pulumi.Map{
					"format": nil,
					"format_json": pulumi.String(
						// Enable Envoy structured logging, see https://www.envoyproxy.io/docs/envoy/latest/operations/cli#cmdoption-log-format & https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/bootstrap/v3/bootstrap.proto#envoy-v3-api-field-config-bootstrap-v3-bootstrap-applicationlogconfig-logformat-json-format
						`{"date":"%Y-%m-%dT%T.%e","level":"%l","logger":"%n","message":"%j","source_line":"%@","thread_id":"%t"}`,
					),
				},
			},
			"loadBalancer": pulumi.Map{
				"l7": pulumi.Map{
					// Use Envoy as L7 load balancer
					"backend": pulumi.String("envoy"),
				},
				// Use native mode XDP acceleration on devices that support it, see https://docs.cilium.io/en/stable/operations/performance/tuning/#xdp-acceleration
				"acceleration": pulumi.String("best-effort"),
			},
			"authentication": pulumi.Map{
				"mutual": pulumi.Map{
					"spire": pulumi.Map{
						// Enable SPIRE integration for mTLS, see https://docs.cilium.io/en/stable/security/network/encryption-wireguard/
						"enabled": pulumi.Bool(true),
						"install": pulumi.Map{
							"server": pulumi.Map{
								"ca": pulumi.Map{
									// Set CA key algorithm, see https://spiffe.io/docs/latest/deploying/spire_server/#server-configuration-file
									"keyType": pulumi.String("ec-p384"),
								},
							},
						},
					},
				},
			},
			"l2announcements": pulumi.Map{
				// Enable L2 announcements (see https://docs.cilium.io/en/stable/network/l2-announcements/), enabling LB IPAM, see https://docs.cilium.io/en/stable/network/lb-ipam/
				"enabled": pulumi.Bool(true),
			},
			"ipv6": pulumi.Map{
				// Enable IPv6
				"enabled": pulumi.Bool(true),
			},
			"endpointRoutes": pulumi.Map{
				// Enable use of per endpoint routes instead of routing via cilium_host interface
				"enabled": pulumi.Bool(true),
			},
			"bpf": pulumi.Map{
				// Enable masquerading, see https://docs.cilium.io/en/stable/network/concepts/masquerading/
				"masquerade": pulumi.Bool(true),
				// Mode for Pod devices for the core datapath (veth, netkit, netkit-l2)
				"dataPathMode": pulumi.String("netkit"),
				// Enables pre-allocation of eBPF map values
				"preallocateMaps": pulumi.Bool(true),
				// Enable eBPF-based TPROXY
				"tproxy": pulumi.Bool(true),
			},
			"bandwidthManager": pulumi.Map{
				// Enable Cilium’s bandwidth manager, see https://docs.cilium.io/en/stable/network/kubernetes/bandwidth-manager/
				"enabled": pulumi.Bool(true),
				// Enable BBR congestion control, see https://docs.cilium.io/en/stable/network/kubernetes/bandwidth-manager/#bbr-for-pods
				"bbr": pulumi.Bool(true),
			},
			// Enable local redirect, see https://docs.cilium.io/en/stable/network/kubernetes/local-redirect-policy/
			"localRedirectPolicy": pulumi.Bool(true),
			// Enable synchronizing Kubernetes EndpointSlice
			"ciliumEndpointSlice": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
			"hostFirewall": pulumi.Map{
				// Enable cilium host firewall
				"enabled": pulumi.Bool(true),
			},
			"maglev": pulumi.Map{
				// Set Maglev table size, see https://docs.cilium.io/en/latest/network/kubernetes/kubeproxy-free/#maglev-consistent-hashing
				"tableSize": pulumi.Int(16381),
			},
			"ipam": pulumi.Map{
				// Let cilium assign per-node PodCIDRs, see https://docs.cilium.io/en/stable/network/concepts/ipam/cluster-pool/
				"mode": pulumi.String("cluster-pool"),
			},
			// Use packet forwarding instead of encapsulation, see https://docs.cilium.io/en/stable/network/concepts/routing/#native-routing
			"routingMode": pulumi.String("native"),
			// Set cluster network CIDR, see https://docs.cilium.io/en/stable/network/concepts/routing/#native-routing
			"ipv4NativeRoutingCIDR": pulumi.String(nativeIPv4CIDR),
			"ipv6NativeRoutingCIDR": pulumi.String(nativeIPv6CIDR),
			// Load routes in Linux kernel, see https://docs.cilium.io/en/stable/network/concepts/routing/#native-routing
			"autoDirectNodeRoutes": pulumi.Bool(true),
			"nat46x64Gateway":      pulumi.Map{
				// TODO Enable NAT gateway, see https://isovalent.com/blog/post/cilium-release-112/#nat46-nat64
				// "enabled": pulumi.Bool(true),
			},
			"ipv4": pulumi.Map{
				// TODO Disable IPv4
				// TODO Also disable in kind
				// "enabled": pulumi.Bool(false),
			},
			// TODO https://medium.com/@nahelou.j/play-with-cilium-native-routing-in-kind-cluster-5a9e586a81ca
		},
	}, pulumi.DependsOn([]pulumi.Resource{gwapiCrd}))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy cni: %w", err)
	}

	return release, nil
}
