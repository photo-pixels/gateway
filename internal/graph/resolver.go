package graph

import (
	"context"

	"github.com/photo-pixels/gateway/internal/auth"
	"github.com/photo-pixels/gateway/internal/clients"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type AuthService interface {
	Logout(ctx context.Context) error
	Login(ctx context.Context, email string, password string) error
	GetAccessSession(ctx context.Context) *auth.AccessSession
	GetTokenSession(ctx context.Context) *auth.TokenSession
}

type Resolver struct {
	*clients.ServiceClientsContainer
	AuthService AuthService
}
