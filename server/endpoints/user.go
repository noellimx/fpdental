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
	"github.com/google/uuid"
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

	r.Post("/", ep.userAppointments)

	r.Post("/release", ep.releaseAppointment)
	r.Get("/avail", ep.availableAppointments)
	r.Post("/book", ep.book)

	return r
}

func (authE *EndpointServiceUser) book(w http.ResponseWriter, r *http.Request) {

	log.Println("[EndpointServiceUser::book]")
	var data RequestBodyBookAppointment

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("[EndpointServiceUser::book] Data %+v", data)

	if err != nil {

		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	err = authE.world.BookAppointment(&data.Token, data.AppointmentId)

	responseJSON := struct {
		Is    bool
		Error error
	}{}
	if err != nil {
		status := http.StatusBadRequest
		responseJSON.Is = false
		responseJSON.Error = err
		http.Error(w, http.StatusText(status), status)
	} else {
		responseJSON.Is = true
	}

	responseString, err := json.Marshal(responseJSON)

	w.Write([]byte(responseString))

	return

}

func (ep *EndpointServiceUser) getAppointments(token *auth.Token) (*appointment.Appointments, error) {
	appointments, err := ep.world.GetUserAppointments(token)

	return appointments, err
}

type RequestBodyAppointments struct {
	Token auth.Token
}

func (authE *EndpointServiceUser) userAppointments(w http.ResponseWriter, r *http.Request) {

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
		log.Printf("[EndpointServiceUser::appointments] getAppointments::Error %s", err)

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
func (authE *EndpointServiceUser) availableAppointments(w http.ResponseWriter, r *http.Request) {

	log.Println("[EndpointServiceUser::availableAppointments]")

	appointments, err := authE.world.GetAvailableAppointments()

	if err != nil {
		log.Printf("[EndpointServiceUser::availableAppointments] getAppointments::Error %s", err)

		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	responseJSON := &struct {
		Appointments []*appointment.Appointment
	}{}
	responseJSON.Appointments = appointments.AsSlice()
	log.Printf("[EndpointServiceUser::availableAppointments] appointments %+v", responseJSON.Appointments)

	responseString, err := json.Marshal(responseJSON)
	if err != nil {
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}
	w.Write([]byte(responseString))
	return

}

type RequestBodyReleaseAppointment struct {
	Token         auth.Token
	AppointmentId uuid.UUID
}

type RequestBodyBookAppointment = RequestBodyReleaseAppointment

func (authE *EndpointServiceUser) releaseAppointment(w http.ResponseWriter, r *http.Request) {

	log.Println("[EndpointServiceUser::releaseAppointment]")
	var data RequestBodyReleaseAppointment

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("[EndpointServiceUser::releaseAppointment] Data %+v", data)

	if err != nil {

		status := http.StatusBadRequest
		http.Error(w, http.StatusText(status), status)
		return
	}

	err = authE.world.ReleaseAppointment(&data.Token, data.AppointmentId)

	responseJSON := struct {
		Is    bool
		Error error
	}{}
	if err != nil {
		status := http.StatusBadRequest
		responseJSON.Is = false
		responseJSON.Error = err
		http.Error(w, http.StatusText(status), status)
	} else {
		responseJSON.Is = true
	}

	responseString, err := json.Marshal(responseJSON)

	w.Write([]byte(responseString))

	return

}
