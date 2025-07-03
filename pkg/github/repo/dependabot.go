package repo

import (
	"github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createDependabot(
	ctx *pulumi.Context,
	provider *github.Provider,
	repo *github.Repository,
	suffix string,
) error {
	repoDependabotSecurityUpdateName := util.FormatResourceName(
		ctx,
		"Dependabot security updates"+suffix,
	)
	_, err := github.NewRepositoryDependabotSecurityUpdates(
		ctx,
		repoDependabotSecurityUpdateName,
		&github.RepositoryDependabotSecurityUpdatesArgs{
			Repository: repo.Name,
			Enabled:    pulumi.Bool(true),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return err
	}
	return nil
}
