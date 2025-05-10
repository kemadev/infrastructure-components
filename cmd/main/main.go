package main

import (
	"log/slog"
	"os"

	"vcs.kema.run/kema/framework-go/pkg/http"
	"vcs.kema.run/kema/framework-go/pkg/log"
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
			slog.String("error.message", err.Error()),
		)
		os.Exit(1)
	}
}
