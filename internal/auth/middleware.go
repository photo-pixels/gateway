package auth

import (
	"context"
	"fmt"
	"net/http"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), httpResponseKey, w)
		getSession, err := getSessionFromCookies(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("error get session: %v", err), http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, sessionKey, getSession)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
