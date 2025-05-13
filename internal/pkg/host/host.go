package host

import (
	"net/url"
	"strings"

	"github.com/kemadev/infrastructure-components/internal/pkg/domain"
)

type (
	// A base host, used to construct URLs
	Host url.URL
	// An URL template, used to distribute traffic
	URL struct {
		// Host where traffic is destinated
		BaseHost Host
		// Path pattern used for `http.ServeMux`, see https://pkg.go.dev/net/http#ServeMux
		// Matching should be done in Gateway API components, replacing templates values with,
		// application-specific values, proxying to application that will match using path template
		PathPattern string
	}
)

func (h Host) String() string {
	return strings.ToLower(h.String())
}

func (u URL) String() string {
	f := url.URL{
		Host: u.BaseHost.Host,
		Path: u.PathPattern,
	}
	return strings.ToLower(f.String())
}

const (
	// HTTPs protocol scheme
	SchemeHTTPS string = "https"
	// `http.ServeMux` matching pattern every application / service should use
	// `entity` is service / application's name
	// `version` is service / application's SemVer version (e.g. `v0.1.12`)
	APIPathPattern string = "{entity}/{version}/"
)

var (
	// Base host for internet-facing public applications
	BaseHostPublicInternetFacingApp Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotDev.String(),
	})
	// Base host for internet-facing public services
	BaseHostPublicInternetFacingService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotCloud.String(),
	})
	// Base host for internet-facing internal services
	BaseHostInternalInternetFacingService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotRun.String(),
	})
	// Base host for non-internet internal services
	BaseHostInternalPrivateService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotInternal.String(),
	})
	// Base host for preview applications
	BaseHostPrivatePreview Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "preview" + BaseHostInternalPrivateService.Host,
	})

	// Host for company's VCS service
	HostVCS Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "vcs." + BaseHostInternalInternetFacingService.Host,
	})
	// Host for company's chat service
	HostChat Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "chat." + BaseHostInternalInternetFacingService.Host,
	})
	// Host for company's main website
	HostMainWebsite Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "www." + domain.DomainKemaDotDev.String(),
	})
	// Host for company's forum
	HostForum Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "discuss." + domain.DomainKemaDotDev.String(),
	})
	// Host for Kubernetes control planes (kube-apiserver access)
	HostKubeControlePlane Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "kube.{clusterName}." + BaseHostInternalPrivateService.Host + ":6443",
	})
	// Host for Kubernetes control planes (kube-apiserver access)
	HostReviewApp Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		// `repoFQDN` is repository's Fully Qualified Domain Name, i.e. go module name, in lowercase alphanum + dashes form (e.g. `host-tld-owner-repo`)
		Host: "{repoFQDN}-{prNumber}." + BaseHostPrivatePreview.Host,
	})
	// Host used to access consoles of internal services
	HostServiceConsole Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		// `service` is service's name
		Host: "{service}." + BaseHostInternalPrivateService.Host,
	})
	// Host for company's main API
	HostMainApi Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "api." + domain.DomainKemaDotInternal.String(),
	})

	// Base URL for all APIs, providing a common structure for all applications / services
	// All applications / services should use the same pattern
	URLMainApi URL = URL{
		BaseHost:    HostMainApi,
		PathPattern: APIPathPattern,
	}
	// URL for company's security guidelines & responsible disclosure procedure
	URLSecurityGuidelines URL = URL{
		BaseHost:    HostForum,
		PathPattern: "/c/security",
	}
)
