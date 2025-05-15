package customer

import "strings"

// A Customer represents a customer.
type Customer string

const (
	// CustomerNone is the customer to be used for internal resources.
	CustomerInternal Customer = "internal"
)

// String returns the string representation of the Customer.
func (c Customer) String() string {
	return strings.ToLower(string(c))
}
