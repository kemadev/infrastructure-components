package dataclassification

import "strings"

// Company Business Unit (BU)
type DataClassification string

const (
	// Infrastructure related taskforce
	DataClassificationNone DataClassification = "none"
)

func (dc DataClassification) String() string {
	return strings.ToLower(string(dc))
}
