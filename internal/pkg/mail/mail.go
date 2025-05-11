package mail

import "vcs.kema.run/kema/infrastructure-components/internal/pkg/domain"

type (
	MailDomain  string
	MailAddress string
)

func (md MailDomain) String() string {
	return string(md)
}

func (ma MailAddress) String() string {
	return string(ma)
}

var PrimaryMailDomain MailDomain = MailDomain(domain.DomainKemaDotDev.String())
