// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-example/app/user/internal/biz"
	"kratos-example/app/user/internal/conf"
	"kratos-example/app/user/internal/data"
	"kratos-example/app/user/internal/server"
	"kratos-example/app/user/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, registry *conf.Registry, jwt *conf.Jwt, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewGormClient(confData, logger)
	client := data.NewRedisClient(confData, logger)
	dataData, cleanup, err := data.NewData(confData, logger, db, client)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, logger)
	userService := service.NewUserService(userUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, userService, logger, jwt)
	registrar, cleanup2, err := data.NewRegistry(registry, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	app := newApp(logger, grpcServer, registrar)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
