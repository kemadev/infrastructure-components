package provider

import (
	"github.com/kemadev/infrastructure-components/pkg/util"
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ProviderArgs struct {
	// Owner is the GitHub organization owner (i.e. organization name).
	Owner string
}

var ProviderDefaultArgs = ProviderArgs{
	Owner: "kemadev",
}

// SetDefaults sets the default values for the provider arguments.
func SetDefaults(args *ProviderArgs) {
	if args.Owner == "" {
		args.Owner = ProviderDefaultArgs.Owner
	}
}

// NewProvider creates a new GitHub provider with the specified owner.
func NewProvider(ctx *pulumi.Context, args ProviderArgs) (*github.Provider, error) {
	providerName := util.FormatResourceName(ctx, "Provider")
	provider, err := github.NewProvider(ctx, providerName, &github.ProviderArgs{
		Owner: pulumi.String(args.Owner),
	})
	if err != nil {
		return nil, err
	}
	return provider, nil
}
