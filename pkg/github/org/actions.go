package org

import (
	"github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ActionsArgs struct {
	// Actions is a list of GitHub Actions patterns that are allowed to run in the organization.
	Actions []string
}

// ActionsDefaultActions is the default list of GitHub Actions patterns that are allowed to run in the organization.
var ActionsDefaultActions = []string{
	"docker://kemadev/*",
	"actions/checkout@*",
}

func createActionsSetDefaults(args *ActionsArgs) {
	if args.Actions == nil {
		args.Actions = ActionsDefaultActions
		return
	}
	if len(args.Actions) == 0 {
		args.Actions = ActionsDefaultActions
		return
	}
	args.Actions = append(args.Actions, ActionsDefaultActions...)
}

func createActions(ctx *pulumi.Context, provider *github.Provider, args ActionsArgs, enablePaidFeatures bool) error {
	actionsOrganizationPermissionsName := util.FormatResourceName(
		ctx,
		"Actions organization permissions",
	)
	if !enablePaidFeatures {
		_, err := github.NewActionsOrganizationPermissions(
			ctx,
			actionsOrganizationPermissionsName,
			&github.ActionsOrganizationPermissionsArgs{
				AllowedActions:      pulumi.String("all"),
				EnabledRepositories: pulumi.String("all"),
			},
			pulumi.Provider(provider),
		)
		if err != nil {
			return err
		}
	} else {
		_, err := github.NewActionsOrganizationPermissions(
			ctx,
			actionsOrganizationPermissionsName,
			&github.ActionsOrganizationPermissionsArgs{
				AllowedActions:      pulumi.String("selected"),
				EnabledRepositories: pulumi.String("all"),
				AllowedActionsConfig: &github.ActionsOrganizationPermissionsAllowedActionsConfigArgs{
					GithubOwnedAllowed: pulumi.Bool(false),
					VerifiedAllowed:    pulumi.Bool(false),
					PatternsAlloweds: func() pulumi.StringArray {
						var patterns pulumi.StringArray
						for _, action := range args.Actions {
							patterns = append(patterns, pulumi.String(action))
						}
						return patterns
					}(),
				},
			},
			pulumi.Provider(provider),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
