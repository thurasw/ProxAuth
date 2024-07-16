package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/thurasw/ProxAuth/src/internal/api/apps"
	"github.com/thurasw/ProxAuth/src/internal/api/auth"
	"github.com/thurasw/ProxAuth/src/internal/api/users"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/auth", auth.Resource{}.Routes())
	r.Mount("/apps", apps.Resource{}.Routes())
	r.Mount("/users", users.Resource{}.Routes())

	r.NotFound(r.NotFoundHandler())
	return r
}
