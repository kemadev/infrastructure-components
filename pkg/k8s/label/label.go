package label

// OrgNs is the namespace for the Kema organization labels.
const OrgNs = "kema.dev"

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

// // DefaultLabels returns a set of default labels for the application instance as per Kubernetes convention,
// // see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
// func DefaultLabels(
// 	appName pulumi.StringInput,
// 	appInstance pulumi.StringInput,
// 	appVersion pulumi.StringInput,
// 	appComponent pulumi.StringInput,
// 	appNamespace pulumi.StringInput,
// ) pulumi.StringMap {
// 	return pulumi.StringMap{
// 		LabelAppNameKey:      appName,
// 		LabelAppInstanceKey:  appInstance,
// 		LabelAppVersionKey:   appVersion,
// 		LabelAppComponentKey: appComponent,
// 		LabelAppNamespaceKey: appNamespace,
// 		LabelAppMangedByKey:  pulumi.String("pulumi"),
// 	}
// }

// // DefaultLabels returns a set of default labels for the application instance as per Kubernetes convention,
// // see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
// func DefaultSelector(
// 	appInstance pulumi.StringInput,
// 	defaultLabels pulumi.StringMap,
// ) pulumi.StringMap {
// 	return pulumi.StringMap{
// 		LabelAppInstanceKey: defaultLabels[LabelAppInstanceKey],
// 	}
// }

// Topology labels
const (
	// LabelTopologyRegionKey is the label key for the region of the node, see https://kubernetes.io/docs/reference/labels-annotations-taints/#topologykubernetesioregion.
	LabelTopologyRegionKey = "topology.kubernetes.io/region"
	// LabelTopologyZoneKey is the label key for the zone of the node, see https://kubernetes.io/docs/reference/labels-annotations-taints/#topologykubernetesiozone.
	LabelTopologyZoneKey = "topology.kubernetes.io/zone"
	// LabelTopologyDatacenterKey is the label key for the datacenter hosting the node.
	LabelTopologyDatacenterKey = "topology." + OrgNs + "/dc"
	// LabelTopologyDatacenterZoneKey is the label key for the datacenter zone hosting the node.
	LabelTopologyDatacenterZoneKey = "topology." + OrgNs + "/dc-zone"
	// LabelTopologyDatacenterAisleKey is the label key for the datacenter aisle hosting the node.
	LabelTopologyDatacenterAisleKey = "topology." + OrgNs + "/dc-aisle"
	// LabelTopologyDatacenterRackKey is the label key for the datacenter rack hosting the node.
	LabelTopologyDatacenterRackKey = "topology." + OrgNs + "/dc-rack"
	// NodeHostnameLabelKey is the label key for the hostname of the node, see https://kubernetes.io/docs/reference/labels-annotations-taints/#kubernetesiohostname.
	LabelTopologyHostnameKey = "kubernetes.io/hostname"
)

// Node roles labels, see https://kubernetes.io/docs/reference/labels-annotations-taints/#node-role-kubernetes-io
const (
	// NodeRoleControlPlaneValue is the label key for control plane nodes.
	NodeRoleControlPlaneLabelKey = "node-role.kubernetes.io/control-plane"
	// NodeRoleControlPlane is the label value for control plane nodes. Empty value, following Kubernetes convention, see https://kubernetes.io/docs/reference/labels-annotations-taints/#node-role-kubernetes-io-control-plane.
	NodeRoleControlPlane = ""

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
	NodeRoleGPULabelKey = "node-role.kubernetes.io/gpu"
	// NodeRoleGenericGPU is the label value for generic GPU nodes.
	NodeRoleGenericGPU = "generic"
)

// Node labels for operating system and architecture
const (
	// NodeOSLabelKey is the label key for the operating system of the node.
	NodeOSLabelKey = "kubernetes.io/os"
	// NodeOSLinux is the label value for Linux operating system.
	NodeOSLinux = "linux"
	// NodeOSWindows is the label value for Windows operating system.
	NodeOSWindows = "windows"

	// NodeArchLabelKey is the label key for the architecture of the node.
	NodeArchLabelKey = "kubernetes.io/arch"
	// NodeArchAMD64 is the label value for AMD64 architecture.
	NodeArchAMD64 = "amd64"
	// NodeArchARM64 is the label value for ARM64 architecture.
	NodeArchARM64 = "arm64"
)

// Labels for shared gateway access, enabling the usage of the shared gateway
const (
	SharedGatewayAccessLabelKey   = "shared-gateway-access"
	SharedGatewayAccessLabelValue = "true"
)

type Taint struct {
	Key    string
	Value  string
	Effect string
}

// Node classic taints
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

// Node custom taints
const (
	// NodeTaintGPUKey is the taint key for nodes that are dedicated to GPU workloads.
	NodeTaintGPUKey = "nodepurpose." + OrgNs + "/gpu"
	// NodeTaintAcceleratedKey is the taint key for nodes that are dedicated to accelerated workloads (e.g. FPGA, TPU).
	NodeTaintAcceleratedKey = "nodepurpose." + OrgNs + "/accelerated"
	// NodeTaintCPUIntensiveKey is the taint key for nodes that are specialized for cpu-intensive workloads.
	NodeTaintCPUIntensiveKey = "nodepurpose." + OrgNs + "/high-cpu"
	// NodeTaintMemoryIntensiveKey is the taint key for nodes that are specialized for memory-intensive workloads.
	NodeTaintMemoryIntensiveKey = "nodepurpose." + OrgNs + "/high-memory"
	// NodeTaintStorageIntensiveKey is the taint key for nodes that are specialized for disk-intensive workloads.
	NodeTaintStorageIntensiveKey = "nodepurpose." + OrgNs + "/high-storage"
	// NodeTaintNetworkIntensiveKey is the taint key for nodes that are specialized for network-intensive workloads.
	NodeTaintNetworkIntensiveKey = "nodepurpose." + OrgNs + "/high-network"
)

// Taints effects
const (
	// TaintEffectNoSchedule is the taint effect for nodes that should not schedule any pods.
	TaintEffectNoSchedule = "NoSchedule"
	// TaintEffectPreferNoSchedule is the taint effect for nodes that should prefer not to schedule any pods.
	TaintEffectPreferNoSchedule = "PreferNoSchedule"
	// TaintEffectNoExecute is the taint effect for nodes that should not execute any pods.
	TaintEffectNoExecute = "NoExecute"
)
