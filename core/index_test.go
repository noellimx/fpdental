package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var wo = &WorldOpts{
	Paths: &Paths{users: "./users_test.json", appointments: "./appointments_test.json", credentials: "./credentials_test.json"},
}

var w *World

func TestInit(t *testing.T) {

	w = Init(wo)

	assert.Equal(t, 1, w.CountAppointmentsAvailable(), "availabile appointments")
	assert.Equal(t, 2, w.CountAppointmentUnavailable(), "unavailable appointments")

	adminUsername := "Admin"
	adminPassword := "Password"
	is := w.IsAuthenticated(adminUsername, adminPassword)

	if !is {
		t.Fatalf("should authenticate admin")
	}

}

func TestSignUp(t *testing.T) {
	w = Init(wo)

	newusername := "newbie"
	newpassword := "123"
	w.SignUp(newusername, newpassword)
	is := w.IsAuthenticated(newusername, newpassword)

	if !is {
		t.Fatalf("should authenticate new user")
	}

}

func TestGetUserAppointmentsJSONBye(t *testing.T) {

	w = Init(wo)

	got, err := w.GetUserAppointmentsJSONByte("u2")

	assert.Equal(t, []byte(`[{"Id":"d557c96c-ae2e-40a1-bc45-bd1b05e52f46","description":"a1"},{"Id":"652d5ac0-dc0f-4763-9b06-4c67bcacf6da","description":"a2"}]`), got, "a")

	if err != nil {
		t.Fatalf("json data mismatch")
	}
}
