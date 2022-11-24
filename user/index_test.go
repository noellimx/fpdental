package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var path string = "./users_test.json"

func TestLoadUsersExtracted(t *testing.T) {

	arr, err := LoadUsersExtracted(path)

	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, 1, len(arr), "len(loaded)")
	assert.Equal(t, "u1", arr[0].Name, "desc of 1st loaded")
	assert.Equal(t, "e3c0f872-0db1-418e-99d3-2221a9afd0cf", arr[0].Id, "id of 1st loaded")

}

