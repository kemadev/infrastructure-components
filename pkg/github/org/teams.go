package org

import (
	"fmt"

	"github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type TeamMemberArgs struct {
	// Username is the GitHub username of the team member.
	Username string
	// Role is the role of the team member in the team. List of available roles can be found in the [documentation].
	//
	// [documentation]: https://www.pulumi.com/registry/packages/github/api-docs/membership/#state_role_go
	Role     string
}

type TeamArgs struct {
	// Name is the name of the team.
	Name        string
	// Description is a short description of the team.
	Description string
	// Privacy is the privacy setting of the team. List of available privacy settings can be found in the [documentation].
	//
	// [documentation]: https://www.pulumi.com/registry/packages/github/api-docs/team/#privacy_go
	Privacy     string
	// ParentTeam is the ID or slug of the parent team. If not set, the team will be a top-level team.
	ParentTeam  string
	// Members is a list of team members.
	Members     []TeamMemberArgs
}

type TeamsArgs struct {
	// Teams is a list of teams to create in the organization.
	Teams []TeamArgs
}

const (
	// AdminTeamName is the name of the team with full access everywhere.
	AdminTeamName       = "admins"
	// MaintainersTeamName is the name of the team that maintains permissions on all repositories.
	MaintainersTeamName = "maintainers"
	// DevelopersTeamName is the name of the parent team for all developers.
	DevelopersTeamName  = "developers"
)

var TeamsDefaultArgs = TeamsArgs{
	Teams: []TeamArgs{
		{
			Name:        AdminTeamName,
			Description: "Full access everywhere",
		},
		{
			Name:        MaintainersTeamName,
			Description: "Maintain permissions on all repositories",
		},
		{
			Name:        DevelopersTeamName,
			Description: "Parent team for all developers",
		},
	},
}

func createTeamsSetDefaults(args *TeamsArgs) {
	if args.Teams == nil {
		args.Teams = TeamsDefaultArgs.Teams
		return
	}
	if len(args.Teams) == 0 {
		args.Teams = TeamsDefaultArgs.Teams
		return
	}
	for _, team := range TeamsDefaultArgs.Teams {
		for i, t := range args.Teams {
			if t.Name == team.Name {
				if t.Description == "" {
					args.Teams[i].Description = team.Description
				}
				if t.Privacy == "" {
					args.Teams[i].Privacy = team.Privacy
				}
				if t.ParentTeam == "" {
					args.Teams[i].ParentTeam = team.ParentTeam
				}
				if t.Members == nil {
					args.Teams[i].Members = team.Members
				}
			}
		}
	}
}

func checkTeamMembersAreMembers(argsTeams TeamsArgs, argsMembers MembersArgs) error {
	for _, t := range argsTeams.Teams {
		if t.Members != nil {
			for _, m := range t.Members {
				found := false
				for _, member := range argsMembers.Members {
					if m.Username == member.Username {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("Team member %s in team %s is not also set to be an organization member", m.Username, t.Name)
				}
			}
		}
	}
	return nil
}

func createTeams(ctx *pulumi.Context, provider *github.Provider, argsTeams TeamsArgs, argsMembers MembersArgs) error {
	err := checkTeamMembersAreMembers(argsTeams, argsMembers)
	if err != nil {
		return err
	}
	for _, t := range argsTeams.Teams {
		teamName := util.FormatResourceName(ctx, "Team "+t.Name)
		team, err := github.NewTeam(ctx, teamName, &github.TeamArgs{
			Name:         pulumi.String(t.Name),
			Description:  pulumi.String(t.Description),
			Privacy:      pulumi.String("closed"),
			ParentTeamId: pulumi.String(t.ParentTeam),
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}
		teamSettingsName := util.FormatResourceName(ctx, "Team "+t.Name+" settings")
		_, err = github.NewTeamSettings(ctx, teamSettingsName, &github.TeamSettingsArgs{
			TeamId: team.ID(),
			ReviewRequestDelegation: &github.TeamSettingsReviewRequestDelegationArgs{
				MemberCount: pulumi.Int(1),
				Algorithm:   pulumi.String("LOAD_BALANCE"),
				Notify:      pulumi.Bool(true),
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}
		if t.Members != nil {
			teamMembersName := util.FormatResourceName(ctx, "Team "+t.Name+" members")
			_, err = github.NewTeamMembers(ctx, teamMembersName, &github.TeamMembersArgs{
				TeamId: team.ID(),
				Members: func() github.TeamMembersMemberArray {
					var members github.TeamMembersMemberArray
					for _, m := range t.Members {
						members = append(members, &github.TeamMembersMemberArgs{
							Username: pulumi.String(m.Username),
							Role:     pulumi.String(m.Role),
						})
					}
					return members
				}(),
			}, pulumi.Provider(provider))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
