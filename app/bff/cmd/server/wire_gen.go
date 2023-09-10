// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-example/app/bff/internal/biz"
	"kratos-example/app/bff/internal/conf"
	"kratos-example/app/bff/internal/data"
	"kratos-example/app/bff/internal/server"
	"kratos-example/app/bff/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(bootstrap *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewGormClient(bootstrap, logger)
	client := data.NewRedisClient(bootstrap, logger)
	discovery, cleanup, err := data.NewDiscovery(bootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	userClient := data.NewUserServiceClient(discovery, logger, bootstrap)
	dataData, cleanup2, err := data.NewData(bootstrap, logger, db, client, userClient)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, logger)
	userService := service.NewUserService(userUsecase, logger)
	grpcServer := server.NewGRPCServer(bootstrap, userService, logger)
	httpServer := server.NewHTTPServer(bootstrap, userService, logger)
	registrar, cleanup3, err := data.NewRegistry(bootstrap, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}