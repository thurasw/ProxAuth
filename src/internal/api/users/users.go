package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thurasw/ProxAuth/src/internal/api/auth"
	"github.com/thurasw/ProxAuth/src/internal/db"
)

type Resource struct{}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(auth.AuthOnly)
	r.Get("/", rs.Get)
	r.Put("/", rs.Edit)

	r.NotFound(r.NotFoundHandler())
	return r
}

func (rs Resource) Get(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(auth.UserCtxKey{}).(int)

	res, err := db.GetUserById(int(uid))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	jData, err := json.Marshal(*res)
	if err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

type editUserBody struct {
	Username string `json:"username"`
}

func (rs Resource) Edit(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(auth.UserCtxKey{}).(int)

	var body editUserBody
	json.NewDecoder(r.Body).Decode(&body)

	err := db.UpdateUser(uid, body.Username)
	if err != nil {
		fmt.Printf("%v\n", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Write([]byte("Updated!"))
}
