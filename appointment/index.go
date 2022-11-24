package appointment

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
)

type Appointment struct {
	id uuid.UUID

	Description string `json:"description"`
}

type AppointmentExtracted struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

func (ap *Appointment) idString() string {
	return ap.id.String()
}

func newAppointment() *Appointment {
	return &Appointment{id: uuid.New()}
}

type Appointments struct {
	m map[string]*Appointment
}

func NewAppointments() *Appointments {
	return &Appointments{m: make(map[string]*Appointment)}
}

func (aps *Appointments) Add(ap *Appointment) {
	aps.m[ap.idString()] = ap
}

func (aps *Appointments) Count() int {
	return len(aps.m)
}

var ErrorLoadAppointmentsFileOpenFail = errors.New("appointment-load: file open failed")
var ErrorLoadAppointmentsFileRead = errors.New("appointment-load: file read failed")
var ErrorLoadAppointmentsUnmarshal = errors.New("appointment-load: unmarshal failed")

func LoadAppointmentsExtracted(path string) ([]*AppointmentExtracted, error) {

	jsonFile, err := os.Open(path)

	if err != nil {
		return nil, ErrorLoadAppointmentsFileOpenFail
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, ErrorLoadAppointmentsFileRead
	}
	var appointmentsExtractedArray []*AppointmentExtracted

	err = json.Unmarshal(byteValue, &appointmentsExtractedArray)
	if err != nil {
		return nil, ErrorLoadAppointmentsUnmarshal
	}
	return appointmentsExtractedArray, nil

}

func loadAppointments(path string) error {

	log.Printf("[LoadAppointments] path <- %s", path)

	_, err := LoadAppointmentsExtracted(path)

	if err != nil {
		return err
	}

	return nil

}
