package host

import (
	"net/url"
	"strings"

	"github.com/kemadev/infrastructure-components/pkg/private/domain"
)

type (
	// A Host is a representation of a DNS host name.
	Host url.URL
	// A URL is a representation of an URL,
	URL struct {
		// Host where traffic is destinated
		BaseHost Host
		// Path pattern used for `http.ServeMux`, see https://pkg.go.dev/net/http#ServeMux
		// Matching should be done in Gateway API components, replacing templates values with,
		// application-specific values, proxying to application that will match using path template
		PathPattern string
	}
)

// String returns the string representation of the Host, lowercased.
func (h Host) String() string {
	return strings.ToLower(h.String())
}

// String returns the string representation of the URL, lowercased.
func (u URL) String() string {
	f := url.URL{
		Host: u.BaseHost.Host,
		Path: u.PathPattern,
	}
	return strings.ToLower(f.String())
}

const (
	// SchemeHTTPS is the HTTPS scheme
	SchemeHTTPS string = "https"
)

var (
	// BaseHostPublicInternetFacingApp is the base host for internet-facing public applications. This is where client-facing
	// applications are hosted.
	BaseHostPublicInternetFacingApp Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotDev.String(),
	})
	// BaseHostPublicInternetFacingService is the base host for internet-facing public services. This is where client-facing
	// services are hosted, service meaning a service for the client such as a forum, a chat, ... , not a microservice.
	BaseHostPublicInternetFacingService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotCloud.String(),
	})
	// BaseHostInternalInternetFacingService is the base host for internet-facing internal services. This is where internal
	// services are hosted, service meaning a service for internal use such as VCS, chat, ... , not a microservice.
	BaseHostInternalInternetFacingService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotRun.String(),
	})
	// BaseHostInternalPrivateService is the base host for internal private services. This is where internal services are hosted, such as
	// Kubernetes control planes, preview applications, ...
	BaseHostInternalPrivateService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotInternal.String(),
	})
	// BaseHostPrivatePreview is the base host for private preview applications.
	BaseHostPrivatePreview Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "preview" + BaseHostInternalPrivateService.Host,
	})

	// HostVCS is the host for company's VCS.
	HostVCS Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "vcs." + BaseHostInternalInternetFacingService.Host,
	})
	// HostChat is the host for company's chat.
	HostChat Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "chat." + BaseHostInternalInternetFacingService.Host,
	})
	// HostMainWebsite is the host for company's main website.
	HostMainWebsite Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "www." + domain.DomainKemaDotDev.String(),
	})
	// HostForum is the host for company's forum.
	HostForum Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "discuss." + domain.DomainKemaDotDev.String(),
	})
	// HostKubeControlePlane is the host for Kubernetes control planes (kube-apiserver access)
	HostKubeControlePlane Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "kube.{clusterName}." + BaseHostInternalPrivateService.Host + ":6443",
	})
	// HostReviewApp is the host for preview applications.
	HostReviewApp Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		// `repoFQDN` is repository's Fully Qualified Domain Name, i.e. go module name, in lowercase alphanum + dashes form (e.g. `host-tld-owner-repo`)
		Host: "{repoFQDN}-{prNumber}." + BaseHostPrivatePreview.Host,
	})
	// HostServiceConsole is the host for service consoles.
	HostServiceConsole Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		// `service` is service's name
		Host: "{service}." + BaseHostInternalPrivateService.Host,
	})
	// HostMainApi is the host for company's main API.
	HostMainApi Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "api." + domain.DomainKemaDotInternal.String(),
	})

	// ServiceNamePathPattern is the path pattern for service name, to be replaced by the service name in path matching.
	ServiceNamePathPattern    = "{service}"
	// ServiceVersionPathPattern is the path pattern for service version, to be replaced by the service version (major) in path matching.
	ServiceVersionPathPattern = "{version}"
	// URLMainApi is the URL for conventional [net/http.ServeMux] matching pattern every application / service should use,
	// providing a common structure for all applications / services.
	// All applications / services should use this pattern.
	URLMainApi URL = URL{
		BaseHost:    HostMainApi,
		PathPattern: "/" + ServiceNamePathPattern + "/" + ServiceVersionPathPattern + "/",
	}
	// URLSecurityGuidelines is the URL for company's security guidelines & responsible disclosure procedure
	URLSecurityGuidelines URL = URL{
		BaseHost:    HostForum,
		PathPattern: "/c/security",
	}
)
