package host

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/kemadev/infrastructure-components/pkg/private/domain"
)

type (
	// A URL is a representation of an URL,
	URL struct {
		// Host where traffic is destinated
		BaseHost url.URL
		// Path pattern used for `http.ServeMux`, see https://pkg.go.dev/net/http#ServeMux
		// Matching should be done in Gateway API components, replacing templates values with,
		// application-specific values, proxying to application that will match using path template
		PathPattern string
	}
)

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
	BaseHostPublicInternetFacingApp url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotDev.String(),
	}
	// BaseHostPublicInternetFacingService is the base host for internet-facing public services. This is where client-facing
	// services are hosted, service meaning a service for the client such as a forum, a chat, ... , not a microservice.
	BaseHostPublicInternetFacingService url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotCloud.String(),
	}
	// BaseHostInternalInternetFacingService is the base host for internet-facing internal services. This is where internal
	// services are hosted, service meaning a service for internal use such as preview app, chat, ... , not a microservice.
	BaseHostInternalInternetFacingService url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotRun.String(),
	}
	// BaseHostInternalPrivateService is the base host for internal private services. This is where internal services are hosted, such as
	// Kubernetes control planes, ...
	BaseHostInternalPrivateService url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotInternal.String(),
	}
	// BaseHostPrivatePreview is the base host for private preview applications.
	BaseHostPrivatePreview url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   "preview." + BaseHostInternalInternetFacingService.Host,
	}

	// HostVCS is the host for company's VCS.
	HostVCS url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   "github.com",
	}
	// HostChat is the host for company's chat.
	HostChat url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   "github.com",
	}

	// HostMainWebsite is the host for company's main website.
	HostMainWebsite url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   "www." + domain.DomainKemaDotDev.String(),
	}

	// HostForum is the host for company's forum.
	HostForum url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   "github.com",
	}

	// HostReviewApp is the host for preview applications.
	HostReviewApp = func(repo url.URL, prNumber int) url.URL {
		repoNoDot := strings.ReplaceAll(
			strings.ToLower(repo.Hostname()),
			".",
			"-",
		)
		repoClean := strings.ReplaceAll(
			repoNoDot,
			"/",
			"-",
		)
		repoFQDN := strings.TrimSuffix(
			repoClean,
			"-",
		)
		return url.URL{
			Scheme: SchemeHTTPS,
			Host: repoFQDN + "-" + strconv.Itoa(
				prNumber,
			) + "." + BaseHostPrivatePreview.Host,
		}
	}

	// HostServiceConsole is the host for service consoles.
	HostServiceConsole = func(serviceName string) url.URL {
		return url.URL{
			Scheme: SchemeHTTPS,
			Host:   serviceName + "." + BaseHostInternalPrivateService.Host,
		}
	}
	// HostMainApi is the host for company's main API.
	HostMainApi url.URL = url.URL{
		Scheme: SchemeHTTPS,
		Host:   "api." + domain.DomainKemaDotInternal.String(),
	}

	// ServiceNamePathPattern is the path pattern for service name, to be replaced by the service name in path matching.
	ServiceNamePathPattern = "{service}"
	// ServiceVersionPathPattern is the path pattern for service version, to be replaced by the service version (major) in path matching.
	ServiceVersionPathPattern = "{version}"
	// URLMainApi is the URL for conventional [net/http.ServeMux] matching pattern every application / service should use,
	// providing a common structure for all applications / services.
	// All applications / services should use this pattern.
	URLMainApi = func(serviceName, serviceVersion string) URL {
		return URL{
			BaseHost:    HostMainApi,
			PathPattern: "/" + serviceName + "/" + serviceVersion + "/",
		}
	}
	// URLSecurityGuidelines is the URL for company's security guidelines & responsible disclosure procedure
	URLSecurityGuidelines URL = URL{
		BaseHost:    HostForum,
		PathPattern: "/c/security",
	}
)
