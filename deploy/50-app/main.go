package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kemadev/infrastructure-components/pkg/k8s/basichttpapp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := basichttpapp.DeployBasicHTTPApp(ctx, basichttpapp.AppParms{
			AppNamespace:        "changeme",
			AppComponent:        "changeme",
			BusinessUnitId:      "changeme",
			CustomerId:          "changeme",
			CostCenter:          "changeme",
			CostAllocationOwner: "changeme",
			OperationsOwner:     "changeme",
			Rpo:                 0 * time.Second,
			MonitoringUrl: url.URL{
				Scheme: "https",
				Host:   "changeme",
				Path:   "changeme",
			},
		})
		if err != nil {
<<<<<<< before updating
			return fmt.Errorf("error deploying basic HTTP app: %w", err)
||||||| last update
			return err
=======
			return fmt.Errorf("failed to deploy basic HTTP app: %w", err)
>>>>>>> after updating
		}
		return nil
	})
}
