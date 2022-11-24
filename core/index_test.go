package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var wo = &WorldOpts{
	Paths: &Paths{users: "./users_test.json", appointments: "./appointments_test.json"},
}

func TestInit(t *testing.T) {

	w := Init(wo)

	assert.Equal(t, 1, w.CountAppointmentsAvailable(), "availabile appointments")
	assert.Equal(t, 2, w.CountAppointmentUnavailable(), "unavailable appointments")

	// is := w.IsAuthenticated("Admin", "Password")

}

func TestSignUp(t *testing.T) {

}
