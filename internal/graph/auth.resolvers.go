package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/photo-pixels/gateway/internal/graph/model"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.Result, error) {
	err := r.AuthService.Login(ctx, input.Email, input.Password)
	if err != nil {
		return nil, handleError("AuthService", err)
	}
	return &model.Result{Success: true}, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (*model.Result, error) {
	session := r.AuthService.GetAccessSession(ctx)
	if session == nil {
		return nil, model.ErrorNoAuth
	}

	err := r.AuthService.Logout(ctx)
	if err != nil {
		return nil, handleError("AuthService", err)
	}
	return &model.Result{Success: true}, err
}
