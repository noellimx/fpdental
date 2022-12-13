package core

import (
	"encoding/json"
	"errors"
	"fpdental/appointment"
	"fpdental/auth"
	"fpdental/user"
	"fpdental/utils"
	"log"

	"github.com/google/uuid"
)

type Paths struct {
	Users        string
	Appointments string
	Credentials  string
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

	w.PatientsOrReceptionist[w.Receptionist.Name] = newPatientOrReceptionnist(w.Receptionist)

	return nil

}

type PatientsOrReceptionist map[string]*PatientOrReceptionist
type World struct {
	PatientsOrReceptionist

	Receptionist *user.User
	*WorldOpts
	Auth *auth.AuthService
}

func (w *World) loadPatients() error {

	log.Println("[::loadPatients]")

	us, err := user.LoadUsers(w.WorldOpts.Paths.Users)

	if err != nil {

		log.Printf("[ERROR] [::loadPatients] %s\n", err)
		return err
	}

	for _, u := range us {
		w.PatientsOrReceptionist[u.Name] = newPatientOrReceptionnist(u)
	}

	return nil

}
func (w *World) RemoveUserSessions(token *auth.Token, userSessions []*auth.UserSessionsBE) {
	w.Auth.RemoveUserSessions(token, userSessions)
}
func (w *World) loadAppointments() error {
	path := w.WorldOpts.Paths.Appointments
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

		if w.PatientsOrReceptionist[apE.Booker] == nil {
			log.Panicf("[w.loadAppointments] User not found:%s", apE.Booker)
		}

		w.PatientsOrReceptionist[apE.Booker].Appointments.Add(ap)

	}
	return nil
}

func (w *World) CountAppointmentsAvailable() int {
	return w.PatientsOrReceptionist[w.Receptionist.Name].Appointments.Count()
}

func (w *World) CountAppointmentUnavailable() int {
	sum := 0

	for _, pOr := range w.PatientsOrReceptionist {
		if pOr.User != w.Receptionist {
			sum += pOr.Appointments.Count()
		}
	}

	return sum
}

func (w *World) loadCredentials() error {
	auth := auth.NewAuth()

	err := auth.InitCredentials(w.WorldOpts.Paths.Credentials)

	if err != nil {
		log.Println("[loadCredentials] error")
		return err
	}
	w.Auth = auth

	return nil

}

func (w *World) IsAuthenticated(username string, password string) (is bool, token *auth.Token) {

	request := &auth.UserCredentialClear{Username: username, Password: auth.PasswordUnhashed(password)}

	return w.Auth.IsAuthUnPw(request.Username, request.Password)
}

func (w *World) IsValidToken(token *auth.Token) bool {
	return w.Auth.IsAssociatedToken(token)
}

func newReceptionist() *user.User {

	return user.NewUser("", uuid.New())

}
func Init(wo *WorldOpts) *World {

	log.Printf("[World::init]")

	world := &World{WorldOpts: wo, PatientsOrReceptionist: make(PatientsOrReceptionist)}

	world.Receptionist = newReceptionist()

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

func (w *World) getUserAppointments(username string) (*appointment.Appointments, error) {

	p := w.PatientsOrReceptionist[username]

	if p == nil {
		return nil, utils.ErrorTODO
	}

	return p.Appointments, nil
}

var ErrorUnauthenticated = errors.New(("Unauthenticated"))
var ErrorUserNotFound = errors.New(("User Not Found"))

func (w *World) GetUserAppointments(token *auth.Token) (*appointment.Appointments, error) {
	is := w.IsValidToken(token)

	if is {

		p := w.PatientsOrReceptionist[token.Username]
		if p == nil {
			return nil, ErrorUserNotFound
		}
		return p.Appointments, nil
	} else {
		return nil, ErrorUnauthenticated
	}

}

func (w *World) GetAvailableAppointments() (*appointment.Appointments, error) {

	p := w.PatientsOrReceptionist[w.Receptionist.Name]
	if p == nil {
		return nil, ErrorUserNotFound
	}
	return p.Appointments, nil

}

func (w *World) GetUserAppointmentsCount(username string) (int, error) {

	aps, err := w.getUserAppointments(username)

	if err != nil {
		return 0, err
	}

	return aps.Count(), nil
}
func (w *World) GetUserAppointmentsJSONByte(username string) ([]byte, error) {

	aps, err := w.getUserAppointments(username)

	if err != nil {
		return nil, err
	}

	appsByte, err := json.Marshal(aps.AsSlice())

	if err != nil {
		return nil, err
	}

	return appsByte, nil

}

func (w *World) GetUserAppointmentById(username, apppointmentId string) (*appointment.Appointment, bool, error) {

	aps, err := w.getUserAppointments(username)

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

	ap, err := w.PatientsOrReceptionist[userNameFrom].Appointments.Remove(appointmentId)

	if err != nil {
		return err
	}

	w.PatientsOrReceptionist[userNameTo].Add(ap) // UNSAFE: Error Check Missing
	log.Printf("[transferAppointmentUNSAFE] #%s from %s to %s ", userNameFrom, userNameTo, appointmentId)
	return nil
}

func (w *World) transferToReceptionist(userNameFrom, appointmentId string) error {

	if w.Receptionist == nil {
		return ErrorReceptionistNotFound
	}

	w.transferAppointmentUNSAFE(userNameFrom, w.Receptionist.Name, appointmentId)

	return nil
}
func (w *World) releaseAppointmentUNSAFE(username, appointmentId string) error {
	return w.transferToReceptionist(username, appointmentId)
}

func (w *World) ReleaseAppointment(token *auth.Token, appointmentId uuid.UUID) error {

	is := w.Auth.IsAssociatedToken(token)

	if !is {

		log.Printf("[Error] [::ReleaseAppointment]")
		return utils.ErrorTODO
	}
	return w.transferToReceptionist(token.Username, appointmentId.String())
}

func (w *World) transferFromReceptionist(userNameTo, appointmentId string) error {
	if w.Receptionist == nil {
		return ErrorReceptionistNotFound
	}

	w.transferAppointmentUNSAFE(w.Receptionist.Name, userNameTo, appointmentId)

	return nil
}
func (w *World) BookAppointment(token *auth.Token, appointmentId uuid.UUID) error {
	log.Println("[BookAppointment]")
	is := w.Auth.IsAssociatedToken(token)

	if !is {
		log.Printf("[Error] [::BookAppointment]")
		return utils.ErrorTODO
	}
	return w.transferFromReceptionist(token.Username, appointmentId.String())
}

func (w *World) GetUserSessionsAll(token *auth.Token) ([]*auth.UserSessionsBE, error) {
	log.Println("[world::GetUserSessionsAll]")
	is := w.Auth.IsAssociatedToken(token)

	if !is {
		log.Printf("[Error] [::GetUserSessionsAll]")
		return nil, utils.ErrorTODO
	}
	return w.Auth.GetUserSessionsAll(), nil
}
