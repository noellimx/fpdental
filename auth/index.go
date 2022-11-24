package auth

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type PasswordUnhashed string
type PasswordHashed string

type UserCredential struct {
	password PasswordHashed
	username string
}

func (u *UserCredential) isAuthByPassword(attempt PasswordUnhashed) bool {

	pw := hashPassword(attempt)
	return pw == u.password
}

type UserCredentialClear struct {
	Password PasswordUnhashed `json:"password"`
	Username string           `json:"username"`
}

func NewUserCredentialClear(username, password string) *UserCredentialClear {

	return &UserCredentialClear{Username: username, Password: PasswordUnhashed(password)}
}

type Auth struct {
	UserCredentials map[string]*UserCredential
}

func (a *Auth) CountCredentials() int {

	return len(a.UserCredentials)
}

func (a *Auth) GetCredentials(name string) *UserCredential {
	log.Printf("%d", a.CountCredentials())

	for _, c := range a.UserCredentials {
		log.Printf("creds checking %s%s", c.username, c.password)

	}

	return a.UserCredentials[name]
}

func (a *Auth) RegisterCredential(uc *UserCredentialClear) error {

	if a.UserCredentials[uc.Username] != nil {

		return ErrorUsernameTaken
	}
	a.UserCredentials[uc.Username] = &UserCredential{username: uc.Username, password: hashPassword(uc.Password)}
	return nil
}

var ErrorInitCredentialsFileOpenFail = errors.New("auth-init-credentials: file open failed")
var ErrorInitCredentialsFileRead = errors.New("auth-init-credentials: file read failed")
var ErrorInitCredentialsUnmarshal = errors.New("auth-init-credentials: unmarshal failed")
var ErrorUsernameTaken = errors.New("register username taken")

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func hashPassword(s PasswordUnhashed) PasswordHashed {
	return PasswordHashed(hash(string(s)))
}
func (auth *Auth) IsAuth(username string, password PasswordUnhashed) bool {
	c := auth.GetCredentials(username)
	if c == nil {
		return false
	}

	return c.isAuthByPassword(password)
}

func (auth *Auth) InitCredentials(path string) error {

	log.Printf("[InitCredentials] path <- %s", path)
	jsonFile, err := os.Open(path)
	if err != nil {
		return ErrorInitCredentialsFileOpenFail
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return ErrorInitCredentialsFileRead
	}
	var credentialsArray []*UserCredentialClear

	err = json.Unmarshal(byteValue, &credentialsArray)
	if err != nil {
		return ErrorInitCredentialsUnmarshal
	}

	credentialsMap := make(map[string]*UserCredential)

	for _, c := range credentialsArray {
		passwordHashed := hashPassword(c.Password)

		c_secure := &UserCredential{
			username: c.Username,
			password: passwordHashed,
		}
		credentialsMap[c.Username] = c_secure

	}

	auth.UserCredentials = credentialsMap

	return nil

}

func NewAuth() *Auth {

	a := &Auth{}
	a.UserCredentials = make(map[string]*UserCredential)

	return a
}
