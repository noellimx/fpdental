package endpoints

import (
	"encoding/json"
	"fpdental/auth"
	"fpdental/core"
	"fpdental/server/middlewares"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type EndpointServiceAdmin struct {
	world *core.World
}

func NewEndpointServiceAdmin(world *core.World) *EndpointServiceAdmin {
	var adminE = &EndpointServiceAdmin{}
	adminE.world = world
	return adminE
}

func (ep *EndpointServiceAdmin) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.RestJSON)

	r.Post("/sessions", ep.userSessions)
	r.Post("/revoke-sessions", ep.revokeUserSessions)

	return r
}

func (epA *EndpointServiceAdmin) getUserSessionsAll(token *auth.Token) ([]*auth.UserSessionsBE, error) {

	return epA.world.GetUserSessionsAll(token)
}

func (epA *EndpointServiceAdmin) removeUserSessions(token *auth.Token, userSessions []*auth.UserSessionsBE) {

	epA.world.RemoveUserSessions(token, userSessions)
}

type RequestBodyUserSession = RequestBodyAppointments

func (authA *EndpointServiceAdmin) userSessions(w http.ResponseWriter, r *http.Request) {

	log.Println("[EndpointServiceAdmin::userSessions]")
	var data RequestBodyUserSession

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("[EndpointServiceAdmin::userSessions] %+v", data)

	userSessions, err := authA.getUserSessionsAll(&data.Token)

	if err != nil {
		log.Printf("[EndpointServiceAdmin::userSessions] getUserSessionsAll::Error %s", err)

		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	responseJSON := &struct {
		UserSessions []*auth.UserSessionsBE
	}{UserSessions: userSessions}
	log.Printf("[EndpointServiceAdmin::userSessions] userSessions %+v", responseJSON.UserSessions)

	responseString, err := json.Marshal(responseJSON)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Write([]byte(responseString))
	return

}

type RequestBodyRevokeUserSession struct {
	Token        auth.Token
	UserSessions []*auth.UserSessionsBE
}

func (authA *EndpointServiceAdmin) revokeUserSessions(w http.ResponseWriter, r *http.Request) {

	log.Println("[EndpointServiceAdmin::revokeUserSessions]")
	var data RequestBodyRevokeUserSession

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("[EndpointServiceAdmin::revokeUserSessions] %+v", data)

	authA.removeUserSessions(&data.Token, data.UserSessions)

	if err != nil {
		log.Printf("[EndpointServiceAdmin::revokeUserSessions] -::Error %s", err)

		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Write([]byte("{}"))
	return

}
