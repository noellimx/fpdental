package user

import (
	"encoding/json"
	"fpdental/utils"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

type User struct {
	Id   uuid.UUID
	Name string
}

type UserExtracted struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Users map[string]*User

func NewUser(name string, id uuid.UUID) *User {
	return &User{Id: id, Name: name}
}

func transformUser(ue *UserExtracted) (*User, error) {
	username := ue.Name

	id, err := uuid.Parse(ue.Id)

	if err != nil {
		return nil, err
	}

	return &User{
		Name: username,
		Id:   id,
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
	var usersExtractedArray []*UserExtracted

	err = json.Unmarshal(byteValue, &usersExtractedArray)
	if err != nil {
		return nil, utils.ErrorTODO
	}
	return usersExtractedArray, nil
}

func LoadUsers(path string) ([]*User, error) {

	ues, err := LoadUsersExtracted(path)
	if err != nil {
		return nil, err
	}

	uemap := []*User{}
	for _, ue := range ues {
		u, err := transformUser(ue)

		if err != nil {
			return nil, err
		}

		uemap = append(uemap, u)
	}

	return uemap, nil

}
