package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"fpdental/auth"
	"fpdental/core"
	"fpdental/server/middlewares"
)

type AuthenticatedEndpoint struct {
	wo *core.World
}

func NewAuthenticatedEndpoint(wo *core.World) *AuthenticatedEndpoint {
	var authE = &AuthenticatedEndpoint{}
	authE.wo = wo
	return authE
}

func (authService *AuthenticatedEndpoint) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.Restful)
	// r.Get("/*", func(w http.ResponseWriter, r *http.Request) {

	// 	log.Println("weolcome auth hi")
	// 	w.Write([]byte("welcome auth"))
	// })
	r.Post("/login", authService.login)

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

func (authE *AuthenticatedEndpoint) login(w http.ResponseWriter, r *http.Request) {

	log.Println("authE.login")
	var data RequestBodyAuth

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("%+v", data)

	is, token := authE.wo.IsAuthenticated(data.Username, data.Password)

	log.Println("authE.login1")

	responseJSON := &struct {
		Token *auth.Token
	}{
		Token: nil,
	}
	if is {
		responseJSON.Token = token
	} else {
		status := http.StatusUnauthorized
		http.Error(w, http.StatusText(status), status)
	}

	responseString, err := json.Marshal(responseJSON)
	if err != nil {

		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Write([]byte(responseString))
	return

}
