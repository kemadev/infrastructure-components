package businessunit

import "strings"

// Company Business Unit (BU)
type BusinessUnit string

const (
	// Infrastructure related taskforce
	BusinessUnitInfrastructure BusinessUnit = "infrastructure"
	// Security related taskforce
	BusinessUnitSecurity BusinessUnit = "security"
	// Engineering related taskforce
	BusinessUnitEngineering BusinessUnit = "engineering"
	// Human-resources related taskforce
	BusinessUnitHumanResources BusinessUnit = "human-resources"
	// Finance related taskforce
	BusinessUnitFinance BusinessUnit = "finance"
	// Marketing related taskforce
	BusinessUnitMarketing BusinessUnit = "marketing"
	// Product related taskforce
	BusinessUnitProduct BusinessUnit = "product"
	// Operations related taskforce
	BusinessUnitOperations BusinessUnit = "operations"
	// Sales related taskforce
	BusinessUnitSales BusinessUnit = "sales"
	// Management related taskforce
	BusinessUnitManagement BusinessUnit = "management"
	// Executive related taskforce
	BusinessUnitExecutive BusinessUnit = "executive"
	// Internal related taskforce
	BusinessUnitInternal BusinessUnit = "internal"
)

func (bu BusinessUnit) String() string {
	return strings.ToLower(string(bu))
}
