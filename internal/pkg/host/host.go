package host

import (
	"net/url"

	"vcs.kema.run/kema/infrastructure-components/internal/pkg/domain"
)

type (
	// A base host, used to construct URLs
	Host url.URL
	// An URL template used to distribute traffic
	URL struct {
		// Host where traffic is destinated
		BaseHost Host
		// Path pattern used for `http.ServeMux`, see https://pkg.go.dev/net/http#ServeMux
		// Please note that `HOST` part won't be used, host matching is done via `BaseHost` in Gateway API components
		PathPattern string
	}
)

func (h Host) String() string {
	return h.String()
}

func (u URL) String() string {
	f := url.URL{
		Host: u.BaseHost.Host,
		Path: u.PathPattern,
	}
	return f.String()
}

const (
	// HTTPs protocol scheme
	SchemeHTTPS string = "https"
)

var (
	// Base host for internet-facing public applications
	HostPublicInternetFacingApp Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotDev.String(),
	})
	// Base host for internet-facing public services
	HostPublicInternetFacingService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotCloud.String(),
	})
	// Base host for internet-facing internal services
	HostInternalInternetFacingService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotRun.String(),
	})
	// Base host for non-internet internal services
	HostInternalPrivateService Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   domain.DomainKemaDotInternal.String(),
	})

	// Host for company's main website
	HostMainWebsite Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "www." + domain.DomainKemaDotDev.String(),
	})
	// Host for company's main API
	HostMainApi Host = Host(url.URL{
		Scheme: SchemeHTTPS,
		Host:   "api." + domain.DomainKemaDotInternal.String(),
	})
	URLMainApi URL = URL{
		BaseHost:    HostMainApi,
		PathPattern: "{service}/v{version}/",
	}
)

// ## Domains definition

// - Internet-facing public production - `kema.dev`
//   - Main website - `www.kema.dev`
//   - API host - `api.kema.dev`
//     - API routes - `api.kema.dev/<service>/v<version>/*`
//       - Microservice name - `service`
//       - Microservice API version - `version`
// - Internet-facing public services - `kema.cloud`
//   - Main website - `www.kema.cloud`
//   - Forum - `discuss.kema.cloud`
//     - Security guidelines & responsible disclosure procedure - `https://discuss.kema.cloud/c/security`
// - `kema.run` - Internet-facing internal services
//   - VCS - `vcs.kema.run`
//   - Internal chat & VOIP - `chat.kema.run`
// - Private internal - `kema.internal`
//   - Kubernetes control planes nodes VIP - `kube.<cluster name>.kema.internal`
//   - Review applications - `preview.kema.internal`
//     - Review applications - `<repo>-<pr number>.preview.kema.internal`
//       - Repository FQDN holding application code - `<repo>`
//       - PR number that triggered deployment - `<pr number>`
//   - Internal services / administration consoles - `<service>.kema.internal`
//     - Service name - `service`
// - Internet-facing personal - `doublej.fr`
