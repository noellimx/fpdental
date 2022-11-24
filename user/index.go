package user

import (
	"encoding/json"
	"fpdental/appointment"
	"fpdental/utils"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type User struct {
	id   uuid.UUID
	name string
	*appointment.Appointments
}

type UserExtracted struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Users map[string]*User

func newUser(name string, id uuid.UUID) *User {
	return &User{id: id, name: name, Appointments: appointment.NewAppointments()}
}

func transformUser(ue *UserExtracted) (*User, error) {
	username := ue.Name

	id, err := uuid.Parse(ue.Id)

	if err != nil {
		return nil, err
	}

	return &User{
		name: username,
		id:   id,
	}, nil

}

func LoadUsersExtracted(path string) ([]*UserExtracted, error) {

	jsonFile, err := os.Open(path)

	if err != nil {
		return nil, utils.ErrorTODO
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, utils.ErrorTODO
	}
	var appointmentsExtractedArray []*UserExtracted

	err = json.Unmarshal(byteValue, &appointmentsExtractedArray)
	if err != nil {
		return nil, utils.ErrorTODO
	}
	return appointmentsExtractedArray, nil
}

func LoadUsers(path string) (map[string]*User, error) {

	ues, err := LoadUsersExtracted(path)
	if err != nil {
		return nil, err
	}

	uemap := make(map[string]*User)
	for _, ue := range ues {
		u, err := transformUser(ue)

		if err != nil {
			return nil, err
		}

		uemap[u.name] = u
	}

	return uemap, nil

}
