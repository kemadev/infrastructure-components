/*
infrastructure-components serves infrastructure-components microservice.
It is deployed in Kubernetes and serves over HTTP.
*/
package main

import (
	"log/slog"
	"os"

	"github.com/kemadev/framework-go/pkg/http"
	"github.com/kemadev/framework-go/pkg/log"
)

func main() {
	// `http.Run()` only returns on init / shutdown failures, where otel logger isn't available
	fallbackLogger := log.CreateFallbackLogger()

	// Define routes to handle
	routes := http.HTTPRoutesToRegister{
		http.HTTPRoute{
			Pattern:     "/rolldice/",
			HandlerFunc: rolldice,
		},
		http.HTTPRoute{
			Pattern:     "/rolldice/{player}",
			HandlerFunc: rolldice,
		},
	}

	// Run HTTP server
	err := http.Run(routes)
	if err != nil {
		fallbackLogger.Error(
			"run",
			slog.String("Body", "http failure"),
			// TODO use semconv value once released, see https://opentelemetry.io/docs/specs/semconv/attributes-registry/error/#error-message
			slog.String("error.message", err.Error()),
		)
		os.Exit(1)
	}
}
