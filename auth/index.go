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

var ADMIN_USERNAME = "Admin"

var ROLE_ADMIN Role = "Admin"
var ROLE_GENERAL Role = "General"

type Role string

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

type AuthService struct {
	UserCredentials map[string]*UserCredential
	AdminUsername   string
}

func (a *AuthService) CountCredentials() int {

	return len(a.UserCredentials)
}

func (a *AuthService) GetCredentials(name string) *UserCredential {
	log.Printf("[GetCredentials] %d", a.CountCredentials())

	log.Printf("[GetCredentials] GOT %s -> %v", name, a.UserCredentials[name])

	return a.UserCredentials[name]
}

func (a *AuthService) RegisterCredential(uc *UserCredentialClear) error {

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
func (authS *AuthService) IsAuthUnPw(username string, password PasswordUnhashed) (bool, *Token) {

	c := authS.GetCredentials(username)

	if c == nil {
		return false, nil
	}
	is := c.isAuthByPassword(password)

	if is {
		if username == authS.AdminUsername {
			token := GenerateToken(username, ROLE_ADMIN, "")
			return is, token
		} else {
			token := GenerateToken(username, ROLE_GENERAL, "")
			return is, token
		}
	} else {
		return is, nil
	}

}

func (auth *AuthService) InitCredentials(path string) error {

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
		log.Printf("[InitCredentials] %s", c.Username)
	}

	auth.UserCredentials = credentialsMap

	return nil

}

func NewAuth() *AuthService {

	a := &AuthService{}
	a.UserCredentials = make(map[string]*UserCredential)
	a.AdminUsername = ADMIN_USERNAME
	return a
}
