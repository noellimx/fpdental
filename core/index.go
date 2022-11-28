package core

import (
	"encoding/json"
	"fpdental/appointment"
	"fpdental/auth"
	"fpdental/user"
	"fpdental/utils"
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

func (w *World) loadReceptionist() error {

	w.PatientOrReceptionists[w.Receptionist.Name] = newPatientOrReceptionnist(w.Receptionist)

	return nil

}

type PatientOrReceptionists map[string]*PatientOrReceptionist
type World struct {
	PatientOrReceptionists

	Receptionist *user.User
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
	return w.PatientOrReceptionists[w.Receptionist.Name].Appointments.Count()
}

func (w *World) CountAppointmentUnavailable() int {
	sum := 0

	for _, pOr := range w.PatientOrReceptionists {
		if pOr.User != w.Receptionist {
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
	var receptionist = user.NewUser("", uuid.New())

	world.Receptionist = receptionist
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

func (w *World) GetUserAppointments(username string) (*appointment.Appointments, error) {

	p := w.PatientOrReceptionists[username]

	if p == nil {

		return nil, utils.ErrorTODO

	}

	return p.Appointments, nil

}

func (w *World) GetUserAppointmentsCount(username string) (int, error) {

	aps, err := w.GetUserAppointments(username)

	if err != nil {

		return 0, err
	}

	return aps.Count(), nil
}
func (w *World) GetUserAppointmentsJSONByte(username string) ([]byte, error) {

	aps, err := w.GetUserAppointments(username)

	if err != nil {
		return nil, err
	}

	appsByte, err := json.Marshal(aps.AsArray())

	if err != nil {

		return nil, err
	}

	return appsByte, nil

}

func (w *World) GetUserAppointmentById(username, apppointmentId string) (*appointment.Appointment, bool, error) {

	aps, err := w.GetUserAppointments(username)

	log.Printf("[GetUserAppoinment] %d", aps.Count())

	aps.Log()

	if err != nil {
		return nil, false, err
	}

	ap, found := aps.GetById(apppointmentId)

	return ap, found, nil

}

var ErrorAppointmentUserMismatch = utils.ErrorTODO
var ErrorReceptionistNotFound = utils.ErrorTODO

func (w *World) transferAppointmentUNSAFE(userNameFrom, userNameTo, appointmentId string) error {

	// UNSAFE: Not thread-safe.
	_, found, err := w.GetUserAppointmentById(userNameFrom, appointmentId)

	if err != nil {
		return err
	}

	if !found {
		return ErrorAppointmentUserMismatch
	}

	ap, err := w.PatientOrReceptionists[userNameFrom].Appointments.Remove(appointmentId)

	if err != nil {
		return err
	}

	w.PatientOrReceptionists[userNameTo].Add(ap) // UNSAFE: Error Check Missing

	return nil
}

func (w *World) transferToReceptionist(userNameFrom, appointmentId string) error {

	if w.Receptionist == nil {
		return ErrorReceptionistNotFound
	}

	w.transferAppointmentUNSAFE(userNameFrom, w.Receptionist.Name, appointmentId)

	return nil
}
func (w *World) ReleaseAppointment(username, appointmentId string) error {

	return w.transferToReceptionist(username, appointmentId)
}
