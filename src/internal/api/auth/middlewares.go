package auth

import (
	"context"
	"net/http"
)

// Key for username context
type UserCtxKey struct{}

// Auth guard - verifies the session token from cookie and adds user id to the context
func AuthOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := authSessions.GetSession(r)

		if userId == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Add user id from claim to context
		newCtx := context.WithValue(r.Context(), UserCtxKey{}, *userId)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
