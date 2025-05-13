package mail

import (
	"strings"

	"github.com/kemadev/infrastructure-components/internal/pkg/domain"
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
