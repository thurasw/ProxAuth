package auth

import (
	"net/http"

	"github.com/thurasw/ProxAuth/src/internal/config"
)

// Controller route that will act as the auth server for traefik
func (rs Resource) Traefik(w http.ResponseWriter, r *http.Request) {
	userId := authSessions.GetSession(r)

	if userId == nil {
		AuthHost := config.Config.AuthHost

		// We don't want to return a 401 for traefik proxy, instead redirect to the auth server
		if AuthHost == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else {
			http.Redirect(w, r, AuthHost, http.StatusFound)
		}
		return
	}
	// Session verified
	w.Write([]byte("Authenticated!"))
}
