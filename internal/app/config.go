package app

import (
	"fmt"
	"github.com/photo-pixels/gateway/internal/auth/jwt_helper"
	"github.com/photo-pixels/gateway/internal/clients"

	"github.com/photo-pixels/gateway/internal/server"
)

const (
	// ServerConfigName конфиг сервера
	ServerConfigName = "server"
	// ServiceClientsName конфиг клиентов
	ServiceClientsName = "clients"
	// JwtHelperName данные для jwt
	JwtHelperName = "jwt_helper"
)

func (a *App) getServerConfig() (server.Config, error) {
	var config server.Config
	err := a.cfgProvider.PopulateByKey(ServerConfigName, &config)
	if err != nil {
		return server.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getServiceClientsContainerConfig() (clients.Config, error) {
	var config clients.Config
	err := a.cfgProvider.PopulateByKey(ServiceClientsName, &config)
	if err != nil {
		return clients.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}

	return config, nil
}

func (a *App) getJwtHelperConfig() (jwt_helper.Config, error) {
	var config jwt_helper.Config
	err := a.cfgProvider.PopulateByKey(JwtHelperName, &config)
	if err != nil {
		return jwt_helper.Config{}, fmt.Errorf("PopulateByKey: %w", err)
	}
	return config, nil
}
