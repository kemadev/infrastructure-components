package app

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/blang/semver"
	"github.com/caarlos0/svu/pkg/svu"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"vcs.kema.cloud/kema/runner-tools/pkg/git"
	"vcs.kema.run/kema/infrastructure-components/internal/pkg/businessunit"
	"vcs.kema.run/kema/infrastructure-components/internal/pkg/customer"
)

type AppParms struct {
	// Image reference (URL)
	ImageRef url.URL
	// Image tag
	ImageTag semver.Version
	// Runtime env
	RuntimeEnv string
	// OpenTelemetry endpoint URL
	OTelEndpointUrl url.URL
	// Application version, i.e. SemVer tag
	AppVersion semver.Version
	// Application name, i.e. repository name
	AppName string
	// Application namespace, which group it belogs to (e.g. shoppingcart, auth, ...)
	AppNamespace string
	// Business unit developing application
	BusinessUnitId businessunit.BusinessUnit
	// Customer intended to use application
	CustomerId customer.Customer
	// Cost center, which i
	CostCenter string
	// Cost allocation owner, who pays for the application, budget holder
	CostAllocationOwner businessunit.BusinessUnit
	// Team  responsible for application
	OperationsOwner businessunit.BusinessUnit
	// Recovery Point Objective (RPO) of resource
	Rpo time.Duration
	// Data classification resource is subject to (e.g. )
	DataClassification string
	// Compliance framework resource is subject to (e.g. )
	ComplianceFramework string
	// Time at which resource should expire, be deleted
	Expiration time.Time
	// Git repository URL
	ProjectUrl url.URL
	// Monitoring URL, (e.g. APM URL)
	MonitoringUrl url.URL
	// Port which application serves
	Port int
}

var (
	ErrNoRemoteURL        = fmt.Errorf("remote URL not found")
	ErrMultipleRemoteURLs = fmt.Errorf("found more than 1 remote URL")
	ErrInvalidUrl         = fmt.Errorf("repository remote URL is invlid")
)

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
	parsedUrl, err := url.Parse(gitUrl)
	if err != nil {
		return "", url.URL{}, fmt.Errorf("error parsing git repository url: %w", err)
	}
	return appName, *parsedUrl, nil
}

func getVersionFromGit() (semver.Version, error) {
	versionString, err := svu.Current(
		svu.WithPrefix("v"),
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

func mergeParams(ctx *pulumi.Context, params *AppParms) error {
	appName, repoUrl, err := getGitInfos()
	if err != nil {
		return fmt.Errorf("error getting git repository information: %w", err)
	}
	appVersion, err := getVersionFromGit()
	if err != nil {
		return fmt.Errorf("error getting app version from git: %w", err)
	}
	defParams := AppParms{
		AppName:    appName,
		ImageRef:   repoUrl,
		ImageTag:   appVersion,
		AppVersion: appVersion,
		RuntimeEnv: ctx.Stack(),
		// TODO stackref to collector project
		OTelEndpointUrl: url.URL{},
		ProjectUrl:      repoUrl,
		Port:            8080,
	}
	err = mergo.Merge(params, defParams)
	if err != nil {
		return fmt.Errorf("error filling app parameters: %w", err)
	}
	return nil
}

func DeployBasicHTTPApp(ctx *pulumi.Context, params AppParms) error {
	mergeParams(ctx, &params)

	// Must match kind mount
	appCodeVolume := "/app-code"

	_, err := corev1.NewNamespace(ctx, "namespace", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(params.AppName),
			Namespace: pulumi.String(params.AppName),
			Labels: func() pulumi.StringMap {
				enforce := "restricted"
				if ctx.Stack() == "dev" {
					// Allow using HostPath volume in dev
					enforce = "privileged"
				}
				return pulumi.StringMap{
					"app":                                pulumi.String(params.AppName),
					"pod-security.kubernetes.io/enforce": pulumi.String(enforce),
					"pod-security.kubernetes.io/enforce-version": pulumi.String("latest"),
					"pod-security.kubernetes.io/audit":           pulumi.String("restricted"),
					"pod-security.kubernetes.io/audit-version":   pulumi.String("latest"),
					"pod-security.kubernetes.io/warn":            pulumi.String("restricted"),
					"pod-security.kubernetes.io/warn-version":    pulumi.String("latest"),
				}
			}(),
		},
	})
	if err != nil {
		return err
	}

	cm, err := corev1.NewConfigMap(ctx, "configmap", &corev1.ConfigMapArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(params.AppName),
			Namespace: pulumi.String(params.AppName),
			Labels: pulumi.StringMap{
				"app": pulumi.String(params.AppName),
			},
		},
		Data: pulumi.StringMap{
			"RUNTIME_ENV":           pulumi.String(params.RuntimeEnv),
			"APP_VERSION":           pulumi.String(params.AppVersion.String()),
			"APP_NAME":              pulumi.String(params.AppName),
			"APP_NAMESPACE":         pulumi.String(params.AppNamespace),
			"OTEL_ENDPOINT_URL":     pulumi.String(params.OTelEndpointUrl.String()),
			"BUSINESS_UNIT_ID":      pulumi.String(params.BusinessUnitId),
			"CUSTOMER_ID":           pulumi.String(params.CustomerId),
			"COST_CENTER":           pulumi.String(params.CostCenter),
			"COST_ALLOCATION_OWNER": pulumi.String(params.CostAllocationOwner),
			"OPERATIONS_OWNER":      pulumi.String(params.OperationsOwner),
			"RPO":                   pulumi.String(params.Rpo.String()),
			"DATA_CLASSIFICATION":   pulumi.String(params.DataClassification),
			"COMPLIANCE_FRAMEWORK":  pulumi.String(params.ComplianceFramework),
			"EXPIRATION":            pulumi.String(params.Expiration.String()),
			"PROJECT_URL":           pulumi.String(params.ProjectUrl.String()),
			"MONITORING_URL":        pulumi.String(params.MonitoringUrl.String()),
		},
	})

	_, err = appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(params.AppName),
			Namespace: pulumi.String(params.AppName),
			Labels: pulumi.StringMap{
				"app": pulumi.String(params.AppName),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			// TODO use HPA / VPA
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(params.AppName),
				},
			},
			ProgressDeadlineSeconds: pulumi.Int(180),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name:      pulumi.String(params.AppName),
					Namespace: pulumi.String(params.AppName),
					Labels: pulumi.StringMap{
						"app": pulumi.String(params.AppName),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							// TODO dev pass args jq
							Args: pulumi.StringArray{},
							EnvFrom: corev1.EnvFromSourceArray{
								corev1.EnvFromSourceArgs{
									ConfigMapRef: corev1.ConfigMapEnvSourceArgs{
										Name: cm.Metadata.Name(),
									},
								},
							},
							LivenessProbe: corev1.ProbeArgs{
								InitialDelaySeconds: pulumi.Int(10),
								HttpGet: corev1.HTTPGetActionArgs{
									HttpHeaders: corev1.HTTPHeaderArray{},
									Path:        pulumi.String("/healthz"),
									Port:        pulumi.Int(params.Port),
								},
							},
							ReadinessProbe: corev1.ProbeArgs{
								HttpGet: corev1.HTTPGetActionArgs{
									HttpHeaders: corev1.HTTPHeaderArray{},
									Path:        pulumi.String("/readyz"),
									Port:        pulumi.Int(params.Port),
								},
							},
							// HACK Enable colorful output for air, remove once https://github.com/air-verse/air/pull/768 is merged
							Stdin: func() pulumi.Bool {
								if ctx.Stack() == "dev" {
									return pulumi.Bool(true)
								}
								return pulumi.Bool(false)
							}(),
							Tty: func() pulumi.Bool {
								if ctx.Stack() == "dev" {
									return pulumi.Bool(true)
								}
								return pulumi.Bool(false)
							}(),
							Image: pulumi.String(params.AppName + ":" + params.ImageTag.String()),
							Name:  pulumi.String(params.AppName),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(params.Port),
									Protocol:      pulumi.String("TCP"),
								},
							},
							VolumeMounts: func() corev1.VolumeMountArray {
								// Mount app code volume in dev
								if ctx.Stack() == "dev" {
									return corev1.VolumeMountArray{
										&corev1.VolumeMountArgs{
											Name:      pulumi.String(params.AppName),
											MountPath: pulumi.String("/app"),
										},
									}
								}
								return corev1.VolumeMountArray{}
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
							ImagePullPolicy: pulumi.String("IfNotPresent"),
						},
					},
					Volumes: func() corev1.VolumeArray {
						// Create app code volume in dev
						if ctx.Stack() == "dev" {
							return corev1.VolumeArray{
								corev1.VolumeArgs{
									Name: pulumi.String(params.AppName),
									HostPath: corev1.HostPathVolumeSourceArgs{
										Path: pulumi.String(appCodeVolume),
										Type: pulumi.String("Directory"),
									},
								},
							}
						}
						return corev1.VolumeArray{}
					}(),
					Resources: corev1.ResourceRequirementsArgs{
						// TODO refine values after benchmarking / load testing
						Requests: pulumi.StringMap{
							"cpu":    pulumi.String("500m"),
							"memory": pulumi.String("100m"),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = corev1.NewService(ctx, "service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(params.AppName),
			Namespace: pulumi.String(params.AppName),
			Labels: pulumi.StringMap{
				"app": pulumi.String(params.AppName),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					AppProtocol: pulumi.String("http"),
					Port:        pulumi.Int(params.Port),
					TargetPort:  pulumi.Any(params.Port),
					Protocol:    pulumi.String("TCP"),
				},
			},
			Selector: pulumi.StringMap{
				"app": pulumi.String(params.AppName),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
