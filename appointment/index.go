package appointment

import (
	"encoding/json"
	"errors"
	"fpdental/utils"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
)

type KeyAppointments string
type Appointments struct {
	m map[KeyAppointments]*Appointment
}

type Appointment struct {
	Id uuid.UUID

	Description string `json:"description"`
}

type AppointmentExtracted struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Booker      string `json:"booker"`
}

func Transform(ap *AppointmentExtracted) (*Appointment, error) {

	id, err := uuid.Parse(ap.Id)

	if err != nil {
		return nil, utils.ErrorTODO
	}
	return &Appointment{Id: id, Description: ap.Description}, nil
}
func (ap *Appointment) idString() string {
	return ap.Id.String()
}

func (aps *Appointments) AsArray() []*Appointment {

	arr := []*Appointment{}

	for _, a := range aps.m {

		arr = append(arr, a)
	}
	return arr
}

func (aps *Appointments) Log() {

	for i, ap := range aps.m {

		log.Printf("appns[%s] <-  %+v ", i, ap)

	}
}
func newAppointment() *Appointment {
	return &Appointment{Id: uuid.New()}
}

func NewAppointments() *Appointments {
	return &Appointments{m: make(map[KeyAppointments]*Appointment)}
}

func (aps *Appointments) Add(ap *Appointment) {
	aps.m[KeyAppointments(ap.idString())] = ap
}

func (aps *Appointments) GetById(id string) (*Appointment, bool) {

	log.Printf("[aps.GetById] %s", id)
	ap, found := aps.m[KeyAppointments(id)]
	log.Printf("[aps.GetById] %t", found)
	log.Printf("[aps.GetById] %+v", ap)

	return ap, found
}
func (aps *Appointments) Count() int {
	return len(aps.m)
}

var ErrorLoadAppointmentsFileOpenFail = errors.New("appointment-load: file open failed")
var ErrorLoadAppointmentsFileRead = errors.New("appointment-load: file read failed")
var ErrorLoadAppointmentsUnmarshal = errors.New("appointment-load: unmarshal failed")

func ExtractFromPath(path string) ([]*AppointmentExtracted, error) {

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
