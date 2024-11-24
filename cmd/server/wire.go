//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"project/internal/biz"
	"project/internal/conf"
	"project/internal/data"
	"project/internal/server"
	"project/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application by injecting the required dependencies via generated code.
//
// The following code is not the final production code, it just declares the dependency providers and the
// injection code is generated in the file `wire_gen.go`, which implements the wiring process.
func wireApp(*conf.Registry, *conf.Server, *conf.Data, *conf.Telemetry) (*kratos.App, func(), error) {
	panic(
		wire.Build( // Finally replaced by the real initialization code, the wire.Build call here is just a placeholder
			server.ProviderSet,  // Server that responses to the client requests
			data.ProviderSet,    // Manage the interaction with database, serving as the DAO layer
			biz.ProviderSet,     // Business logics providing the functionalities of performing operations
			service.ProviderSet, // Expose the interface of manipulating the business
			newApp,              // Set up the common foundation (log, register, etc.) and then instantiate the application
		),
	)
}
