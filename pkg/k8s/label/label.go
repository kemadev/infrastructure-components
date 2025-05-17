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
