package backend

import (
	"testing"
)

func TestInitializeCredentials(t *testing.T) {
	InitCredentials()
	if !(len(credentials.User) > 0) {
		t.Error("Should get user configfile")
	}
	if !(len(credentials.Password) > 0) {
		t.Error("Should get user from configfile")
	}
}
