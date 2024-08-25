package auth

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/photo-pixels/gateway/internal/graph/model"
	userAccount "github.com/photo-pixels/gateway/pkg/gen/user_account"
)

// IsToken достает данные сессии из куки, и пытается из AccessToken получить AccessSession
func (s *Auth) IsToken(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	apiToken, ok := ctx.Value(tokenKey).(string)
	if !ok || apiToken == "" {
		return nil, model.ErrorTokenIsNotValid
	}

	tokenSession, err := s.getToken(ctx, apiToken)
	if err != nil {
		return nil, fmt.Errorf("getToken: %w", err)
	}

	ctx = context.WithValue(ctx, tokenSessionKey, &tokenSession)
	return next(ctx)
}

// GetTokenSession получить TokenSession из контекста
func (s *Auth) GetTokenSession(ctx context.Context) *TokenSession {
	v, ok := ctx.Value(tokenSessionKey).(*TokenSession)
	if !ok {
		return nil
	}
	return v
}

func (s *Auth) getToken(ctx context.Context, token string) (TokenSession, error) {
	response, err := s.tokenServiceClient.GetToken(ctx, &userAccount.GetTokenRequest{
		Token: "",
	})
	if err != nil {
		return TokenSession{}, fmt.Errorf("s.tokenServiceClient.GetToken: %w", err)
	}

	userID, err := uuid.Parse(response.Token.UserId)
	if err != nil {
		return TokenSession{}, fmt.Errorf("userID is invalid: %w", err)
	}

	return TokenSession{
		UserID:    userID,
		TokenType: response.Token.TokenType,
	}, nil
}
