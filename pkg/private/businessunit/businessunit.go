package businessunit

import "strings"

// A BusinessUnit represents a taskforce within the organization.
type BusinessUnit string

const (
	// BusinessUnitInfrastructure is the business unit in charge of infrastructure related tasks
	BusinessUnitInfrastructure BusinessUnit = "infrastructure"
	// BusinessUnitSecurity is the business unit in charge of security related tasks
	BusinessUnitSecurity BusinessUnit = "security"
	// BusinessUnitEngineering is the business unit in charge of engineering related tasks
	BusinessUnitEngineering BusinessUnit = "engineering"
	// BusinessUnitHumanResources is the business unit in charge of human-resources related tasks
	BusinessUnitHumanResources BusinessUnit = "human-resources"
	// BusinessUnitFinance is the business unit in charge of finance related tasks
	BusinessUnitFinance BusinessUnit = "finance"
	// BusinessUnitMarketing is the business unit in charge of marketing related tasks
	BusinessUnitMarketing BusinessUnit = "marketing"
	// BusinessUnitProduct is the business unit in charge of product related tasks
	BusinessUnitProduct BusinessUnit = "product"
	// BusinessUnitOperations is the business unit in charge of operations related tasks
	BusinessUnitOperations BusinessUnit = "operations"
	// BusinessUnitSales is the business unit in charge of sales related tasks
	BusinessUnitSales BusinessUnit = "sales"
	// BusinessUnitManagement is the business unit in charge of management related tasks
	BusinessUnitManagement BusinessUnit = "management"
	// BusinessUnitExecutive is the business unit in charge of executive related tasks
	BusinessUnitExecutive BusinessUnit = "executive"
	// BusinessUnitInternal is the business unit in charge of internal related tasks
	BusinessUnitInternal BusinessUnit = "internal"
)

// String returns the string representation of the BusinessUnit
func (bu BusinessUnit) String() string {
	return strings.ToLower(string(bu))
}
