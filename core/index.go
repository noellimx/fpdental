package core

import (
	"fpdental/appointment"
	"fpdental/auth"
	"fpdental/user"
	"log"

	"github.com/google/uuid"
)

type Paths struct {
	users        string
	appointments string
	credentials  string
}
type WorldOpts struct {
	*Paths
}

type PatientOrReceptionist struct {
	*user.User
	*appointment.Appointments
}

func newPatientOrReceptionnist(u *user.User) *PatientOrReceptionist {
	return &PatientOrReceptionist{User: u, Appointments: appointment.NewAppointments()}
}

var receptionist = user.NewUser("", uuid.New())

func (w *World) loadReceptionist() error {

	w.PatientOrReceptionists[receptionist.Name] = newPatientOrReceptionnist(receptionist)

	return nil

}

type PatientOrReceptionists map[string]*PatientOrReceptionist
type World struct {
	PatientOrReceptionists
	*WorldOpts
	*auth.Auth
}

func (w *World) loadPatients() error {
	us, err := user.LoadUsers(w.WorldOpts.Paths.users)

	if err != nil {
		return err
	}

	for _, u := range us {
		w.PatientOrReceptionists[u.Name] = newPatientOrReceptionnist(u)
	}

	return nil

}

func (w *World) loadAppointments() error {
	path := w.WorldOpts.Paths.appointments
	log.Printf("[w.loadAppointments] path <- %s", path)

	apsE, err := appointment.ExtractFromPath(path)
	if err != nil {
		return err
	}
	for _, apE := range apsE {
		ap, err := appointment.Transform(apE)
		if err != nil {
			return err
		}

		if w.PatientOrReceptionists[apE.Booker] == nil {
			log.Panicf(":%s", apE.Booker)
		}

		w.PatientOrReceptionists[apE.Booker].Appointments.Add(ap)

	}
	return nil
}

func (w *World) CountAppointmentsAvailable() int {
	return w.PatientOrReceptionists[receptionist.Name].Appointments.Count()
}

func (w *World) CountAppointmentUnavailable() int {
	sum := 0

	for _, pOr := range w.PatientOrReceptionists {
		if pOr.User != receptionist {
			sum += pOr.Appointments.Count()
		}
	}

	return sum
}

func (w *World) loadCredentials() error {
	auth := auth.NewAuth()

	err := auth.InitCredentials(w.WorldOpts.Paths.credentials)

	if err != nil {
		return err
	}

	w.Auth = auth

	return nil

}

func (w *World) IsAuthenticated(username string, password string) (is bool) {

	request := &auth.UserCredentialClear{Username: username, Password: auth.PasswordUnhashed(password)}

	return w.Auth.IsAuth(request.Username, request.Password)
}

func Init(wo *WorldOpts) *World {

	world := &World{WorldOpts: wo, PatientOrReceptionists: make(PatientOrReceptionists)}

	world.loadReceptionist()
	world.loadPatients()

	world.loadAppointments()

	world.loadCredentials()
	return world
}

func (w *World) SignUp(username, password string) error {
	cred := auth.NewUserCredentialClear(username, password)
	err := w.Auth.RegisterCredential(cred)

	if err != nil {
		return err
	}
	return nil
}
