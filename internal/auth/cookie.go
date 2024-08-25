package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ContextKey string

const (
	accessSessionKey ContextKey = "accessSession"
	sessionKey       ContextKey = "session"
	httpResponseKey  ContextKey = "httpResponse"
	tokenKey         ContextKey = "token"
	tokenSessionKey  ContextKey = "tokenSession"
)

const (
	userIDCookieName       = "user_id"
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
)

func getSessionFromCookies(r *http.Request) (*session, error) {
	userID, err := r.Cookie(userIDCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user id: %w", err)
	}

	accessToken, err := r.Cookie(accessTokenCookieName)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	refreshToken, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	return &session{
		UserID: userID.Value,
		AccessToken: func() *string {
			if accessToken == nil {
				return nil
			}
			return &accessToken.Value
		}(),
		RefreshToken: refreshToken.Value,
	}, nil
}

func getResponse(ctx context.Context) http.ResponseWriter {
	v, ok := ctx.Value(httpResponseKey).(http.ResponseWriter)
	if !ok {
		return nil
	}
	return v
}

func setCookie(ctx context.Context, name string, value string, expires time.Time) {
	w := getResponse(ctx)
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		// Secure:   true,
	})
}

func deleteCookie(ctx context.Context, name string) {
	w := getResponse(ctx)
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
}
