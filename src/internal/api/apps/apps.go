package apps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thurasw/ProxAuth/src/internal/api/auth"
	"github.com/thurasw/ProxAuth/src/internal/db"
)

type Resource struct{}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(auth.AuthOnly)
	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Put("/{appId}", rs.Edit)
	r.Delete("/{appId}", rs.Delete)

	r.NotFound(r.NotFoundHandler())
	return r
}

func (rs Resource) List(w http.ResponseWriter, r *http.Request) {
	res, err := db.GetApps()
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	jData, err := json.Marshal(res)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

type createAppBody struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Logo  string `json:"logo"`
}

func (rs Resource) Create(w http.ResponseWriter, r *http.Request) {
	var body createAppBody
	json.NewDecoder(r.Body).Decode(&body)

	err := db.CreateApp(db.AppItem{
		Name:  body.Name,
		Color: body.Color,
		Logo:  []byte(body.Logo),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Created!"))
}

func (rs Resource) Edit(w http.ResponseWriter, r *http.Request) {
	appId := chi.URLParam(r, "appId")

	id, err := strconv.Atoi(appId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var body createAppBody
	json.NewDecoder(r.Body).Decode(&body)

	err = db.UpdateApp(db.AppItem{
		Id:    id,
		Name:  body.Name,
		Color: body.Color,
		Logo:  []byte(body.Logo),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Updated!"))
}

func (rs Resource) Delete(w http.ResponseWriter, r *http.Request) {
	appId := chi.URLParam(r, "appId")

	id, err := strconv.Atoi(appId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = db.DeleteApp(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Deleted!"))
}
