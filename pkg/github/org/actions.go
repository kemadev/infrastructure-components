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
	// Internal workflows and actions
	"kemadev/workflows-and-actions/.github/workflows/*",
	"kemadev/workflows-and-actions/.github/actions/*",
	// Actions from reusable workflows and actions
	"actions/cache@*",
	"actions/checkout@*",
	"actions/download-artifact@*",
	"actions/github-script@*",
	"actions/labeler@*",
	"actions/setup-go@*",
	"actions/stale@*",
	"actions/upload-artifact@*",
	"anchore/sbom-action@*",
	"anchore/scan-action@*",
	"aws-actions/configure-aws-credentials@*",
	"DavidAnson/markdownlint-cli2-action@*",
	"docker://rhysd/actionlint:*",
	"golangci/golangci-lint-action@*",
	"googleapis/release-please-action@*",
	"goreleaser/goreleaser-action@*",
	"hadolint/hadolint-action@*",
	"ibiqlik/action-yamllint@*",
	"peter-evans/create-or-update-comment@*",
	"pulumi/actions@*",
	"semgrep/semgrep@*",
	"trufflesecurity/trufflehog@*",
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

func createActions(ctx *pulumi.Context, provider *github.Provider, args ActionsArgs) error {
	actionsOrganizationPermissionsName := util.FormatResourceName(
		ctx,
		"Actions organization permissions",
	)
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
	return nil
}
