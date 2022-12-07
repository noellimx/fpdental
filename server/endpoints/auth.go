package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"fpdental/server/middlewares"
)

type Auth struct {
}

func (auth *Auth) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.Restful)

	r.Post("/", auth.login) // POST /users - create a new user and persist it

	// r.Route("/{id}", func(r chi.Router) {
	// 	// r.Use(TodoCtx) // lets have a users map, and lets actually load/manipulate
	// 	r.Get("/", Get)       // GET /users/{id} - read a single user by :id
	// 	r.Put("/", Update)    // PUT /users/{id} - update a single user by :id
	// 	r.Delete("/", Delete) // DELETE /users/{id} - delete a single user by :id
	// })

	return r
}

type RequestBodyAuth struct {
	Username string
	Password string
}

func (a *Auth) login(w http.ResponseWriter, r *http.Request) {
	var data RequestBodyAuth

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)

	responseJSON := &struct {
		Username      string
		StatusMessage string
	}{
		Username:      data.Username,
		StatusMessage: http.StatusText(http.StatusNotImplemented),
	}
	responseString, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Write([]byte(responseString))
}

func NewAuth() *Auth {
	return &Auth{}
}
