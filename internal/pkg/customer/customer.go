package customer

import "strings"

type Customer string

const (
	CustomerInternal Customer = "internal"
)

func (c Customer) String() string {
	return strings.ToLower(string(c))
}
