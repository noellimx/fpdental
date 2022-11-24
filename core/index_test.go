package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {

	wo := &WorldOpts{
		Paths: &Paths{users: "./users_test.json", appointments: "./appointments_test.json"},
	}

	w := Init(wo)

	assert.Equal(t, 1, w.CountAppointmentsAvailable(), "availabile appointments")
	assert.Equal(t, 2, w.CountAppointmentUnavailable(), "unavailable appointments")

}
