package main

import (
	"net/url"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"vcs.kema.run/kema/infrastructure-components/pkg/k8s/basichttpapp"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := basichttpapp.DeployBasicHTTPApp(ctx, basichttpapp.AppParms{
			// TODO
			AppNamespace:        "input",
			BusinessUnitId:      "input",
			CustomerId:          "input",
			CostCenter:          "input",
			CostAllocationOwner: "input",
			OperationsOwner:     "input",
			Rpo:                 1 * time.Hour,
			DataClassification:  "input",
			ComplianceFramework: "input",
			Expiration:          time.Time{},
			MonitoringUrl:       url.URL{},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
