package label

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// Application labels
const (
	// LabelAppNameKey is the label key for the name of the application, see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels.
	LabelAppNameKey = "app.kubernetes.io/name"
	// LabelAppInstanceKey is the label key for the name of the application instance, see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels.
	LabelAppInstanceKey = "app.kubernetes.io/instance"
	// LabelAppVersionKey is the label key for the version of the application, see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels.
	LabelAppVersionKey = "app.kubernetes.io/version"
	// LabelAppComponentKey is the label key for the component of the application, see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels.
	LabelAppComponentKey = "app.kubernetes.io/component"
	// LabelAppNamespaceKey is the label key for the namespace of the application, see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels.
	LabelAppNamespaceKey = "app.kubernetes.io/part-of"
	// LabelAppMangedByKey is the label key for the name of the application manager, see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels.
	LabelAppMangedByKey = "app.kubernetes.io/managed-by"
)

// DefaultLabels returns a set of default labels for the application instance as per Kubernetes convention,
// see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
func DefaultLabels(
	appName pulumi.StringInput,
	appInstance pulumi.StringInput,
	appVersion pulumi.StringInput,
	appComponent pulumi.StringInput,
	appNamespace pulumi.StringInput,
) pulumi.StringMap {
	return pulumi.StringMap{
		LabelAppNameKey:      appName,
		LabelAppInstanceKey:  appInstance,
		LabelAppVersionKey:   appVersion,
		LabelAppComponentKey: appComponent,
		LabelAppNamespaceKey: appNamespace,
		LabelAppMangedByKey:  pulumi.String("pulumi"),
	}
}

// DefaultLabels returns a set of default labels for the application instance as per Kubernetes convention,
// see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
func DefaultSelector(
	appInstance pulumi.StringInput,
	defaultLabels pulumi.StringMap,
) pulumi.StringMap {
	return pulumi.StringMap{
		LabelAppInstanceKey: defaultLabels[LabelAppInstanceKey],
	}
}

// Topology labels
const (
	// LabelTopologyRegionKey is the label key for the region of the node, see https://kubernetes.io/docs/reference/labels-annotations-taints/#topologykubernetesioregion.
	LabelTopologyRegionKey = "topology.kubernetes.io/region"
	// LabelTopologyZoneKey is the label key for the zone of the node, see https://kubernetes.io/docs/reference/labels-annotations-taints/#topologykubernetesiozone.
	LabelTopologyZoneKey = "topology.kubernetes.io/zone"
	// LabelTopologyDatacenterKey is the label key for the datacenter hosting the node.
	LabelTopologyDatacenterKey = "topology.kema.dev/dc"
	// LabelTopologyDatacenterZoneKey is the label key for the datacenter zone hosting the node.
	LabelTopologyDatacenterZoneKey = "topology.kema.dev/dc-zone"
	// LabelTopologyDatacenterAisleKey is the label key for the datacenter aisle hosting the node.
	LabelTopologyDatacenterAisleKey = "topology.kema.dev/dc-aisle"
	// LabelTopologyDatacenterRackKey is the label key for the datacenter rack hosting the node.
	LabelTopologyDatacenterRackKey = "topology.kema.dev/dc-rack"
	// NodeHostnameLabelKey is the label key for the hostname of the node, see https://kubernetes.io/docs/reference/labels-annotations-taints/#kubernetesiohostname.
	LabelTopologyHostnameKey = "kubernetes.io/hostname"
)

// Node roles labels, see https://kubernetes.io/docs/reference/labels-annotations-taints/#node-role-kubernetes-io
const (
	// NodeRoleControlPlaneValue is the label key for control plane nodes.
	NodeRoleControlPlaneLabelKey = "node-role.kubernetes.io/control-plane"
	// NodeRoleControlPlaneLabelValue is the label value for control plane nodes. Empty value, following Kubernetes convention, see https://kubernetes.io/docs/reference/labels-annotations-taints/#node-role-kubernetes-io-control-plane.
	NodeRoleControlPlaneLabelValue = ""

	// NodeRoleComputeLabelKey is the label key for compute nodes.
	NodeRoleComputeLabelKey = "node-role.kubernetes.io/compute"
	// NodeRoleComputeGeneric is the label value for generic compute nodes.
	NodeRoleComputeGeneric = "generic"

	// NodeRoleCPUIntensiveLabelKey is the label key for cpu intensive compute nodes.
	NodeRoleCPUIntensiveLabelKey = "node-role.kubernetes.io/compute-cpu-intensive"
	// NodeRoleComputeGeneric is the label value for generic cpu intensive compute nodes.
	NodeRoleComputeCPUIntensiveGeneric = "generic"

	// NodeRoleMemoryIntensiveLabelKey is the label key for memory intensive compute nodes.
	NodeRoleMemoryIntensiveLabelKey = "node-role.kubernetes.io/compute-memory-intensive"
	// NodeRoleComputeMemoryIntensiveGeneric is the label value for generic memory intensive compute nodes.
	NodeRoleComputeMemoryIntensiveGeneric = "generic"

	// NodeRoleNetworkIntensiveLabelKey is the label key for network intensive compute nodes.
	NodeRoleNetworkIntensiveLabelKey = "node-role.kubernetes.io/compute-network-intensive"
	// NodeRoleComputeNetworkIntensiveGeneric is the label value for generic network intensive compute nodes.
	NodeRoleComputeNetworkIntensiveGeneric = "generic"

	// NodeRoleStorageLabelKey is the label key for storage nodes.
	NodeRoleStorageLabelKey = "node-role.kubernetes.io/storage"
	// NodeRoleStorageGeneric is the label value for generic storage nodes.
	NodeRoleStorageGeneric = "generic"

	// NodeRoleAcceleratedLabelKey is the label key for accelerated nodes (e.g. FPGA, TPU).
	NodeRoleAcceleratedLabelKey = "node-role.kubernetes.io/accelerated"
	// NodeRoleAcceleratedGeneric is the label value for generic accelerated nodes.
	NodeRoleAcceleratedGeneric = "generic"

	// NodeRoleGPUValue is the label key for GPU nodes.
	NodeRoleGenericGPULabelKey = "node-role.kubernetes.io/gpu"
	// NodeRoleGenericGPULabelValue is the label value for generic GPU nodes.
	NodeRoleGenericGPULabelValue = "generic"
)

// Labels for shared gateway access, enabling the usage of the shared gateway
const (
	SharedGatewayAccessLabelKey   = "shared-gateway-access"
	SharedGatewayAccessLabelValue = "true"
)

// Node taints
const (
	// NodeTaintNotReadyKey is the taint key for nodes that are not ready.
	NodeTaintNotReadyKey = "node.kubernetes.io/not-ready"
	// NodeTaintUnreachableKey is the taint key for nodes that are unreachable.
	NodeTaintUnreachableKey = "node.kubernetes.io/unreachable"
	// NodeTaintDiskPressureKey is the taint key for nodes that are under disk pressure.
	NodeTaintDiskPressureKey = "node.kubernetes.io/disk-pressure"
	// NodeTaintMemoryPressureKey is the taint key for nodes that are under memory pressure.
	NodeTaintMemoryPressureKey = "node.kubernetes.io/memory-pressure"
	// NodeTaintPIDPressureKey is the taint key for nodes that are under PID pressure.
	NodeTaintPIDPressureKey = "node.kubernetes.io/pid-pressure"
	// NodeTaintUnschedulableKey is the taint key for nodes that are unschedulable.
	NodeTaintUnschedulableKey = "node.kubernetes.io/unschedulable"
	// NodeTaintNetworkUnavailableKey is the taint key for nodes that are under network pressure.
	NodeTaintNetworkUnavailableKey = "node.kubernetes.io/network-unavailable"
	// NodeTaintUninitializedKey is the taint key for nodes that are uninitialized.
	NodeTaintUninitializedKey = "node.cloudprovider.kubernetes.io/uninitialized"
	// NodeTaintControlPlaneKey is the taint key for control plane nodes.
	NodeTaintControlPlaneKey = NodeRoleControlPlaneLabelKey
)
