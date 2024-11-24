// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"project/internal/biz"
	"project/internal/conf"
	"project/internal/data"
	"project/internal/server"
	"project/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application by injecting the required dependencies via generated code.
//
// The following code is not the final production code, it just declares the dependency providers and the
// injection code is generated in the file `wire_gen.go`, which implements the wiring process.
func wireApp(registry *conf.Registry, confServer *conf.Server, confData *conf.Data, telemetry *conf.Telemetry) (*kratos.App, func(), error) {
	registrar := server.NewRegistry(registry)
	dataData, cleanup, err := data.NewData(confData)
	if err != nil {
		return nil, nil, err
	}
	cache, cleanup2, err := data.NewCache(confData)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	projectRepository := data.NewProjectRepository(dataData, cache)
	projectManager := biz.NewProjectManager(projectRepository)
	projectService := service.NewProjectService(projectManager)
	middlewares := server.NewMiddlewares(telemetry)
	grpcServer := server.NewGRPCServer(confServer, projectService, middlewares)
	httpServer := server.NewHTTPServer(confServer, projectService, middlewares)
	app := newApp(registrar, grpcServer, httpServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
