package dataclassification

import "strings"

// A DataClassification represents a data classification, used to
// determine the data classification requirements for a given resource within the organization.
type DataClassification string

const (
	// DataClassificationNone is the data classification to be used
	// when no data classification is required.
	DataClassificationNone DataClassification = "none"
)

// String returns the string representation of the DataClassification.
func (dc DataClassification) String() string {
	return strings.ToLower(string(dc))
}
