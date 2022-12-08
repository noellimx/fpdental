package auth

import "github.com/google/uuid"

type Token struct {
	Id       string
	Username string
	Role     Role
	Expiry   string
}

func GenerateToken(username string, role Role, expiry string) *Token {

	id := uuid.New().String()
	return &Token{
		Id:       id,
		Username: username,
		Role:     role,
		Expiry:   "",
	}
}
