package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var wo = &WorldOpts{
	Paths: &Paths{users: "./users_test.json", appointments: "./appointments_test.json", credentials: "./credentials_test.json"},
}

func TestInit(t *testing.T) {

	w := Init(wo)

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

}
