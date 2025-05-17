package label

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

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
		"app.kubernetes.io/name":       appName,
		"app.kubernetes.io/instance":   appInstance,
		"app.kubernetes.io/version":    appVersion,
		"app.kubernetes.io/component":  appComponent,
		"app.kubernetes.io/part-of":    appNamespace,
		"app.kubernetes.io/managed-by": pulumi.String("pulumi"),
	}
}

// DefaultLabels returns a set of default labels for the application instance as per Kubernetes convention,
// see https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels
func DefaultSelector(
	appInstance pulumi.StringInput,
	defaultLabels pulumi.StringMap,
) pulumi.StringMap {
	return pulumi.StringMap{
		"app.kubernetes.io/instance": defaultLabels["app.kubernetes.io/instance"],
	}
}
