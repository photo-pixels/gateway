package app

import (
	"context"
	"fmt"

	"github.com/photo-pixels/gateway/internal/auth"
	"github.com/photo-pixels/gateway/internal/auth/jwt_helper"
	"github.com/photo-pixels/gateway/internal/clients"
	"github.com/photo-pixels/gateway/internal/server"
	"github.com/photo-pixels/platform/config"
	"github.com/photo-pixels/platform/log"
)

// App приложение
type App struct {
	cfgProvider config.Provider
	logger      log.Logger
	// cfg
	serverCfg server.Config
	//
	cc *clients.ServiceClientsContainer
	//
	jwtHelper *jwt_helper.JwtHelper
	auth      *auth.Auth
}

// NewApp новое приложение
func NewApp(cfgProvider config.Provider) *App {
	return &App{cfgProvider: cfgProvider}
}

// Create создание сервисов
func (a *App) Create(ctx context.Context) error {
	a.logger = log.NewLogger()

	var err error
	a.serverCfg, err = a.getServerConfig()
	if err != nil {
		return fmt.Errorf("getServerConfig: %w", err)
	}

	serviceClientsContainerCfg, err := a.getServiceClientsContainerConfig()
	if err != nil {
		return fmt.Errorf("getServiceClientsContainerConfig: %w", err)
	}
	a.cc, err = clients.NewServiceClientsContainer(serviceClientsContainerCfg)
	if err != nil {
		return fmt.Errorf("clients.NewServiceClientsContainer: %w", err)
	}

	jwtHelperCfg, err := a.getJwtHelperConfig()
	if err != nil {
		return fmt.Errorf("getJwtHelperConfig: %w", err)
	}
	a.jwtHelper, err = jwt_helper.NewHelper(jwtHelperCfg)
	if err != nil {
		return fmt.Errorf("jwt_helper.NewHelper: %w", err)
	}

	a.auth = auth.NewAuth(
		a.logger,
		a.jwtHelper,
		a.cc.GetAuthClient(),
		a.cc.GetTokenClient(),
	)

	return err
}

// GetLogger получить логер
func (a *App) GetLogger() log.Logger {
	return a.logger
}

// GetServerConfig конфигуратор сервера
func (a *App) GetServerConfig() server.Config {
	return a.serverCfg
}

// GetServiceClientsContainer список grpc клиентов
func (a *App) GetServiceClientsContainer() *clients.ServiceClientsContainer {
	return a.cc
}

func (a *App) GetAuth() *auth.Auth {
	return a.auth
}
