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
	return
}
