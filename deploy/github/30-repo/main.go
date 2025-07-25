package main

import (
	"github.com/kemadev/infrastructure-components/pkg/github/provider"
	p "github.com/kemadev/infrastructure-components/pkg/github/provider"
	"github.com/kemadev/infrastructure-components/pkg/github/repo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var repositories = []repo.WrapperArgs{
	{
		Repository: repo.RepositoryArgs{
			Name:        ".github",
			Description: "Organization wide files",
			Visibility:  "public",
			HomepageUrl: "https://www.kema.dev",
			Topics:      []string{"github", "organization", "files"},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "discussions",
			Description: "Organization wide discussions",
			Visibility:  "public",
			HomepageUrl: "https://www.kema.dev",
			Topics:      []string{"github", "organization", "discussions"},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "repo-template",
			Description: "Repository template",
			Visibility:  "public",
			Topics:      []string{"repository", "template", "github", "pulumi", "go", "copier"},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "go-framework",
			Description: "Go framework, ensuring best practices and security",
			Visibility:  "public",
			Topics:      []string{"go", "framework", "best-practices", "security"},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "ci-cd",
			Description: "CI/CD tooling for repositories",
			Visibility:  "public",
			Topics:      []string{"go", "ci-cd", "github", "pulumi", "docker", "runner"},
		},
	},
	// NOTE This one has been initially imported using `pulumi import 'github:index/repository:Repository' '<resource id>' <repo name> --provider '<provider urn>'`
	{
		Repository: repo.RepositoryArgs{
			Name:        "infrastructure-components",
			Description: "Infrastructure components, ensuring homegenous and performant standards",
			Visibility:  "public",
			Topics: []string{
				"go",
				"pulumi",
				"infrastructure",
				"components",
				"kubernetes",
				"security",
			},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "kemutil",
			Description: "CLI utility for everyday tasks",
			Visibility:  "public",
			Topics:      []string{"go", "framework", "cli", "utility", "everyday-tasks"},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "server-bootstrap",
			Description: "Server boostrapping, from PXE to Ignition",
			Visibility:  "private",
			Topics:      []string{"server", "bootstrap", "pxe", "ignition", "bare-metal"},
		},
	},
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		provider, err := p.NewProvider(ctx, provider.ProviderArgs{
			Owner: "kemadev",
		})
		if err != nil {
			return err
		}
		for _, repoArgs := range repositories {
			repoArgs.Provider = provider
			err := repo.Wrapper(ctx, repoArgs)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
