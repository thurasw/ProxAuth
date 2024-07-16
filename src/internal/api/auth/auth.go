package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thurasw/ProxAuth/src/internal/api/sessions"
	"github.com/thurasw/ProxAuth/src/internal/config"
	"github.com/thurasw/ProxAuth/src/internal/db"
)

type Resource struct{}

var authSessions *sessions.SessionStore[int]

func (rs Resource) Routes() chi.Router {
	/* Create the session store for user sessions */
	authSessions = sessions.New[int](
		"X-Proxy-Auth",
		time.Hour,
		config.Config.Secure,
		config.Config.Domain,
	)

	r := chi.NewRouter()
	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", rs.Login)
	})
	//Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthOnly)
		r.Get("/logout", rs.Logout)
		r.Put("/password", rs.UpdatePassword)
		r.Get("/traefik", rs.Traefik) // Auth server route for traefik
	})

	r.NotFound(r.NotFoundHandler())
	return r
}

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (rs Resource) Login(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody
	json.NewDecoder(r.Body).Decode(&body)

	user, err := db.GetUserByName(body.Username)
	if err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, "The username or password is incorrect", http.StatusUnauthorized)
		return
	}

	check := checkPasswordHash(body.Password, user.Password)
	if user.Username != body.Username || !check {
		fmt.Printf("Invalid login attempt by user: %s\n", body.Username)
		http.Error(w, "The username or password is incorrect", http.StatusUnauthorized)
		return
	}

	NewUserSession(&user.Id, w, r)
	w.Write([]byte("Success!"))
}

type updatePwBody struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (rs Resource) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(UserCtxKey{}).(int)

	var body updatePwBody
	json.NewDecoder(r.Body).Decode(&body)

	user, err := db.GetUserById(uid)
	if err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	check := checkPasswordHash(body.OldPassword, user.Password)
	if !check {
		fmt.Printf("Invalid password attempt by user: %d\n", uid)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	hashedPw, err := hashPassword(body.NewPassword)
	if err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = db.UpdatePassword(uid, hashedPw)
	if err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Updated!"))
}

func (rs Resource) Logout(w http.ResponseWriter, r *http.Request) {
	authSessions.DeleteSession(w, r)
	w.Write([]byte("Logged out!"))
}
