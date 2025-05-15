package mail

import (
	"strings"

	"github.com/kemadev/infrastructure-components/pkg/private/domain"
)

type (
	// A MailDomain represents a mail domain name.
	MailDomain string
	// A MailAddress represents a mail address.
	MailAddress string
)

// String returns the string representation of the MailDomain
func (md MailDomain) String() string {
	return strings.ToLower(string(md))
}

// String returns the string representation of the MailAddress
func (ma MailAddress) String() string {
	return strings.ToLower(string(ma))
}

// PrimaryMailDomain is the primary mail domain used in the organization.
var PrimaryMailDomain MailDomain = MailDomain(domain.DomainKemaDotDev.String())
