package repo

import (
	"github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type EnvsArgs struct {
	// Dev is the name of the development environment.
	Dev string
	// Next is the name of the next environment, typically used for staging or pre-production.
	Next string
	// Prod is the name of the production environment.
	Prod string
}

type TEnvironmentsCreated struct {
	dev  *github.RepositoryEnvironment
	next *github.RepositoryEnvironment
	prod *github.RepositoryEnvironment
}

var EnvsDefaultArgs = EnvsArgs{
	Dev:  "dev",
	Next: "next",
	Prod: "main",
}

func createEnvironmentsSetDefaults(args *EnvsArgs) {
	if args.Dev == "" {
		args.Dev = EnvsDefaultArgs.Dev
	}
	if args.Next == "" {
		args.Next = EnvsDefaultArgs.Next
	}
	if args.Prod == "" {
		args.Prod = EnvsDefaultArgs.Prod
	}
}

func createEnvironments(
	ctx *pulumi.Context,
	provider *github.Provider,
	repo *github.Repository,
	argsEnvs EnvsArgs,
	suffix string,
) (TEnvironmentsCreated, error) {
	deploymentEnvironmentDevName := util.FormatResourceName(
		ctx,
		"Deployment environment dev"+suffix,
	)
	deploymentEnvironmentDev, err := github.NewRepositoryEnvironment(
		ctx,
		deploymentEnvironmentDevName,
		&github.RepositoryEnvironmentArgs{
			Repository:        repo.Name,
			Environment:       pulumi.String(argsEnvs.Dev),
			CanAdminsBypass:   pulumi.Bool(true),
			PreventSelfReview: pulumi.Bool(false),
			WaitTimer:         pulumi.Int(0),
			// Trust PR reviews
			Reviewers: github.RepositoryEnvironmentReviewerArray{},
			DeploymentBranchPolicy: github.RepositoryEnvironmentDeploymentBranchPolicyArgs{
				CustomBranchPolicies: pulumi.Bool(true),
				ProtectedBranches:    pulumi.Bool(false),
			},
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return TEnvironmentsCreated{}, err
	}

	deploymentEnvironmentNextName := util.FormatResourceName(
		ctx,
		"Deployment environment next"+suffix,
	)
	deploymentEnvironmentNext, err := github.NewRepositoryEnvironment(
		ctx,
		deploymentEnvironmentNextName,
		&github.RepositoryEnvironmentArgs{
			Repository:        repo.Name,
			Environment:       pulumi.String(argsEnvs.Next),
			CanAdminsBypass:   pulumi.Bool(true),
			PreventSelfReview: pulumi.Bool(false),
			WaitTimer:         pulumi.Int(0),
			// Trust PR reviews
			Reviewers: github.RepositoryEnvironmentReviewerArray{},
			DeploymentBranchPolicy: github.RepositoryEnvironmentDeploymentBranchPolicyArgs{
				CustomBranchPolicies: pulumi.Bool(true),
				ProtectedBranches:    pulumi.Bool(false),
			},
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return TEnvironmentsCreated{}, err
	}

	deploymentEnvironmentProdName := util.FormatResourceName(
		ctx,
		"Deployment environment prod"+suffix,
	)
	deploymentEnvironmentProd, err := github.NewRepositoryEnvironment(
		ctx,
		deploymentEnvironmentProdName,
		&github.RepositoryEnvironmentArgs{
			Repository:        repo.Name,
			Environment:       pulumi.String(argsEnvs.Prod),
			CanAdminsBypass:   pulumi.Bool(true),
			PreventSelfReview: pulumi.Bool(false),
			WaitTimer:         pulumi.Int(0),
			// Trust PR reviews
			Reviewers: github.RepositoryEnvironmentReviewerArray{},
			DeploymentBranchPolicy: github.RepositoryEnvironmentDeploymentBranchPolicyArgs{
				CustomBranchPolicies: pulumi.Bool(true),
				ProtectedBranches:    pulumi.Bool(false),
			},
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return TEnvironmentsCreated{}, err
	}

	repositoryEnvironmentDeploymentPolicyDevName := util.FormatResourceName(
		ctx,
		"Repository environment deployment policy dev"+suffix)
	_, err = github.NewRepositoryEnvironmentDeploymentPolicy(
		ctx,
		repositoryEnvironmentDeploymentPolicyDevName,
		&github.RepositoryEnvironmentDeploymentPolicyArgs{
			Repository:    repo.Name,
			Environment:   deploymentEnvironmentDev.Environment,
			BranchPattern: pulumi.String("*"),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return TEnvironmentsCreated{}, err
	}

	repositoryEnvironmentDeploymentPolicyNextName := util.FormatResourceName(
		ctx,
		"Repository environment deployment policy next"+suffix)
	_, err = github.NewRepositoryEnvironmentDeploymentPolicy(
		ctx,
		repositoryEnvironmentDeploymentPolicyNextName,
		&github.RepositoryEnvironmentDeploymentPolicyArgs{
			Repository:    repo.Name,
			Environment:   deploymentEnvironmentNext.Environment,
			BranchPattern: pulumi.String("*"),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return TEnvironmentsCreated{}, err
	}

	repositoryEnvironmentDeploymentPolicyProdName := util.FormatResourceName(
		ctx,
		"Repository environment deployment policy prod"+suffix)
	_, err = github.NewRepositoryEnvironmentDeploymentPolicy(
		ctx,
		repositoryEnvironmentDeploymentPolicyProdName,
		&github.RepositoryEnvironmentDeploymentPolicyArgs{
			Repository:    repo.Name,
			Environment:   deploymentEnvironmentProd.Environment,
			BranchPattern: pulumi.String("*"),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return TEnvironmentsCreated{}, err
	}

	return TEnvironmentsCreated{
		dev:  deploymentEnvironmentDev,
		next: deploymentEnvironmentNext,
		prod: deploymentEnvironmentProd,
	}, nil
}
