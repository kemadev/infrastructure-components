package org

import (
	"github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type User struct {
	// Username is the GitHub username of the user.
	Username string
	// Role is the role of the user in the organization. List of available roles can be found in the [documentation].
	//
	// [documentation]: https://www.pulumi.com/registry/packages/github/api-docs/membership/#state_role_go
	Role string
}

type MembersArgs struct {
	// Members is a list of users to add to the organization.
	Members []User
	// Admins is a list of users to add as admins to the organization.
	Admins []User
}

func createMembers(ctx *pulumi.Context, provider *github.Provider, argsMembers MembersArgs) error {
	for _, t := range argsMembers.Members {
		memberName := util.FormatResourceName(ctx, "Member")
		_, err := github.NewMembership(ctx, memberName, &github.MembershipArgs{
			Username: pulumi.String(t.Username),
			Role:     pulumi.String(t.Role),
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}
	}
	return nil
}
