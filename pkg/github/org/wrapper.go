package org

import (
	p "github.com/kemadev/infrastructure-components/pkg/github/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type WrapperArgs struct {
	// Provider is the GitHub provider configuration.
	Provider p.ProviderArgs
	// Settings contains the settings for the organization.
	Settings SettingsArgs
	// Teams contains the teams to create in the organization.
	Teams TeamsArgs
	// Actions contains the GitHub Actions patterns that are allowed to run in the organization.
	Actions ActionsArgs
	// Members contains the members to add to the organization.
	Members MembersArgs
	// GitHubPlan is the GitHub plan subscribed for the organization. It is used to determine whether to create resources for paid features. Default to "free".
	GitHubPlan string
}

func setDefaultArgs(args *WrapperArgs) {
	p.SetDefaults(&args.Provider)
	createTeamsSetDefaults(&args.Teams)
	createActionsSetDefaults(&args.Actions)
}

// Wrapper creates a GitHub organization with the provided settings, members, teams, and actions.
func Wrapper(ctx *pulumi.Context, args WrapperArgs) error {
	if args.GitHubPlan == "" {
		args.GitHubPlan = "free"
	}
	enablePaidFeatures := args.GitHubPlan != "free"
	setDefaultArgs(&args)
	provider, err := p.NewProvider(ctx, args.Provider)
	if err != nil {
		return err
	}
	err = createSettings(ctx, provider, args.Settings)
	if err != nil {
		return err
	}
	err = createMembers(ctx, provider, args.Members)
	if err != nil {
		return err
	}
	err = createTeams(ctx, provider, args.Teams, args.Members)
	if err != nil {
		return err
	}
	err = createActions(ctx, provider, args.Actions, enablePaidFeatures)
	if err != nil {
		return err
	}
	return nil
}
