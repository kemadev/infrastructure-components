package complianceframework

import "strings"

// A ComplianceFramework represents a compliance framework, used to
// determine the compliance requirements for a given resource within the organization.
type ComplianceFramework string

const (
	// ComplianceFrameworkNone is the compliance framework to be used
	// when no compliance framework is required
	ComplianceFrameworkNone ComplianceFramework = "none"
	// ComplianceFrameworkRGPD is the compliance framework to be used
	// when the RGPD compliance framework is required
	ComplianceFrameworkRGPD ComplianceFramework = "rgpd"
)

// String returns the string representation of the ComplianceFramework
func (dc ComplianceFramework) String() string {
	return strings.ToLower(string(dc))
}
