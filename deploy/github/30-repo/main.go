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
			Name:        "server-bootstrap",
			Description: "Server boostrapping, from PXE to Ignition",
			Topics:      []string{"server", "bootstrap", "pxe", "ignition", "bare-metal"},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "repo-template",
			Description: "Repository template",
			Topics:      []string{"repository", "template", "github", "pulumi", "go", "copier"},
		},
	},
	{
		Repository: repo.RepositoryArgs{
			Name:        "go-framework",
			Description: "Go framework, ensuring best practices and security",
			Topics:      []string{"go", "framework", "best-practices", "security"},
		},
	},
	// NOTE This one has been initially imported using `pulumi import 'github:index/repository:Repository' '<resource id>' <repo name> --provider '<provider urn>'`
	{
		Repository: repo.RepositoryArgs{
			Name:        "infrastructure-components",
			Description: "Infrastructure components, ensuring homegenous and performant standards",
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
