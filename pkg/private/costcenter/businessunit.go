package costcenter

import (
	"strings"

	"github.com/kemadev/infrastructure-components/pkg/private/businessunit"
)

// Cost center
type CostCenter string

const (
	// Infrastructure related taskforce
	CostCenterInfrastructure CostCenter = CostCenter(businessunit.BusinessUnitInfrastructure)
	// Security related taskforce
	CostCenterSecurity CostCenter = CostCenter(businessunit.BusinessUnitSecurity)
	// Engineering related taskforce
	CostCenterEngineering CostCenter = CostCenter(businessunit.BusinessUnitEngineering)
	// Human-resources related taskforce
	CostCenterHumanResources CostCenter = CostCenter(businessunit.BusinessUnitHumanResources)
	// Finance related taskforce
	CostCenterFinance CostCenter = CostCenter(businessunit.BusinessUnitFinance)
	// Marketing related taskforce
	CostCenterMarketing CostCenter = CostCenter(businessunit.BusinessUnitMarketing)
	// Product related taskforce
	CostCenterProduct CostCenter = CostCenter(businessunit.BusinessUnitProduct)
	// Operations related taskforce
	CostCenterOperations CostCenter = CostCenter(businessunit.BusinessUnitOperations)
	// Sales related taskforce
	CostCenterSales CostCenter = CostCenter(businessunit.BusinessUnitSales)
	// Management related taskforce
	CostCenterManagement CostCenter = CostCenter(businessunit.BusinessUnitManagement)
	// Executive related taskforce
	CostCenterExecutive CostCenter = CostCenter(businessunit.BusinessUnitExecutive)
	// Internal related taskforce
	CostCenterInternal CostCenter = CostCenter(businessunit.BusinessUnitInternal)
)

func (cc CostCenter) String() string {
	return strings.ToLower(string(cc))
}
