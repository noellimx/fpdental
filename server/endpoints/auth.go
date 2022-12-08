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

type EndpointServiceAuthentication struct {
	wo *core.World
}

func NewEndpointServiceAuthentication(wo *core.World) *EndpointServiceAuthentication {
	var authE = &EndpointServiceAuthentication{}
	authE.wo = wo
	return authE
}

func (authService *EndpointServiceAuthentication) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.RestJSON)

	r.Post("/login", authService.login)
	r.Post("/is-valid-token", authService.isValidToken)

	return r
}

type RequestBodyAuthLogin struct {
	Username string
	Password string
}

func (authE *EndpointServiceAuthentication) login(w http.ResponseWriter, r *http.Request) {

	log.Println("authE.login")
	var data RequestBodyAuthLogin

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

type RequestBodyAuthIsValidToken struct {
	Token auth.Token
}

func (authE *EndpointServiceAuthentication) isValidToken(w http.ResponseWriter, r *http.Request) {

	log.Println("[authE.isValidToken]")
	var data RequestBodyAuthIsValidToken

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("[authE.isValidToken] data <- %+v", data)

	is := authE.wo.IsValidToken(&data.Token)
	log.Printf("[authE.isValidToken] %t", is)

	type Response struct {
		Is bool
	}
	responseJSON := &Response{}
	responseJSON.Is = is

	if is {
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

	log.Printf("[authE.isValidToken] responseString %s", responseString)

	w.Write([]byte(responseString))
	return

}
