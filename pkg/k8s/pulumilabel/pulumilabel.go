package pulumilabel

import (
	"github.com/kemadev/infrastructure-components/pkg/k8s/label"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
		label.LabelAppNameKey:      appName,
		label.LabelAppInstanceKey:  appInstance,
		label.LabelAppVersionKey:   appVersion,
		label.LabelAppComponentKey: appComponent,
		label.LabelAppNamespaceKey: appNamespace,
		label.LabelAppMangedByKey:  pulumi.String("pulumi"),
	}
}

// DefaultLabels returns a set of default labels for the application instance as per Kubernetes convention,
// see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
func DefaultSelector(
	appInstance pulumi.StringInput,
	defaultLabels pulumi.StringMap,
) pulumi.StringMap {
	return pulumi.StringMap{
		label.LabelAppInstanceKey: defaultLabels[label.LabelAppInstanceKey],
	}
}
