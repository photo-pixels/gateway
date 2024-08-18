package auth

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/photo-pixels/gateway/internal/auth/jwt_helper"
	"github.com/photo-pixels/gateway/internal/graph/model"
	userAccount "github.com/photo-pixels/gateway/pkg/gen/user_account"
	"github.com/photo-pixels/platform/log"
)

// JWTHelper хелпер для работы с jwt
type JWTHelper interface {
	Parse(token string, claims jwt_helper.Claims) error
}

type Auth struct {
	logger            log.Logger
	jwtHelper         JWTHelper
	authServiceClient userAccount.AuthServiceClient
}

// NewAuth новый менеджер работы с токенам
func NewAuth(
	logger log.Logger,
	jwtHelper JWTHelper,
	authServiceClient userAccount.AuthServiceClient,
) *Auth {
	return &Auth{
		logger:            logger.Named("auth"),
		jwtHelper:         jwtHelper,
		authServiceClient: authServiceClient,
	}
}

// IsAuthenticated достает данные сессии из куки, и пытается из AccessToken получить AccessSession
func (s *Auth) IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	getSession, ok := ctx.Value(sessionKey).(*session)
	if !ok || getSession == nil {
		// Пользователь не залоген
		return nil, model.ErrorSessionNoFound
	}

	if getSession.AccessToken == nil {
		accessToken, err := s.refresh(ctx, getSession.RefreshToken)
		ctx = context.WithValue(ctx, sessionKey, getSession)
		if err != nil {
			return nil, model.ErrorNoAuth
		}
		getSession.AccessToken = &accessToken
		return s.IsAuthenticated(ctx, obj, next)
	}

	claims := new(accessSessionClaims)
	err = s.jwtHelper.Parse(*getSession.AccessToken, claims)
	if err != nil {
		return nil, model.ErrorNoAuth
	}

	ctx = context.WithValue(ctx, accessSessionKey, &claims.AccessSession)
	return next(ctx)
}

// GetAccessSession получить AccessSession из контекста
func (s *Auth) GetAccessSession(ctx context.Context) *AccessSession {
	v, ok := ctx.Value(accessSessionKey).(*AccessSession)
	if !ok {
		return nil
	}
	return v
}

func (s *Auth) setSession(ctx context.Context, session *userAccount.AuthData) {
	setCookie(ctx, userIDCookieName, session.UserId, session.RefreshTokenExpiration.AsTime())
	setCookie(ctx, accessTokenCookieName, session.AccessToken, session.AccessTokenExpiration.AsTime())
	setCookie(ctx, refreshTokenCookieName, session.RefreshToken, session.RefreshTokenExpiration.AsTime())
}

func (s *Auth) deleteSession(ctx context.Context) {
	deleteCookie(ctx, userIDCookieName)
	deleteCookie(ctx, accessTokenCookieName)
	deleteCookie(ctx, refreshTokenCookieName)
}

func (s *Auth) refresh(ctx context.Context, refreshToken string) (string, error) {
	authData, err := s.authServiceClient.RefreshToken(ctx, &userAccount.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return "", fmt.Errorf("s.authServiceClient.RefreshToken: %w", err)
	}
	s.setSession(ctx, authData)
	return authData.AccessToken, nil
}

func (s *Auth) Login(ctx context.Context, email string, password string) error {
	authData, err := s.authServiceClient.Login(ctx, &userAccount.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return fmt.Errorf("s.authServiceClient.Login: %w", err)
	}

	s.setSession(ctx, authData)
	return nil
}

func (s *Auth) Logout(ctx context.Context) error {
	getSession, ok := ctx.Value(sessionKey).(*session)
	if !ok || getSession == nil {
		// Пользователь не залоген
		return model.ErrorSessionNoFound
	}

	_, err := s.authServiceClient.Logout(ctx, &userAccount.LogoutRequest{
		RefreshToken: getSession.RefreshToken,
	})
	if err != nil {
		return fmt.Errorf("s.authServiceClient.Logout: %w", err)
	}
	s.deleteSession(ctx)
	return nil
}
