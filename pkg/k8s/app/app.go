package app

import (
	"net/url"

	"github.com/blang/semver"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AppParms struct {
	ImageRef string
	ImageTag string
	RuntimeEnv string
	OTelEndpointUrl url.URL
	AppVersion semver.Version
	AppName string
	AppNamespace string
	BusinessUnitId 
	CustomerId
	CostCenter
	CostAllocationOwner
	OperationsOwner
	Rpo
	DataClassification
	ComplianceFramework
	Expiration
	ProjectUrl
	MonitoringUrl
}

func DeployDevApp(ctx *pulumi.Context)
