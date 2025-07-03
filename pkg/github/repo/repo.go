package repo

import (
	"fmt"

	"github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DirectMember struct {
	// Username is the GitHub username of the user to add as a direct member of the repository.
	Username string
	// Role is the role of the user in the repository. List of available roles can be found in the [documentation].
	//
	// [documentation]: https://www.pulumi.com/registry/packages/github/api-docs/repositorycollaborators/#permission_go
	Role string
}

type Team struct {
	// Name is the name of the team to add as a collaborator to the repository.
	Name string
	// Role is the role of the team in the repository. List of available roles can be found in the [documentation].
	//
	// [documentation]: https://www.pulumi.com/registry/packages/github/api-docs/repositorycollaborators/#permission_go
	Role string
}

type RepositoryArgs struct {
	// Name is the name of the repository.
	Name string
	// Description is a short description of the repository.
	Description string
	// HomepageUrl is the URL of the repository's homepage.
	HomepageUrl string
	// defaultBranch is the name of the default branch for the repository.
	defaultBranch string
	// Topics is a list of topics to associate with the repository.
	Topics []string
	// Visibility is the visibility of the repository. Defaults to "private". List of available visibility settings can be found in the [documentation].
	//
	// [documentation]: https://www.pulumi.com/registry/packages/github/api-docs/repository/#visibility_go
	Visibility string
	// Archived indicates whether the repository should be archived. Defaults to false.
	Archived bool
	// IsTemplate indicates whether the repository is a template repository. Defaults to false.
	IsTemplate bool
	// Teams is a list of teams to add as collaborators to the repository.
	Teams []Team
	// DirectMembers is a list of direct members to add to the repository with specific roles.
	DirectMembers []DirectMember
}

var RepositoryDefaultArgs = RepositoryArgs{
	Visibility:    "private",
	Archived:      false,
	IsTemplate:    false,
	defaultBranch: "main",
}

func createRepositorySetDefaults(args *RepositoryArgs) error {
	if args.Description == "" {
		return fmt.Errorf("Repository Description is required")
	} else if args.Description == "CHANGEME" {
		return fmt.Errorf("Repository Description must be changed from the default value")
	}
	if args.Visibility == "" {
		args.Visibility = RepositoryDefaultArgs.Visibility
	} else if args.Visibility == "CHANGEME" {
		return fmt.Errorf("Repository Visibility must be changed from the default value")
	}
	if args.defaultBranch == "" {
		args.defaultBranch = RepositoryDefaultArgs.defaultBranch
	}
	return nil
}

func createRepo(
	ctx *pulumi.Context,
	provider *github.Provider,
	argsRepo RepositoryArgs,
	suffix string,
) (*github.Repository, error) {
	repoName := util.FormatResourceName(ctx, "Repository"+suffix)
	repo, err := github.NewRepository(ctx, repoName, &github.RepositoryArgs{
		Name:        pulumi.String(argsRepo.Name),
		Description: pulumi.String(argsRepo.Description),
		HomepageUrl: pulumi.String(argsRepo.HomepageUrl),
		Topics: func() pulumi.StringArrayInput {
			var topics pulumi.StringArray
			for _, topic := range argsRepo.Topics {
				topics = append(topics, pulumi.String(topic))
			}
			return topics
		}(),
		Visibility: pulumi.String(argsRepo.Visibility),
		IsTemplate: pulumi.Bool(argsRepo.IsTemplate),

		// Prevent accidental deletion
		ArchiveOnDestroy: pulumi.Bool(true),
		// Allow non-admins read access from pulumi
		IgnoreVulnerabilityAlertsDuringRead: pulumi.Bool(true),

		AllowSquashMerge:         pulumi.Bool(true),
		SquashMergeCommitTitle:   pulumi.String("PR_TITLE"),
		SquashMergeCommitMessage: pulumi.String("PR_BODY"),
		AllowMergeCommit:         pulumi.Bool(false),
		AllowRebaseMerge:         pulumi.Bool(false),
		AllowUpdateBranch:        pulumi.Bool(true),
		AllowAutoMerge:           pulumi.Bool(true),
		DeleteBranchOnMerge:      pulumi.Bool(true),
		HasDiscussions:           pulumi.Bool(true),
		HasIssues:                pulumi.Bool(true),
		HasProjects:              pulumi.Bool(true),
		HasWiki:                  pulumi.Bool(true),
		HasDownloads:             pulumi.Bool(false),
		Archived:                 pulumi.Bool(argsRepo.Archived),
		WebCommitSignoffRequired: pulumi.Bool(false),
		AutoInit:                 pulumi.Bool(true),

		VulnerabilityAlerts: func() pulumi.Bool {
			if argsRepo.Visibility == "public" {
				return pulumi.Bool(true)
			}
			// Advanced Security is required for private repositories
			return pulumi.Bool(false)
		}(),
		SecurityAndAnalysis: func() *github.RepositorySecurityAndAnalysisArgs {
			if argsRepo.Visibility == "public" {
				return &github.RepositorySecurityAndAnalysisArgs{
					SecretScanning: github.RepositorySecurityAndAnalysisSecretScanningArgs{
						Status: pulumi.String("enabled"),
					},
					SecretScanningPushProtection: github.RepositorySecurityAndAnalysisSecretScanningPushProtectionArgs{
						Status: pulumi.String("enabled"),
					},
				}
			}
			// Advanced Security is required for private repositories
			return nil
		}(),
	}, pulumi.Provider(provider), pulumi.IgnoreChanges([]string{"template"}))
	if err != nil {
		return nil, err
	}
	defaulBranchName := util.FormatResourceName(ctx, "Repository default branch"+suffix)
	_, err = github.NewBranch(ctx, defaulBranchName, &github.BranchArgs{
		Repository: repo.Name,
		Branch:     pulumi.String(argsRepo.defaultBranch),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}
	branchDefaultName := util.FormatResourceName(ctx, "Repository default branch"+suffix)
	_, err = github.NewBranchDefault(ctx, branchDefaultName, &github.BranchDefaultArgs{
		Repository: repo.Name,
		Branch:     pulumi.String(argsRepo.defaultBranch),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	repoCollaboratorsName := util.FormatResourceName(ctx, "Repository collaborators"+suffix)
	_, err = github.NewRepositoryCollaborators(
		ctx,
		repoCollaboratorsName,
		&github.RepositoryCollaboratorsArgs{
			Repository: repo.Name,
			Users: func() github.RepositoryCollaboratorsUserArray {
				var members github.RepositoryCollaboratorsUserArray
				for _, m := range argsRepo.DirectMembers {
					members = append(members, &github.RepositoryCollaboratorsUserArgs{
						Username:   pulumi.String(m.Username),
						Permission: pulumi.String(m.Role),
					})
				}
				return members
			}(),
			Teams: func() github.RepositoryCollaboratorsTeamArray {
				var teams github.RepositoryCollaboratorsTeamArray
				for _, t := range argsRepo.Teams {
					teams = append(teams, &github.RepositoryCollaboratorsTeamArgs{
						TeamId:     pulumi.String(t.Name),
						Permission: pulumi.String(t.Role),
					})
				}
				return teams
			}(),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
