package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "project/api/project/v1"
	"project/internal/conf"
	"project/internal/service"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer news an HTTP server. For gRPC requests, refer to the function [NewGRPCServer].
//
// This function would read the configuration to configure the HTTP server well,
// and then register the service to the HTTP server.
func NewHTTPServer(c *conf.Server, s *service.ProjectService, m Middlewares) *http.Server {
	// Here we tell the framework that we need these middlewares, and the framework would provide them automatically.
	opts := []http.ServerOption{
		http.Middleware(m...),
	}
	// Reflect the options on the configuration into the memory
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	// Instantiate a new HTTP server listening to a specific port to serve the requests.
	srv := http.NewServer(opts...)
	srv.Handle("/metrics", promhttp.Handler())     // We shall register the Prometheus handler to the server as well
	v1.RegisterProjectManagementHTTPServer(srv, s) // Register the service handlers as well
	return srv
}
