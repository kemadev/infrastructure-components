package main

import (
	"net/url"
	"time"

<<<<<<< before updating
	"github.com/kemadev/infrastructure-components/pkg/k8s/basichttpapp"
||||||| last update
=======
	"github.com/kemadev/infrastructure-components/pkg/k8s/basichttpapp"
	"github.com/kemadev/infrastructure-components/pkg/private/businessunit"
	"github.com/kemadev/infrastructure-components/pkg/private/costcenter"
	"github.com/kemadev/infrastructure-components/pkg/private/customer"
>>>>>>> after updating
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := basichttpapp.DeployBasicHTTPApp(ctx, basichttpapp.AppParms{
			// TODO
			AppNamespace:        "changeme",
			AppComponent:        "changeme",
			BusinessUnitId:      businessunit.BusinessUnitEngineering,
			CustomerId:          customer.CustomerInternal,
			CostCenter:          costcenter.CostCenterInternal,
			CostAllocationOwner: businessunit.BusinessUnitEngineering,
			OperationsOwner:     businessunit.BusinessUnitEngineering,
			Rpo:                 1 * time.Hour,
			MonitoringUrl: url.URL{
				Scheme: "https",
				Host:   "changeme",
				Path:   "changeme",
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
