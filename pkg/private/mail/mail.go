package mail

import (
	"strings"

	"github.com/kemadev/infrastructure-components/pkg/private/domain"
)

type (
	MailDomain  string
	MailAddress string
)

func (md MailDomain) String() string {
	return strings.ToLower(string(md))
}

func (ma MailAddress) String() string {
	return strings.ToLower(string(ma))
}

var PrimaryMailDomain MailDomain = MailDomain(domain.DomainKemaDotDev.String())
