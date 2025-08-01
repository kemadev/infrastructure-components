package cni

import (
	"fmt"

	"github.com/kemadev/infrastructure-components/pkg/k8s/priorityclass"
	"github.com/kemadev/infrastructure-components/pkg/k8s/pulumilabel"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	Namespace = "cilium"
)

// DeployCNI deploys the Cilium CNI using Helm, returning the corresponding Release object and an error if any.
func DeployCNI(
	ctx *pulumi.Context,
	gwapiCrd *yamlv2.ConfigFile,
	clusterName string,
) (*helm.Release, error) {
	const cniName = "cilium"
	// TODO add renovate tracking
	const cniVersion = "1.17.4"

	clusterNativeRoutingCIDR, err := RandomIPv6ULARoutingPrefix(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random IPv6 ULA: %w", err)
	}

	sharedLabels := pulumilabel.DefaultLabels(
		pulumi.String(cniName),
		pulumi.String(cniName),
		pulumi.String(cniVersion),
		pulumi.String("cni"),
		pulumi.String("network"),
	)

	ns, err := corev1.NewNamespace(ctx, Namespace, &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(Namespace),
			Namespace: pulumi.String(Namespace),
			Labels:    sharedLabels,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace %s: %w", cniName, err)
	}

	release, err := helm.NewRelease(ctx, cniName, &helm.ReleaseArgs{
		Name:        pulumi.String(cniName),
		Description: pulumi.String("Pretty much all the networking stuff"),
		Namespace:   ns.Metadata.Name(),
		Timeout:     pulumi.Int(600),
		RepositoryOpts: &helm.RepositoryOptsArgs{
			Repo: pulumi.String("https://helm.cilium.io/"),
		},
		Chart:   pulumi.String(cniName),
		Version: pulumi.String(cniVersion),
		Values: pulumi.All(
			clusterNativeRoutingCIDR,
		).ApplyT(func(cidr interface{}) pulumi.Map {
			nativeRoutingSubnet := cidr.(string)
			return pulumi.Map{
				// Add labels to all resources
				"commonLabels": sharedLabels,
				"image": pulumi.Map{
					// Don't pull if image already present
					"pullPolicy": pulumi.String("IfNotPresent"),
				},
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
					// TODO Force pod-to-pod encrpytion in all case, see https://docs.cilium.io/en/stable/security/network/encryption/#egress-traffic-to-not-yet-discovered-remote-endpoints-may-be-unencrypted (IPv6 not supported)
					// "strictMode":     pulumi.String("enabled"),
				},
				"externalIPs": pulumi.Map{
					// Enable ExternalIPs, see https://docs.cilium.io/en/stable/network/kubernetes/external-ips/
					"enabled": pulumi.Bool(true),
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
				},
				"hubble": pulumi.Map{
					// Enable Hubble
					"enabled": pulumi.Bool(true),
					"relay": pulumi.Map{
						// Enable Hubble relay
						"enabled": pulumi.Bool(true),
						// Rollout pods on ConfigMap change
						"rollOutPods": pulumi.Bool(true),
						// Set Hubble as moderate priority
						"priorityClassName": pulumi.String(priorityclass.PriorityClassModerate),
						"prometheus": pulumi.Map{
							// Expose Hubble relay metrics
							"enabled": pulumi.Bool(true),
						},
					},
					"ui": pulumi.Map{
						// Enable Hubble UI
						"enabled": pulumi.Bool(true),
						// Rollout pods on ConfigMap change
						"rollOutPods": pulumi.Bool(true),
						// Set Hubble as moderate priority
						"priorityClassName": pulumi.String(priorityclass.PriorityClassModerate),
						"livenessProbe": pulumi.Map{
							// Enable Hubble UI liveness probe
							"enabled": pulumi.Bool(true),
						},
						"readinessProbe": pulumi.Map{
							// Enable Hubble UI readiness probe
							"enabled": pulumi.Bool(true),
						},
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
				"nodePort": pulumi.Map{
					// Enable NoodePort, required for Gateway API Support
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
						// Enable Envoy structured logging, see https://www.envoyproxy.io/docs/envoy/latest/operations/cli#cmdoption-log-format & https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/bootstrap/v3/bootstrap.proto#envoy-v3-api-field-config-bootstrap-v3-bootstrap-applicationlogconfig-logformat-json-format
						// Inspired from OpenTelemetry format
						"format_json": pulumi.Map{
							"Timestamp":                       pulumi.String("%Y-%m-%dT%T.%e%z"),
							"SeverityText":                    pulumi.String("%l"),
							"Resource":                        pulumi.String("%n"),
							"Body":                            pulumi.String("%j"),
							string(semconv.CodeFilepathKey):   pulumi.String("%g"),
							string(semconv.CodeLineNumberKey): pulumi.String("%#"),
							string(semconv.CodeFunctionKey):   pulumi.String("%!"),
							string(semconv.ThreadIDKey):       pulumi.String("%t"),
						},
						"format": nil,
					},
				},
				"loadBalancer": pulumi.Map{
					"l7": pulumi.Map{
						// Use Envoy as L7 load balancer
						"backend": pulumi.String("envoy"),
					},
					// Use native mode XDP acceleration on devices that support it, see https://docs.cilium.io/en/stable/operations/performance/tuning/#xdp-acceleration
					"acceleration": pulumi.String("best-effort"),
					// Use hybrid DSR / SNAT, see https://docs.cilium.io/en/stable/network/kubernetes/kubeproxy-free/#hybrid-dsr-and-snat-mode
					"mode": pulumi.String("hybrid"),
					// Use Maglev consistent hashing, see https://docs.cilium.io/en/stable/network/kubernetes/kubeproxy-free/#maglev-consistent-hashing
					"algorithm": pulumi.String("maglev"),
				},
				// TODO enable mTLS
				// "authentication": pulumi.Map{
				// 	"mutual": pulumi.Map{
				// 		"spire": pulumi.Map{
				// 			// Enable SPIRE integration for mTLS, see https://docs.cilium.io/en/stable/security/network/encryption-wireguard/
				// 			"enabled": pulumi.Bool(true),
				// 			"install": pulumi.Map{
				// 				"server": pulumi.Map{
				// 					"ca": pulumi.Map{
				// 						// Set CA key algorithm, see https://spiffe.io/docs/latest/deploying/spire_server/#server-configuration-file
				// 						"keyType": pulumi.String("ec-p384"),
				// 					},
				// 				},
				// 			},
				// 		},
				// 	},
				// },
				"l2announcements": pulumi.Map{
					// Enable L2 announcements (see https://docs.cilium.io/en/stable/network/l2-announcements/), enabling LB IPAM, see https://docs.cilium.io/en/stable/network/lb-ipam/
					"enabled": pulumi.Bool(true),
				},
				"bpf": pulumi.Map{
					// Enable masquerading, see https://docs.cilium.io/en/stable/network/concepts/masquerading/
					// "masquerade": pulumi.Bool(true),
					// Mode for Pod devices for the core datapath
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
				// Use packet forwarding instead of encapsulation, see https://docs.cilium.io/en/stable/network/concepts/routing/#native-routing
				"routingMode": pulumi.String("native"),
				// Load routes in Linux kernel, see https://docs.cilium.io/en/stable/network/concepts/routing/#native-routing
				"autoDirectNodeRoutes": pulumi.Bool(true),
				"ipv4": pulumi.Map{
					// Disable IPv4
					"enabled": pulumi.Bool(false),
				},
				// "ipv4NativeRoutingCIDR": pulumi.String(nativeIPv4CIDR.String()),
				"ipv6": pulumi.Map{
					// Enable IPv6
					"enabled": pulumi.Bool(true),
				},
				"k8s": pulumi.Map{
					// Wait for PodCIDR allocation
					"requireIPv6PodCIDR": pulumi.Bool(true),
				},
				// Set cluster network CIDR, see https://docs.cilium.io/en/stable/network/concepts/routing/#native-routing
				"ipv6NativeRoutingCIDR": pulumi.String(nativeRoutingSubnet + "::/64"),
				"ipam": pulumi.Map{
					"operator": pulumi.Map{
						// Use cilium managed native routing
						"clusterPoolIPv6PodCIDRList": pulumi.StringArray{
							pulumi.String(nativeRoutingSubnet + "::/104"),
						},
					},
				},
			}
		}).(pulumi.MapOutput),
	}, pulumi.DependsOn([]pulumi.Resource{gwapiCrd}))
	if err != nil {
		return nil, fmt.Errorf("failed to deploy cni: %w", err)
	}

	return release, nil
}

// RandomIPv6ULARoutingPrefix generates a random IPv6 Unique Local Address (ULA) routing prefix, 64 bits masked.
func RandomIPv6ULARoutingPrefix(ctx *pulumi.Context) (pulumi.StringInput, error) {
	ula, err := random.NewRandomId(ctx, "ipv6-ula", &random.RandomIdArgs{
		ByteLength: pulumi.Int(7),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate random IPv6 ULA: %w", err)
	}

	ip := ula.Hex.ApplyT(func(hex string) (string, error) {
		b := make([]byte, 7)
		_, err := fmt.Sscanf(hex, "%16x", &b)
		if err != nil {
			return "", fmt.Errorf("failed to parse hex: %w", err)
		}
		buf := make([]byte, 16)
		buf[0] = 0xfd
		copy(buf[1:7], b)
		return fmt.Sprintf("%x:%x:%x:%x", buf[0:2], buf[2:4], buf[4:6], buf[6:8]), nil
	}).(pulumi.StringOutput)

	return ip, nil
}
