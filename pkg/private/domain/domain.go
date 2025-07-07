
/*
Package domain provides a representation of a DNS domain name.

It provides static definitions for the domains used in the organization,
and a way to represent them as strings.
*/
package domain

import "strings"

// A Domain represents a DNS domain name.
type Domain string

const (
	// DomainKemaDotFr is the domain name for kema.dev
	// DomainKemaDotDev is the domaine name for kema.dev
	DomainKemaDotDev      Domain = "kema.dev"
	// DomainKemaDotCloud is the domaine name for kema.cloud
	DomainKemaDotCloud    Domain = "kema.cloud"
	// DomainKemaDotRun is the domaine name for kema.run
	DomainKemaDotRun      Domain = "kema.run"
	// DomainDoubleJDotFr is the domaine name for doublej.fr
	DomainDoubleJDotFr    Domain = "doublej.fr"
	// DomainKemaDevDotFr is the domaine name for kemadev.fr
	DomainKemaDevDotFr    Domain = "kemadev.fr"
	// DomainKemaDevDotCom is the domaine name for kemadev.com
	DomainKemaDevDotCom   Domain = "kemadev.com"
	// DomainKemaDotInternal is the domaine name for kema.internal
	DomainKemaDotInternal Domain = "kema.internal"
)

// String returns the string representation of the Domain
func (d Domain) String() string {
	return strings.ToLower(string(d))
}

var (
	// AwsRegisteredDomain are the domains registered to the AWS registar.
	AwsRegisteredDomain []Domain = []Domain{
		DomainKemaDevDotCom,
		DomainKemaDevDotFr,
		DomainDoubleJDotFr,
		DomainKemaDotRun,
	}
	// SquarespaceRegisteredDomain are the domains registered to the Squarespace registar.
	SquarespaceRegisteredDomain []Domain = []Domain{
		DomainKemaDotDev,
	}
	// CloudflareRegisteredDomain are the domains registered to the Cloudflare registar.
	CloudflareRegisteredDomain []Domain = []Domain{
		DomainKemaDotCloud,
	}
	// InternalDomain are the domains that are not registered to any registar, but used internally.
	InternalDomain []Domain = []Domain{
		DomainKemaDotInternal,
	}
)
