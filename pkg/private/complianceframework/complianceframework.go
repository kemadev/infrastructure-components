package complianceframework

import "strings"

// Company Business Unit (BU)
type ComplianceFramework string

const (
	// Infrastructure related taskforce
	ComplianceFrameworkNone ComplianceFramework = "none"
	ComplianceFrameworkRGPD ComplianceFramework = "rgpd"
)

func (dc ComplianceFramework) String() string {
	return strings.ToLower(string(dc))
}
