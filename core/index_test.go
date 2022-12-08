package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var wo = &WorldOpts{
	Paths: &Paths{Users: "./users_test.json", Appointments: "./appointments_test.json", Credentials: "./credentials_test.json"},
}

var w *World

func TestInit(t *testing.T) {

	w = Init(wo)

	assert.Equal(t, 1, w.CountAppointmentsAvailable(), "availabile appointments")
	assert.Equal(t, 2, w.CountAppointmentUnavailable(), "unavailable appointments")

	adminUsername := "Admin"
	adminPassword := "Password"
	is, _ := w.IsAuthenticated(adminUsername, adminPassword)

	if !is {
		t.Fatalf("should authenticate admin")
	}

}

func TestSignUp(t *testing.T) {
	w = Init(wo)

	newusername := "newbie"
	newpassword := "123"
	w.SignUp(newusername, newpassword)
	is, _ := w.IsAuthenticated(newusername, newpassword)

	if !is {
		t.Fatalf("should authenticate new user")
	}

}

func TestGetUserAppointmentsJSONBye(t *testing.T) {

	w = Init(wo)

	targetUser := "u2"
	gotApp, err := w.GetUserAppointments(targetUser)

	if err != nil {

		t.Fatalf("Error retrieving user's appointments")
	}

	gotCount := gotApp.Count()
	wantCount := 2
	assert.Equal(t, wantCount, gotCount, "have count")
	gotByte, errByte := w.GetUserAppointmentsJSONByte(targetUser)

	assert.Equal(t, []byte(`[{"Id":"d557c96c-ae2e-40a1-bc45-bd1b05e52f46","description":"a1"},{"Id":"652d5ac0-dc0f-4763-9b06-4c67bcacf6da","description":"a2"}]`), gotByte, "a")

	if errByte != nil {
		t.Fatalf("json data mismatch")
	}
}

func TestGetUserAppointmentById(t *testing.T) {
	w = Init(wo)

	targetUserId := "u2"
	wantTargetUserHasBookedAppointmentId := "d557c96c-ae2e-40a1-bc45-bd1b05e52f46"

	app, found, err := w.GetUserAppointmentById(targetUserId, wantTargetUserHasBookedAppointmentId)
	if err != nil {
		t.Fatalf("Error retrieving user's appointments")
	}

	if !found {
		t.Fatalf("appointment should be booked by targetUser")
	}

	wantDescription := "a1"

	gotDescription := app.Description

	assert.Equal(t, wantDescription, gotDescription, "descriptions")

}

func TestReleaseAppointment(t *testing.T) {

	w := Init(wo)

	targetUserId := "u2"
	wantTargetUserHasBookedAppointmentId := "d557c96c-ae2e-40a1-bc45-bd1b05e52f46"

	app, found, err := w.GetUserAppointmentById(targetUserId, wantTargetUserHasBookedAppointmentId)
	if err != nil {
		t.Fatalf("Error retrieving user's appointments")
	}

	if !found {
		t.Fatalf("appointment should be booked by targetUser")
	}

	wantDescription := "a1"

	gotDescription := app.Description

	assert.Equal(t, wantDescription, gotDescription, "descriptions")

	err = w.ReleaseAppointment(targetUserId, wantTargetUserHasBookedAppointmentId)

	if err != nil {
		t.Fatalf(err.Error())
	}

}
