package repo

import (
	p "github.com/kemadev/infrastructure-components/pkg/github/provider"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type WrapperArgs struct {
	// ProviderOpts contains the provider settings for the GitHub repository. It can be omitted if the provider is already is passed as "Provider" in arguments.
	ProviderOpts p.ProviderArgs
	// Envs contains the settings for environments in the repository.
	Envs EnvsArgs
	// Rulesets contains the settings for rulesets in the repository.
	Rulesets RulesetsArgs
	// Repository contains the settings for the repository itself.
	Repository RepositoryArgs
	// Provider is the GitHub provider to use for the repository. If provided, ProviderOpts will be ignored.
	Provider *github.Provider
	// GitHubPlan is the GitHub plan subscribed for the organization. It is used to determine whether to create resources for paid features. Default to "free".
	GitHubPlan string
}

func setDefaultArgs(args *WrapperArgs) error {
	p.SetDefaults(&args.ProviderOpts)
	createEnvironmentsSetDefaults(&args.Envs)
	createRulesetsSetDefaults(&args.Rulesets)
	err := createRepositorySetDefaults(&args.Repository)
	if err != nil {
		return err
	}
	return nil
}

// Wrapper creates a GitHub repository with the provided settings, environments, rulesets, codeowners, and files.
func Wrapper(ctx *pulumi.Context, args WrapperArgs) error {
	if args.GitHubPlan == "" {
		args.GitHubPlan = "free"
	}
	enablePaidFeatures := args.GitHubPlan != "free" || args.Repository.Visibility == "public"
	// targetBranch := "repo-as-code-update"
	err := setDefaultArgs(&args)
	if err != nil {
		return err
	}
	var provider *github.Provider
	if args.Provider != nil {
		provider = args.Provider
	} else {
		prov, err := p.NewProvider(ctx, args.ProviderOpts)
		if err != nil {
			return err
		}
		provider = prov
	}
	suffix := " " + args.Repository.Name
	repo, err := createRepo(ctx, provider, args.Repository, enablePaidFeatures, suffix)
	if err != nil {
		return err
	}
	_, err = createEnvironments(ctx, provider, repo, args.Envs, suffix)
	if err != nil {
		return err
	}
	if enablePaidFeatures {
		err = createRulesets(ctx, provider, repo, args.Rulesets, args.Envs, suffix)
		if err != nil {
			return err
		}
	}
	err = createDependabot(ctx, provider, repo, suffix)
	if err != nil {
		return err
	}
	err = createIssues(ctx, provider, repo, suffix)
	if err != nil {
		return err
	}
	return nil
}
