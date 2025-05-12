package app

import (
	"fmt"
	"net/url"
	"strings"

	"dario.cat/mergo"
	"github.com/blang/semver"
	"github.com/caarlos0/svu/pkg/svu"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"vcs.kema.cloud/kema/runner-tools/pkg/git"
)

type AppParms struct {
	ImageRef            string
	ImageTag            string
	RuntimeEnv          string
	OTelEndpointUrl     url.URL
	AppVersion          semver.Version
	AppName             string
	AppNamespace        string
	BusinessUnitId      string
	CustomerId          string
	CostCenter          string
	CostAllocationOwner string
	OperationsOwner     string
	Rpo                 string
	DataClassification  string
	ComplianceFramework string
	Expiration          string
	ProjectUrl          string
	MonitoringUrl       string
}

var (
	ErrNoRemoteURL        = fmt.Errorf("remote URL not found")
	ErrMultipleRemoteURLs = fmt.Errorf("found more than 1 remote URL")
	ErrInvalidUrl         = fmt.Errorf("repository remote URL is invlid")
)

func getGitInfos() (string, string, error) {
	repo, err := git.GetGitRepo()
	if err != nil {
		return "", "", fmt.Errorf("error getting git repository: %w", err)
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", "", fmt.Errorf("error getting git remote origin: %w", err)
	}
	urls := remote.Config().URLs
	if len(urls) < 1 {
		return "", "", ErrNoRemoteURL
	} else if len(urls) > 1 {
		return "", "", ErrMultipleRemoteURLs
	}
	gitUrl, err := git.GetGitBasePathWithRepo(repo)
	if err != nil {
		return "", "", fmt.Errorf("error getting git repository base path: %w", err)
	}
	urlParts := strings.Split(gitUrl, "/")
	if len(urlParts) < 2 {
		return "", "", fmt.Errorf("remote url %s: %w", gitUrl, ErrInvalidUrl)
	}
	appName := strings.Join(urlParts[len(urlParts)-1:], "")
	return appName, gitUrl, nil
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
		ImageTag:   appVersion.String(),
		AppVersion: appVersion,
		RuntimeEnv: ctx.Stack(),
		// TODO stackref to collector project
		OTelEndpointUrl: url.URL{},
		ProjectUrl:      "https://" + repoUrl,
	}
	err = mergo.Merge(params, defParams)
	if err != nil {
		return fmt.Errorf("error filling app parameters: %w", err)
	}
	return nil
}

func DeployApp(ctx *pulumi.Context, params AppParms) error {
	mergeParams(ctx, &params)

	// Must match kind mount
	appCodeVolume := "/app-code"

	_, err := corev1.NewNamespace(ctx, "namespace", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(params.AppName),
			Namespace: pulumi.String(params.AppName),
			Labels: pulumi.StringMap{
				"app": pulumi.String(params.AppName),
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(params.AppName),
			Namespace: pulumi.String(params.AppName),
			Labels: pulumi.StringMap{
				"app": pulumi.String(params.AppName),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(params.AppName),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String(params.AppName),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Image: pulumi.String(params.AppName + ":" + params.ImageTag),
							Name:  pulumi.String(params.AppName),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8080),
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
					Port:       pulumi.Int(8080),
					Protocol:   pulumi.String("TCP"),
					TargetPort: pulumi.Any(8080),
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
