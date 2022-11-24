package appointment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndOwn(t *testing.T) {
	aps := NewAppointments()
	assert.Equal(t, aps.Count(), 0, "start with nil appointments")

	ap := newAppointment()

	aps.Add(ap)

	assert.Equal(t, aps.Count(), 1, "after 1 add")
}

func TestLoadAppointmentsExtracted(t *testing.T) {

	path := "./appointments_test.json"

	arr, err := LoadAppointmentsExtracted(path)

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, 3, len(arr), "len(loaded)")
	assert.Equal(t, "a1", arr[0].Description, "desc of 1st loaded")
	assert.Equal(t, "d557c96c-ae2e-40a1-bc45-bd1b05e52f46", arr[0].Id, "id of 1st loaded")

}
