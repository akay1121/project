package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "project/api/project/v1"
	"project/internal/conf"
	"project/internal/service"
)

// NewGRPCServer news a gRPC server. For the HTTP server references, check the documentation of [NewHTTPServer].
//
// The server only handles the gRPC calls, which are more commonly used among services, reducing the overall
// overhead cost and communication cost.
func NewGRPCServer(
	c *conf.Server, s *service.ProjectService, m Middlewares) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(m...),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterProjectManagementServer(srv, s)
	return srv
}
