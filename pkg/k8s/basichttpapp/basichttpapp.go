package basichttpapp

import (
	"fmt"
	"maps"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/blang/semver"
	"github.com/caarlos0/svu/pkg/svu"
	"github.com/kemadev/framework-go/pkg/config"
	"github.com/kemadev/infrastructure-components/pkg/k8s/label"
	"github.com/kemadev/infrastructure-components/pkg/k8s/priorityclass"
	"github.com/kemadev/infrastructure-components/pkg/private/businessunit"
	"github.com/kemadev/infrastructure-components/pkg/private/complianceframework"
	"github.com/kemadev/infrastructure-components/pkg/private/costcenter"
	"github.com/kemadev/infrastructure-components/pkg/private/customer"
	"github.com/kemadev/infrastructure-components/pkg/private/dataclassification"
	"github.com/kemadev/runner-tools/pkg/git"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	autoscalingv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/autoscaling/v2"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A AppParms contains all the parameters needed to deploy a basic HTTP application
type AppParms struct {
	// ImageRef is the base image reference, e.g. registry.host.tld/repo/imagename.
	ImageRef url.URL
	// ImageTag is the image tag, as a SemVer tag. It should not be manually set, as it
	// is automatically set to AppVersion.
	ImageTag semver.Version
	// RuntimeEnv is the runtime environment, i.e. Pulumi stack name. It is used as a suffix to the application name
	// in application instance name, ensuring uniqueness across environments. Please note that [github.com/kemadev/framework-go/pkg/config.Env_dev]
	// is a special case, tweaking the application deployment to ease local development.
	RuntimeEnv string
	// OTelEndpointUrl is the OpenTelemetry collector endpoint URL.
	OTelEndpointUrl url.URL
	// OtelExporterCompression is the OpenTelemetry exporter compression method.
	OtelExporterCompression string
	// AppVersion is the application version, as a SemVer tag.
	AppVersion semver.Version
	// AppName is the application name, i.e. the name of the repository.
	AppName string
	// AppNamespace is the application namespace, i.e. which group it belongs to (e.g. shoppingcart, auth, ...)
	AppNamespace string
	// AppComponent is the application role, e.g. frontend, api, database, ...
	AppComponent string
	// BusinessUnitId is the business unit developing application.
	BusinessUnitId businessunit.BusinessUnit
	// CustomerId is the customer using the application.
	CustomerId customer.Customer
	// CostCenter is the cost center to which the application belongs.
	CostCenter costcenter.CostCenter
	// CostAllocationOwner is the business unit allocating resources to the application, i.e. the budget holder.
	CostAllocationOwner businessunit.BusinessUnit
	// OperationsOwner is the business unit responsible for developing and maintaining the application.
	OperationsOwner businessunit.BusinessUnit
	// Rpo is the recovery point objective, i.e. the maximum amount of data that can be lost in case of a failure.
	Rpo time.Duration
	// DataClassification is the data classification the application is subject to.
	DataClassification dataclassification.DataClassification
	// ComplianceFramework is the compliance framework the application is subject to.
	ComplianceFramework complianceframework.ComplianceFramework
	// Expiration is the expiration date of the application, i.e. when should be decommissioned.
	Expiration time.Time
	// ProjectUrl is the URL of the project, i.e. the URL of the repository.
	ProjectUrl url.URL
	// MonitoringUrl is the URL of the monitoring system, e.g. the URL of the APM.
	MonitoringUrl url.URL
	// Port is the port on which the application is listening.
	Port int
	// HTTPReadTimeout is the HTTP read timeout, in seconds.
	HTTPReadTimeout int
	// HTTPWriteTimeout is the HTTP write timeout, in seconds.
	HTTPWriteTimeout int
	// CPURequestMiliCPU is the CPU request for the pod, in mili vCPU (will be set as `strconv.Itoa(CPURequestMiliCPU) + "m"`)
	CPURequestMiliCPU int
	// CPULimitMiliCPU is the CPU limit for the pod, in mili vCPU (will be set as `strconv.Itoa(CPULimitMiliCPU) + "m"`). It will also be used to
	// set GOMAXPROCS to 1/1000th of this value, floored (you should only specify multiples of 1000)
	CPULimitMiliCPU int
	// MemoryRequestMiB is the memory request for the pod, in MiB (will be set as `strconv.Itoa(MemoryRequestMiB) + "MiB"`)
	MemoryRequestMiB int
	// MemoryLimitMiB is the memory limit for the pod, in MiB (will be set as `strconv.Itoa(MemoryLimitMiB) + "MiB"`). It will also be used to
	// set GOMEMLIMIT to 95% of this value.
	MemoryLimitMiB int
	// MinReplicas is the minimum number of replicas for the pod, used for HPA
	MinReplicas int
	// MaxReplicas is the maximum number of replicas for the pod, used for HPA
	MaxReplicas int
	// ProgressDeadlineSeconds is the maximum time in seconds for the deployment to be ready.
	ProgressDeadlineSeconds int
	// DevImagePullPolicy is the image pull policy to use.
	ImagePullPolicy string
	// PodAffinity is the pod affinity to use for the pod. Should be set when know pods communicate alot with the application.
	PodAffinity corev1.AffinityPtrInput
	// NodeSelectors is the node selectors to use for the pod.
	NodeSelectors pulumi.StringMapInput
	// PriorityClassName is the name of the priority class to use for the pod.
	PriorityClassName string
	// TopologySpreadConstraints is the list of topology spread constraints to use for the pod.
	TopologySpreadConstraints corev1.TopologySpreadConstraintArray
	// HorizontalPodAutoscalerBehavior is the behavior of the HPA.
	HorizontalPodAutoscalerBehavior autoscalingv2.HorizontalPodAutoscalerBehaviorPtrInput
	// HorizontalPodAutoscalerBehaviorMetricSpec is the metric spec for the HPA behavior.
	HorizontalPodAutoscalerBehaviorMetricSpec autoscalingv2.MetricSpecArray
	// DevDnsAdditionalNameservers is the list of additional nameservers use in pods (dev stack only).
	DevDnsAdditionalNameservers []string
	// DevGoPrivateString is the GOPRIVATE string to set (dev stack only).
	DevGoPrivateString string
	// DevVolumeMountPath is the path where the organization code is mounted on node (dev stack only).
	DevVolumeMountPath string
}

var (
	// ErrNoRemoteURL is a sentinel error indicating that no remote URL was found in the git repository.
	ErrNoRemoteURL = fmt.Errorf("remote URL not found")
	// ErrMultipleRemoteURLs is a sentinel error indicating that multiple remote URLs were found in the git repository.
	ErrMultipleRemoteURLs = fmt.Errorf("found more than 1 remote URL")
	// ErrInvalidUrl is a sentinel error indicating that the remote URL is invalid.
	ErrInvalidUrl = fmt.Errorf("repository remote URL is invlid")
	// ErrGitTokenNotFound is a sentinel error indicating that the git token was not found in the pulumi config.
	ErrGitTokenNotFound = fmt.Errorf(
		"git-token not found in pulumi config, please set it using `pulumi config set --secret git-token`, with a personal access token (with read access to organization's repositories) as value",
	)
)

// getGitInfos returns the application name and the remote URL of the git repository, based on the git remote "origin", and
// an error if any.
func getGitInfos() (string, url.URL, error) {
	repo, err := git.GetGitRepo()
	if err != nil {
		return "", url.URL{}, fmt.Errorf("error getting git repository: %w", err)
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", url.URL{}, fmt.Errorf("error getting git remote origin: %w", err)
	}
	urls := remote.Config().URLs
	if len(urls) < 1 {
		return "", url.URL{}, ErrNoRemoteURL
	} else if len(urls) > 1 {
		return "", url.URL{}, ErrMultipleRemoteURLs
	}
	gitUrl, err := git.GetGitBasePathWithRepo(repo)
	if err != nil {
		return "", url.URL{}, fmt.Errorf("error getting git repository base path: %w", err)
	}
	urlParts := strings.Split(gitUrl, "/")
	if len(urlParts) < 2 {
		return "", url.URL{}, fmt.Errorf("remote url %s: %w", gitUrl, ErrInvalidUrl)
	}
	appName := strings.Join(urlParts[len(urlParts)-1:], "")
	gitUrlWithScheme := "https://" + gitUrl
	parsedUrl, err := url.Parse(gitUrlWithScheme)
	if err != nil {
		return "", url.URL{}, fmt.Errorf("error parsing git repository url: %w", err)
	}
	return appName, *parsedUrl, nil
}

// getVersionFromGit returns the application version from the git repository, based on the current tag, and an error if any.
func getVersionFromGit() (semver.Version, error) {
	versionString, err := svu.Current(
		svu.StripPrefix(),
	)
	if err != nil {
		return semver.Version{}, fmt.Errorf("error getting app version from git: %w", err)
	}
	version, err := semver.Parse(versionString)
	if err != nil {
		return semver.Version{}, fmt.Errorf("error parsing app version from git: %w", err)
	}
	return version, nil
}

// validateParams validates the application parameters, returning an error if any of them is invalid.
// Not all parameters are enforced, as some of them are optional.
// NOTE(maintainers): When adding new parameters, add them to this function, even if they are not enforced, by commenting them out.
func validateParams(params *AppParms) error {
	// Enforce parameters, with commented-out non-enforced values
	if params.ImageRef.String() == "" {
		return fmt.Errorf("ImageRef cannot be empty")
	}
	if params.ImageTag.String() == "" {
		return fmt.Errorf("ImageTag cannot be empty")
	}
	if params.RuntimeEnv == "" {
		return fmt.Errorf("RuntimeEnv cannot be empty")
	}
	if params.OTelEndpointUrl.String() == "" {
		return fmt.Errorf("OTelEndpointUrl cannot be empty")
	}
	if params.OtelExporterCompression == "" {
		return fmt.Errorf("OtelCompression cannot be empty")
	}
	if params.AppVersion.String() == "" {
		return fmt.Errorf("AppVersion cannot be empty")
	}
	if params.AppName == "" {
		return fmt.Errorf("AppName cannot be empty")
	}
	if params.AppNamespace == "" {
		return fmt.Errorf("AppNamespace cannot be empty")
	}
	if params.AppComponent == "" {
		return fmt.Errorf("AppComponent cannot be empty")
	}
	if params.BusinessUnitId == "" {
		return fmt.Errorf("BusinessUnitId cannot be empty")
	}
	if params.CustomerId == "" {
		return fmt.Errorf("CustomerId cannot be empty")
	}
	if params.CostCenter == "" {
		return fmt.Errorf("CostCenter cannot be empty")
	}
	if params.CostAllocationOwner == "" {
		return fmt.Errorf("CostAllocationOwner cannot be empty")
	}
	if params.OperationsOwner == "" {
		return fmt.Errorf("OperationsOwner cannot be empty")
	}
	if params.Rpo == 0 {
		return fmt.Errorf("Rpo cannot be zero")
	}
	if params.DataClassification == "" {
		return fmt.Errorf("DataClassification cannot be empty")
	}
	if params.ComplianceFramework == "" {
		return fmt.Errorf("ComplianceFramework cannot be empty")
	}
	// if params.Expiration.IsZero() {
	// 	return fmt.Errorf("Expiration cannot be zero")
	// }
	if params.ProjectUrl.String() == "" {
		return fmt.Errorf("ProjectUrl cannot be empty")
	}
	if params.MonitoringUrl.String() == "" {
		return fmt.Errorf("MonitoringUrl cannot be empty")
	}
	if params.Port == 0 {
		return fmt.Errorf("Port cannot be zero")
	}
	if params.HTTPReadTimeout == 0 {
		return fmt.Errorf("HTTPReadTimeout cannot be zero")
	}
	if params.HTTPWriteTimeout == 0 {
		return fmt.Errorf("HTTPWriteTimeout cannot be zero")
	}
	if params.CPURequestMiliCPU == 0 {
		return fmt.Errorf("CPURequest cannot be zero")
	}
	// if params.CPULimit == 0 {
	// 	return fmt.Errorf("CPULimit cannot be zero")
	// }
	if params.MemoryRequestMiB == 0 {
		return fmt.Errorf("MemoryRequest cannot be zero")
	}
	// if params.MemoryLimit == 0 {
	// 	return fmt.Errorf("MemoryLimit cannot be zero")
	// }
	if params.MinReplicas == 0 {
		return fmt.Errorf("MinReplicas cannot be zero")
	}
	if params.MaxReplicas == 0 {
		return fmt.Errorf("MaxReplicas cannot be zero")
	}
	if params.ProgressDeadlineSeconds == 0 {
		return fmt.Errorf("ProgressDeadlineSeconds cannot be zero")
	}
	if params.ImagePullPolicy == "" {
		return fmt.Errorf("ImagePullPolicy cannot be empty")
	}
	// if params.PodAffinity == nil {
	// 	return fmt.Errorf("PodAffinity cannot be nil")
	// }
	// if params.NodeSelectors == nil {
	// 	return fmt.Errorf("NodeSelectors cannot be nil")
	// }
	if params.PriorityClassName == "" {
		return fmt.Errorf("PriorityClassName cannot be empty")
	}
	if params.TopologySpreadConstraints == nil {
		return fmt.Errorf("TopologySpreadConstraints cannot be nil")
	}
	if params.HorizontalPodAutoscalerBehavior == nil {
		return fmt.Errorf("HorizontalPodAutoscalerBehavior cannot be nil")
	}
	if params.HorizontalPodAutoscalerBehaviorMetricSpec == nil {
		return fmt.Errorf("HorizontalPodAutoscalerBehaviorMetricSpec cannot be nil")
	}
	// if len(params.DevDnsAdditionalNameservers) == 0 {
	// 	return fmt.Errorf("DevDnsAdditionalNameservers cannot be empty")
	// }
	// if params.DevGoPrivateString == "" {
	// 	return fmt.Errorf("DevGoPrivateString cannot be empty")
	// }
	if params.DevVolumeMountPath == "" {
		return fmt.Errorf("DevVolumeMountPath cannot be empty")
	}
	return nil
}

// mergeParams merges the default parameters with the provided parameters, returning an error if any of them is invalid.
func mergeParams(
	params *AppParms,
	appName string,
	appInstance string,
	repoUrl url.URL,
	runtimeEnv string,
) error {
	appVersion, err := getVersionFromGit()
	if err != nil {
		return fmt.Errorf("error getting app version from git: %w", err)
	}
	defParams := AppParms{
		AppName:             appName,
		ImageRef:            repoUrl,
		ImageTag:            appVersion,
		AppVersion:          appVersion,
		DataClassification:  dataclassification.DataClassificationNone,
		ComplianceFramework: complianceframework.ComplianceFrameworkNone,
		RuntimeEnv:          runtimeEnv,
		OTelEndpointUrl: url.URL{
			Scheme: "grpc",
			Host:   "string",
			Path:   "string",
		},
		OtelExporterCompression: "gzip",
		ProjectUrl: func() url.URL {
			t := repoUrl
			t.Scheme = "https"
			return t
		}(),
		Port:                    8080,
		HTTPReadTimeout:         5,
		HTTPWriteTimeout:        10,
		CPURequestMiliCPU:       500,
		MemoryRequestMiB:        500,
		MinReplicas:             1,
		MaxReplicas:             10,
		ImagePullPolicy:         "IfNotPresent",
		ProgressDeadlineSeconds: 180,
		PriorityClassName:       priorityclass.PriorityClassNormal,
		TopologySpreadConstraints: corev1.TopologySpreadConstraintArray{
			// Spread pods across regions, best effort
			corev1.TopologySpreadConstraintArgs{
				MaxSkew: pulumi.Int(1),
				LabelSelector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app.kubernetes.io/instance": pulumi.String(appInstance),
					},
				},
				MatchLabelKeys: pulumi.StringArray{
					pulumi.String("pod-template-hash"),
				},
				TopologyKey:       pulumi.String(label.LabelTopologyRegionKey),
				WhenUnsatisfiable: pulumi.String("ScheduleAnyway"),
			},
			// Spread pods across zones, best effort
			corev1.TopologySpreadConstraintArgs{
				MaxSkew: pulumi.Int(1),
				LabelSelector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app.kubernetes.io/instance": pulumi.String(appInstance),
					},
				},
				MatchLabelKeys: pulumi.StringArray{
					pulumi.String("pod-template-hash"),
				},
				TopologyKey:       pulumi.String(label.LabelTopologyZoneKey),
				WhenUnsatisfiable: pulumi.String("ScheduleAnyway"),
			},
			// Spread pods across datacenters, best effort
			corev1.TopologySpreadConstraintArgs{
				MaxSkew: pulumi.Int(1),
				LabelSelector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app.kubernetes.io/instance": pulumi.String(appInstance),
					},
				},
				MatchLabelKeys: pulumi.StringArray{
					pulumi.String("pod-template-hash"),
				},
				TopologyKey:       pulumi.String(label.LabelTopologyDatacenterKey),
				WhenUnsatisfiable: pulumi.String("ScheduleAnyway"),
			},
			// Spread pods across datacenter zones, best effort
			corev1.TopologySpreadConstraintArgs{
				MaxSkew: pulumi.Int(1),
				LabelSelector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app.kubernetes.io/instance": pulumi.String(appInstance),
					},
				},
				MatchLabelKeys: pulumi.StringArray{
					pulumi.String("pod-template-hash"),
				},
				TopologyKey:       pulumi.String(label.LabelTopologyDatacenterZoneKey),
				WhenUnsatisfiable: pulumi.String("ScheduleAnyway"),
			},
			// Spread pods across aisles, best effort
			corev1.TopologySpreadConstraintArgs{
				MaxSkew: pulumi.Int(1),
				LabelSelector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app.kubernetes.io/instance": pulumi.String(appInstance),
					},
				},
				MatchLabelKeys: pulumi.StringArray{
					pulumi.String("pod-template-hash"),
				},
				TopologyKey:       pulumi.String(label.LabelTopologyDatacenterAisleKey),
				WhenUnsatisfiable: pulumi.String("ScheduleAnyway"),
			},
			// Spread pods across racks, best effort
			corev1.TopologySpreadConstraintArgs{
				MaxSkew: pulumi.Int(1),
				LabelSelector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app.kubernetes.io/instance": pulumi.String(appInstance),
					},
				},
				MatchLabelKeys: pulumi.StringArray{
					pulumi.String("pod-template-hash"),
				},
				TopologyKey:       pulumi.String(label.LabelTopologyDatacenterRackKey),
				WhenUnsatisfiable: pulumi.String("ScheduleAnyway"),
			},
			// Spread pods across nodes, best effort
			corev1.TopologySpreadConstraintArgs{
				MaxSkew: pulumi.Int(1),
				LabelSelector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app.kubernetes.io/instance": pulumi.String(appInstance),
					},
				},
				MatchLabelKeys: pulumi.StringArray{
					pulumi.String("pod-template-hash"),
				},
				TopologyKey:       pulumi.String(label.LabelTopologyHostnameKey),
				WhenUnsatisfiable: pulumi.String("ScheduleAnyway"),
			},
		},
		HorizontalPodAutoscalerBehavior: &autoscalingv2.HorizontalPodAutoscalerBehaviorArgs{
			ScaleDown: &autoscalingv2.HPAScalingRulesArgs{
				// Downscale max 30%/minute
				Policies: autoscalingv2.HPAScalingPolicyArray{
					&autoscalingv2.HPAScalingPolicyArgs{
						Type:          pulumi.String("Percent"),
						PeriodSeconds: pulumi.Int(60),
						Value:         pulumi.Int(30),
					},
				},
				SelectPolicy: pulumi.String("Min"),
			},
			ScaleUp: &autoscalingv2.HPAScalingRulesArgs{
				// Upscale max 30%/minute
				Policies: autoscalingv2.HPAScalingPolicyArray{
					&autoscalingv2.HPAScalingPolicyArgs{
						Type:          pulumi.String("Percent"),
						PeriodSeconds: pulumi.Int(60),
						Value:         pulumi.Int(30),
					},
				},
				SelectPolicy: pulumi.String("Max"),
			},
		},
		HorizontalPodAutoscalerBehaviorMetricSpec: autoscalingv2.MetricSpecArray{
			&autoscalingv2.MetricSpecArgs{
				Type: pulumi.String("Resource"),
				Resource: &autoscalingv2.ResourceMetricSourceArgs{
					Name: pulumi.String("cpu"),
					Target: &autoscalingv2.MetricTargetArgs{
						Type:               pulumi.String("Utilization"),
						AverageUtilization: pulumi.Int(70),
					},
				},
			},
		},
		DevDnsAdditionalNameservers: []string{
			"1.1.1.1",
		},
		DevVolumeMountPath: "/app",
		DevGoPrivateString: "github.com/kemadev",
	}
	err = mergo.Merge(params, defParams)
	if err != nil {
		return fmt.Errorf("error filling app parameters: %w", err)
	}
	err = validateParams(params)
	if err != nil {
		return fmt.Errorf("error validating app parameters: %w", err)
	}
	return nil
}

// DeployBasicHTTPApp deploys a basic HTTP application to the Kubernetes cluster, using the provided parameters merged with the default ones,
// and returns an error if any of the parameters is invalid or if the deployment fails.
func DeployBasicHTTPApp(ctx *pulumi.Context, params AppParms) error {
	appName, repoUrl, err := getGitInfos()
	if err != nil {
		return fmt.Errorf("error getting git repository information: %w", err)
	}

	// Runtime environment, i.e. Pulumi stack name
	runtimeEnv := ctx.Stack()

	// Application instance to use, using runtime env as suffix to distinguish different stacks, e.g. to distinguish review applications using their stack name (i.e. branch name)
	appInstance := appName + "-" + runtimeEnv

	err = mergeParams(&params, appName, appInstance, repoUrl, runtimeEnv)
	if err != nil {
		return fmt.Errorf("failed to apply default application parameters: %w", err)
	}

	sharedLabels := label.DefaultLabels(
		pulumi.String(params.AppName),
		pulumi.String(appInstance),
		pulumi.String(params.AppVersion.String()),
		pulumi.String(params.AppComponent),
		pulumi.String(params.AppNamespace),
	)
	basicSelector := label.DefaultSelector(
		pulumi.String(appInstance),
		sharedLabels,
	)

	// Namespace to deploy to
	namespace := appInstance

	// Where to mount organization's code, must match kind mount
	orgCodeVolume := "/git-vcs-org"

	// Application namespace
	_, err = corev1.NewNamespace(ctx, "namespace", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(namespace),
			Namespace: pulumi.String(namespace),
			Labels: func() pulumi.StringMap {
				enforce := "restricted"
				if ctx.Stack() == config.Env_dev {
					// Allow using HostPath volume in dev, as well as other side-effects, but should be mitigated by SecurityContext
					enforce = "privileged"
				}
				labels := pulumi.StringMap{
					// See https://kubernetes.io/docs/concepts/security/pod-security-admission/#pod-security-admission-labels-for-namespaces
					"pod-security.kubernetes.io/enforce":         pulumi.String(enforce),
					"pod-security.kubernetes.io/enforce-version": pulumi.String("latest"),
					"pod-security.kubernetes.io/audit":           pulumi.String("restricted"),
					"pod-security.kubernetes.io/audit-version":   pulumi.String("latest"),
					"pod-security.kubernetes.io/warn":            pulumi.String("restricted"),
					"pod-security.kubernetes.io/warn-version":    pulumi.String("latest"),
				}
				maps.Copy(labels, sharedLabels)
				// Allow shared gateway access to this namespace
				gatewayAttachmentEnableLabel := pulumi.StringMap{
					label.SharedGatewayAccessLabelKey: pulumi.String(
						label.SharedGatewayAccessLabelValue,
					),
				}
				maps.Copy(labels, gatewayAttachmentEnableLabel)
				return labels
			}(),
		},
	})
	if err != nil {
		return err
	}

	// ConfigMap providing common environment variable to containers
	cm, err := corev1.NewConfigMap(ctx, "env-configmap", &corev1.ConfigMapArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(appInstance),
			Namespace: pulumi.String(namespace),
			Labels:    sharedLabels,
		},
		Data: func() pulumi.StringMap {
			envMap := pulumi.StringMap{
				config.EnvVarKeyRuntimeEnv:   pulumi.String(params.RuntimeEnv),
				config.EnvVarKeyAppVersion:   pulumi.String(params.AppVersion.String()),
				config.EnvVarKeyAppName:      pulumi.String(appInstance),
				config.EnvVarKeyAppNamespace: pulumi.String(params.AppNamespace),
				config.EnvVarKeyOtelEndpointURL: pulumi.String(
					params.OTelEndpointUrl.String(),
				),
				config.EnvVarKeyOtelExporterCompression: pulumi.String(
					params.OtelExporterCompression,
				),
				config.EnvVarKeyHTTPServePort: pulumi.String(
					strconv.Itoa(params.Port),
				),
				config.EnvVarKeyHTTPReadTimeout: pulumi.String(
					strconv.Itoa(params.HTTPReadTimeout),
				),
				config.EnvVarKeyHTTPWriteTimeout: pulumi.String(
					strconv.Itoa(params.HTTPWriteTimeout),
				),
				config.EnvVarKeyBusinessUnitId:      pulumi.String(params.BusinessUnitId),
				config.EnvVarKeyCustomerId:          pulumi.String(params.CustomerId),
				config.EnvVarKeyCostCenter:          pulumi.String(params.CostCenter),
				config.EnvVarKeyCostAllocationOwner: pulumi.String(params.CostAllocationOwner),
				config.EnvVarKeyOperationsOwner:     pulumi.String(params.OperationsOwner),
				config.EnvVarKeyRpo:                 pulumi.String(params.Rpo.String()),
				config.EnvVarKeyDataClassification:  pulumi.String(params.DataClassification),
				config.EnvVarKeyComplianceFramework: pulumi.String(params.ComplianceFramework),
				config.EnvVarKeyProjectUrl:          pulumi.String(params.ProjectUrl.String()),
				config.EnvVarKeyMonitoringUrl: pulumi.String(
					params.MonitoringUrl.String(),
				),
			}
			if !params.Expiration.IsZero() {
				envMap[config.EnvVarKeyExpiration] = pulumi.String(params.Expiration.String())
			}
			if params.CPULimitMiliCPU != 0 {
				// Match allocated CPUs, floored
				envMap["GOMAXPROCS"] = pulumi.String(strconv.Itoa(params.CPULimitMiliCPU / 1000))
			}
			if params.MemoryLimitMiB != 0 {
				envMap["GOMEMLIMIT"] = pulumi.String(
					// Match allocated memory, with little room, floored
					pulumi.String(strconv.Itoa(params.MemoryLimitMiB*95/100) + "MiB"),
				)
			}
			return envMap
		}(),
	})

	netRcFileName := ".netrc"
	var gitSecret *corev1.Secret
	if ctx.Stack() == config.Env_dev {
		gitTokenConfigKey := "git-token"
		gitToken, present := ctx.GetConfig(
			ctx.Project() + ":" + gitTokenConfigKey,
		)
		if !present || gitToken == "" {
			return ErrGitTokenNotFound
		}
		gitSecret, err = corev1.NewSecret(ctx, "gitSecret", &corev1.SecretArgs{
			Type: pulumi.String("Opaque"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(appInstance + "-git-secret"),
				Namespace: pulumi.String(namespace),
				Labels:    sharedLabels,
			},
			StringData: pulumi.StringMap{
				netRcFileName: pulumi.String("machine " + params.ProjectUrl.Host + `
login git
password ` + gitToken),
			},
		})
		if err != nil {
			return err
		}
	}

	// Application deployment
	deployment, err := appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(appInstance),
			Namespace: pulumi.String(namespace),
			Labels:    sharedLabels,
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: basicSelector,
			},
			ProgressDeadlineSeconds: pulumi.Int(params.ProgressDeadlineSeconds),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name:      pulumi.String(appInstance),
					Namespace: pulumi.String(namespace),
					Labels:    sharedLabels,
				},
				Spec: &corev1.PodSpecArgs{
					PriorityClassName:         pulumi.String(params.PriorityClassName),
					TopologySpreadConstraints: params.TopologySpreadConstraints,
					NodeSelector:              params.NodeSelectors,
					Affinity:                  params.PodAffinity,
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							EnvFrom: corev1.EnvFromSourceArray{
								corev1.EnvFromSourceArgs{
									ConfigMapRef: corev1.ConfigMapEnvSourceArgs{
										Name: cm.Metadata.Name(),
									},
								},
							},
							Env: func() corev1.EnvVarArray {
								if ctx.Stack() == config.Env_dev {
									if params.DevGoPrivateString == "" {
										return nil
									}
									// Set GOPRIVATE for dev environment
									return corev1.EnvVarArray{
										corev1.EnvVarArgs{
											Name:  pulumi.String("GOPRIVATE"),
											Value: pulumi.String(params.DevGoPrivateString),
										},
									}
								}
								return nil
							}(),
							WorkingDir: func() pulumi.StringPtrInput {
								if ctx.Stack() == config.Env_dev {
									// Use mounted volume's project dir as working dir in dev
									return pulumi.String(
										params.DevVolumeMountPath + "/" + params.AppName,
									)
								}
								return nil
							}(),
							LivenessProbe: func() corev1.ProbePtrInput {
								if ctx.Stack() == config.Env_dev {
									return nil
								}
								return corev1.ProbeArgs{
									InitialDelaySeconds: pulumi.Int(10),
									HttpGet: corev1.HTTPGetActionArgs{
										Path: pulumi.String(config.HTTPLivenessCheckPath),
										Port: pulumi.Int(params.Port),
									},
								}
							}(),
							ReadinessProbe: func() corev1.ProbePtrInput {
								if ctx.Stack() == config.Env_dev {
									return nil
								}
								return corev1.ProbeArgs{
									HttpGet: corev1.HTTPGetActionArgs{
										Path: pulumi.String(config.HTTPReadinessCheckPath),
										Port: pulumi.Int(params.Port),
									},
								}
							}(),
							// HACK Enable colorful output for air, remove once https://github.com/air-verse/air/pull/768 is merged
							Stdin: func() pulumi.Bool {
								if ctx.Stack() == config.Env_dev {
									return pulumi.Bool(true)
								}
								return pulumi.Bool(false)
							}(),
							// HACK Enable colorful output for air, remove once https://github.com/air-verse/air/pull/768 is merged
							Tty: func() pulumi.Bool {
								if ctx.Stack() == config.Env_dev {
									return pulumi.Bool(true)
								}
								return pulumi.Bool(false)
							}(),
							Image: pulumi.String(
								params.ImageRef.Host + params.ImageRef.Path + ":" + params.ImageTag.String(),
							),
							Name: pulumi.String(appInstance),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(params.Port),
									Protocol:      pulumi.String("TCP"),
								},
							},
							VolumeMounts: func() corev1.VolumeMountArrayInput {
								if ctx.Stack() == config.Env_dev {
									return corev1.VolumeMountArray{
										// Mount organization code volume in dev
										&corev1.VolumeMountArgs{
											Name:      pulumi.String(appInstance),
											MountPath: pulumi.String(params.DevVolumeMountPath),
										},
										// Mountnetrc file volume in dev
										&corev1.VolumeMountArgs{
											Name: gitSecret.Metadata.Name().Elem(),
											MountPath: pulumi.String(
												"/home/nonroot/" + netRcFileName,
											),
											SubPath: pulumi.String(netRcFileName),
										},
									}
								}
								return nil
							}(),
							SecurityContext: corev1.SecurityContextArgs{
								AllowPrivilegeEscalation: pulumi.Bool(false),
								RunAsNonRoot:             pulumi.Bool(true),
								SeccompProfile: corev1.SeccompProfileArgs{
									Type: pulumi.String("RuntimeDefault"),
								},
								Capabilities: corev1.CapabilitiesArgs{
									Drop: pulumi.StringArray{
										pulumi.String("ALL"),
									},
								},
							},
							ImagePullPolicy: pulumi.String(params.ImagePullPolicy),
							Resources: corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"cpu": pulumi.String(
										strconv.Itoa(params.CPURequestMiliCPU) + "m",
									),
									"memory": pulumi.String(
										strconv.Itoa(params.MemoryRequestMiB) + "Mi",
									),
								},
								Limits: func() pulumi.StringMapInput {
									l := pulumi.StringMap{}
									if params.CPULimitMiliCPU != 0 {
										l["cpu"] = pulumi.String(
											strconv.Itoa(params.CPULimitMiliCPU) + "m",
										)
									}
									if params.MemoryLimitMiB != 0 {
										l["memory"] = pulumi.String(
											strconv.Itoa(params.MemoryLimitMiB) + "Mi",
										)
									}
									return l
								}(),
							},
						},
					},
					DnsConfig: func() corev1.PodDNSConfigPtrInput {
						if ctx.Stack() == config.Env_dev {
							// Add resolver for public DNS in dev
							return corev1.PodDNSConfigArgs{
								Nameservers: func() pulumi.StringArrayInput {
									if len(params.DevDnsAdditionalNameservers) == 0 {
										return nil
									}
									nameservers := make(
										pulumi.StringArray,
										len(params.DevDnsAdditionalNameservers),
									)
									for i, ns := range params.DevDnsAdditionalNameservers {
										nameservers[i] = pulumi.String(ns)
									}
									return nameservers
								}(),
							}
						}
						return nil
					}(),
					Volumes: func() corev1.VolumeArrayInput {
						if ctx.Stack() == config.Env_dev {
							return corev1.VolumeArray{
								// Create organization code volume in dev
								corev1.VolumeArgs{
									Name: pulumi.String(appInstance),
									HostPath: corev1.HostPathVolumeSourceArgs{
										Path: pulumi.String(orgCodeVolume),
										Type: pulumi.String("Directory"),
									},
								},
								// Create netrc file volume in dev
								corev1.VolumeArgs{
									Name: gitSecret.Metadata.Name().Elem(),
									Projected: corev1.ProjectedVolumeSourceArgs{
										Sources: corev1.VolumeProjectionArray{
											corev1.VolumeProjectionArgs{
												Secret: corev1.SecretProjectionArgs{
													Name: gitSecret.Metadata.Name(),
													Items: corev1.KeyToPathArray{
														corev1.KeyToPathArgs{
															Key:  pulumi.String(netRcFileName),
															Path: pulumi.String(netRcFileName),
														},
													},
												},
											},
										},
									},
								},
							}
						}
						return nil
					}(),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = autoscalingv2.NewHorizontalPodAutoscaler(
		ctx,
		"hpa",
		&autoscalingv2.HorizontalPodAutoscalerArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(appInstance),
				Namespace: pulumi.String(namespace),
				Labels:    sharedLabels,
			},
			Spec: &autoscalingv2.HorizontalPodAutoscalerSpecArgs{
				MinReplicas: pulumi.Int(params.MinReplicas),
				MaxReplicas: pulumi.Int(params.MaxReplicas),
				ScaleTargetRef: &autoscalingv2.CrossVersionObjectReferenceArgs{
					Kind:       deployment.Kind,
					ApiVersion: deployment.ApiVersion,
					Name:       deployment.Metadata.Name().Elem(),
				},
				Behavior: params.HorizontalPodAutoscalerBehavior,
				Metrics:  params.HorizontalPodAutoscalerBehaviorMetricSpec,
			},
		},
	)
	if err != nil {
		return err
	}

	// Application service
	_, err = corev1.NewService(ctx, "service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(appInstance),
			Namespace: pulumi.String(namespace),
			Labels:    sharedLabels,
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					AppProtocol: pulumi.String("kubernetes.io/h2c"),
					Port:        pulumi.Int(params.Port),
					TargetPort:  pulumi.Any(params.Port),
					Protocol:    pulumi.String("TCP"),
				},
			},
			Selector: basicSelector,
			// Prioritize close endpoints, best-effort, see https://kubernetes.io/docs/reference/networking/virtual-ips/#traffic-distribution
			TrafficDistribution: pulumi.String("PreferClose"),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
