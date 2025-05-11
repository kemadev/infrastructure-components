package domain

import "strings"

type Domain string

const (
	DomainKemaDotDev      Domain = "kema.dev"
	DomainKemaDotCloud    Domain = "kema.cloud"
	DomainKemaDotRun      Domain = "kema.run"
	DomainDoubleJDotFr    Domain = "doublej.fr"
	DomainKemaDevDotFr    Domain = "kemadev.fr"
	DomainKemaDevDotCom   Domain = "kemadev.com"
	DomainKemaDotInternal Domain = "kema.internal"
)

func (d Domain) String() string {
	return strings.ToLower(string(d))
}

var (
	// Domains registered to AWS registar
	AwsRegisteredDomain []Domain = []Domain{
		DomainKemaDevDotCom,
		DomainKemaDevDotFr,
		DomainDoubleJDotFr,
		DomainKemaDotRun,
	}
	// Domains registered to Squarespace registar
	SquarespaceRegisteredDomain []Domain = []Domain{
		DomainKemaDotDev,
	}
	// Domains registered to Cloudflare registar
	CloudflareRegisteredDomain []Domain = []Domain{
		DomainKemaDotCloud,
	}
	// Internal (non registered) domains
	InternalDomain []Domain = []Domain{
		DomainKemaDotInternal,
	}
)
