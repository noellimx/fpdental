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
	tokens   map[string]*Token
}

type UserSessionsBE struct {
	Username string
	TokenIds []string
}

func (u *UserCredential) IsTokenInPossession(token *Token) bool {

	return u.tokens[token.Id] != nil

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

type UserCredentials struct {
	Map map[string]*UserCredential
}
type AuthService struct {
	UserCredentials
	AdminUsername string
}

func (ucS UserCredentials) AsUserSessions() []*UserSessionsBE {

	userSessions := []*UserSessionsBE{}

	for username, value := range ucS.Map {

		us := &UserSessionsBE{
			Username: username,
		}

		_tokenIds := []string{}
		for _, token := range value.tokens {
			_tokenIds = append(_tokenIds, token.Id)
		}

		us.TokenIds = _tokenIds

		log.Printf("[::AsUserSessions] %v", us)

		userSessions = append(userSessions, us)
	}

	return userSessions
}
func (a *AuthService) GetUserSessionsAll() []*UserSessionsBE {

	return a.UserCredentials.AsUserSessions()

}
func (a *AuthService) CountCredentials() int {

	return len(a.UserCredentials.Map)
}
func (a *AuthService) IsAssociatedToken(token *Token) bool {
	log.Printf("[IsAssociatedToken] %s", token.Username)

	credential := a.UserCredentials.Map[token.Username]

	if credential != nil {
		return credential.IsTokenInPossession(token)
	}

	return false

}

func (a *AuthService) GetCredentials(name string) *UserCredential {
	log.Printf("[GetCredentials] %d", a.CountCredentials())

	log.Printf("[GetCredentials] GOT %s -> %v", name, a.UserCredentials.Map[name])

	return a.UserCredentials.Map[name]
}

func (a *AuthService) RegisterCredential(uc *UserCredentialClear) error {

	if a.UserCredentials.Map[uc.Username] != nil {
		return ErrorUsernameTaken
	}
	a.UserCredentials.Map[uc.Username] = &UserCredential{username: uc.Username, password: hashPassword(uc.Password)}
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
			token := authS.GenerateAndAssociateTokenToUser(username, ROLE_ADMIN, "")
			return is, token
		} else {
			token := authS.GenerateAndAssociateTokenToUser(username, ROLE_GENERAL, "")
			return is, token
		}
	} else {
		return is, nil
	}

}

func (auth *AuthService) GenerateAndAssociateTokenToUser(username string, role Role, expiry string) *Token {
	token := GenerateToken(username, role, expiry)
	auth.UserCredentials.Map[token.Username].tokens[token.Id] = token
	return token
}

func NewCredential(username string, passwordHashed PasswordHashed) *UserCredential {

	return &UserCredential{
		username: username,
		password: passwordHashed,
		tokens:   make(map[string]*Token),
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
		c_secure := NewCredential(c.Username, passwordHashed)

		credentialsMap[c.Username] = c_secure
		log.Printf("[InitCredentials] %s", c.Username)
	}

	auth.UserCredentials.Map = credentialsMap

	return nil

}

func (a *AuthService) IsAssociatedTokenAdmin(token *Token) bool {
	log.Printf("[IsAssociatedTokenAdmin] %s", token.Username)

	credential := a.UserCredentials.Map[token.Username]

	if credential != nil {
		return credential.IsTokenInPossession(token) && token.Username == ADMIN_USERNAME
	}

	return false

}

func (auth *AuthService) RemoveUserSessions(token *Token, userSessions []*UserSessionsBE) {

	is := auth.IsAssociatedTokenAdmin(token)
	if !is {
		return
	}

	for _, us := range userSessions {
		username := us.Username
		for _, id := range us.TokenIds {
			delete(auth.UserCredentials.Map[username].tokens, id)
		}
	}

}
func NewAuth() *AuthService {
	a := &AuthService{}
	a.UserCredentials.Map = make(map[string]*UserCredential) // TODO: To remove this line and test. Map Initialization should be done at ::InitCredentials
	a.AdminUsername = ADMIN_USERNAME
	return a
}
