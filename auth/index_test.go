package auth

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var authT *Auth = &Auth{}

// TODO: automatic teardown after every test
func teardown() {
	authT = NewAuth()
}

var path_rel string = "./credentials_test.json"

func TestInitCredentialsCount(t *testing.T) {

	want := 0
	got := authT.CountCredentials()

	assert.Equal(t, want, got, "len(credentials)")

	path_abs, err := filepath.Abs(path_rel)

	if err != nil {
		t.Fatalf("please check valid path relative: %s absolute %s", path_rel, path_abs)
	}

	t.Logf("[TestInitCredentials] %s %s", path_rel, path_abs)

	authT.InitCredentials(path_abs) // 1.5:be

	wantAdminCount := 1
	wantUserCount := 1
	want = wantAdminCount + wantUserCount
	got = authT.CountCredentials()
	assert.Equal(t, want, got, "len(credentials) after init. ")

	teardown()
}

func TestPasswordHash(t *testing.T) {

	var p PasswordUnhashed = "a"

	hp := PasswordHashed(fmt.Sprintf("%x", sha256.Sum256([]byte(p))))

	assert.Equal(t, hp, hashPassword(p), "error")

}

func TestAuthentications(t *testing.T) {

	uc := &UserCredentialClear{Username: "someu", Password: "somep"}

	wantAuth := false
	gotAuth := authT.IsAuth(uc.Username, uc.Password)

	assert.Equal(t, wantAuth, gotAuth, "Before storing credentials - isauth")

	authT.RegisterCredential(uc) // 1.1:be // 1.2:be
	wantAuth = true
	gotAuth = authT.IsAuth(uc.Username, uc.Password) // 1.3:be // 1.4:be

	assert.Equal(t, wantAuth, gotAuth, "After storing credentials - isauth")

}
