package endpoints

import (
	"encoding/json"
	"fpdental/appointment"
	"fpdental/auth"
	"fpdental/core"
	"fpdental/server/middlewares"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type EndpointServiceUser struct {
	world *core.World
}

func NewEndpointServiceUser(wo *core.World) *EndpointServiceUser {
	var authE = &EndpointServiceUser{}
	authE.world = wo
	return authE
}

func (ep *EndpointServiceUser) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.RestJSON)

	r.Post("/", ep.appointments)
	return r
}

func (ep *EndpointServiceUser) getAppointments(token *auth.Token) (*appointment.Appointments, error) {
	appointments, err := ep.world.GetUserAppointments(token)

	return appointments, err
}

type RequestBodyAppointments struct {
	Token auth.Token
}

func (authE *EndpointServiceUser) appointments(w http.ResponseWriter, r *http.Request) {

	log.Println("[EndpointServiceUser::appointments]")
	var data RequestBodyAppointments

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("[EndpointServiceUser::appointments] %+v", data)

	appointments, err := authE.getAppointments(&data.Token)

	if err != nil {

		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	responseJSON := &struct {
		Appointments []*appointment.Appointment
	}{}
	responseJSON.Appointments = appointments.AsSlice()
	log.Printf("[EndpointServiceUser::appointments] appointments %+v", responseJSON.Appointments)

	responseString, err := json.Marshal(responseJSON)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Write([]byte(responseString))
	return

}
