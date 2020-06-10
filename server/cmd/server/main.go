package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ritchie-server/server"
	"ritchie-server/server/starter"
)

var h Handler

type Handler struct {
	LoginHandler            server.DefaultHandler
	CredentialConfigHandler server.DefaultHandler
	ConfigHealth            server.DefaultHandler
	UsageLoggerHandler      server.DefaultHandler
	CliVersionHandler       server.DefaultHandler
	RepositoryHandler       server.DefaultHandler
	TreeHandler             server.DefaultHandler
	FormulasHandler         server.DefaultHandler
	MiddlewareHandler       server.MiddlewareHandler
	CredentialHandler       server.CredentialHandler
	HelloHandler            server.DefaultHandler
}

func init() {
	i, err := starter.NewConfiguration()
	if err != nil {
		log.Fatalf("Failed to load server configuration: %v", err)
	}
	h = Handler{
		LoginHandler:            i.LoadLoginHandler(),
		CredentialConfigHandler: i.LoadCredentialConfigHandler(),
		ConfigHealth:            i.LoadConfigHealth(),
		UsageLoggerHandler:      i.LoadUsageLoggerHandler(),
		CliVersionHandler:       i.LoadCliVersionHandler(),
		RepositoryHandler:       i.LoadRepositoryHandler(),
		TreeHandler:             i.LoadTreeHandler(),
		FormulasHandler:         i.LoadFormulasHandler(),
		MiddlewareHandler:       i.LoadMiddlewareHandler(),
		CredentialHandler:       i.LoadCredentialHandler(),
		HelloHandler:            i.LoadHelloHandler(),
	}
}

func main() {

	log.Info("Starting server")
	http.Handle("/login", h.MiddlewareHandler.Filter(h.LoginHandler.Handler()))
	http.Handle("/credentials/admin", h.MiddlewareHandler.Filter(h.CredentialHandler.HandleAdmin()))
	http.Handle("/credentials/org", h.MiddlewareHandler.Filter(h.CredentialHandler.HandleOrg()))
	http.Handle("/credentials/me", h.MiddlewareHandler.Filter(h.CredentialHandler.HandleMe()))
	http.Handle("/credentials/me/", h.MiddlewareHandler.Filter(h.CredentialHandler.HandleMe()))
	http.Handle("/credentials/config", h.MiddlewareHandler.Filter(h.CredentialConfigHandler.Handler()))
	http.Handle("/metrics", h.MiddlewareHandler.Filter(promhttp.Handler()))
	http.Handle("/usage", h.MiddlewareHandler.Filter(h.UsageLoggerHandler.Handler()))
	http.Handle("/health", h.MiddlewareHandler.Filter(h.ConfigHealth.Handler()))
	http.Handle("/cli-version", h.MiddlewareHandler.Filter(h.CliVersionHandler.Handler()))
	http.Handle("/repositories", h.MiddlewareHandler.Filter(h.RepositoryHandler.Handler()))
	http.Handle("/tree/", h.MiddlewareHandler.Filter(h.TreeHandler.Handler()))
	http.Handle("/formulas/", h.MiddlewareHandler.Filter(h.FormulasHandler.Handler()))
	http.Handle("/", h.MiddlewareHandler.Filter(h.HelloHandler.Handler()))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
