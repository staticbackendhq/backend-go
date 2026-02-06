package backend_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/staticbackendhq/backend-go"
)

func TestRegisterAndLogin(t *testing.T) {
	email, pass := fmt.Sprintf("unit_%d@test.com", time.Now().UnixNano()), "unit"

	if _, err := backend.Register(email, pass); err != nil {
		t.Error(err)
	}

	authToken, err := backend.Login(email, pass)
	if err != nil {
		t.Error(err)
	}

	if _, err := backend.AddUser(authToken, "user2@ok.com", "passwd1234"); err != nil {
		t.Fatal(err)
	}

	users, err := backend.Users(authToken)
	if err != nil {
		t.Fatal(err)
	}

	var userID string
	for _, user := range users {
		if user.Email == "user2@ok.com" {
			userID = user.UserID
			break
		}
	}

	t.Log("userID", userID)
	if err := backend.RemoveUser(token, userID); err != nil {
		t.Fatal(err)
	}

	users, err = backend.Users(token)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, user := range users {
		if user.UserID == userID {
			found = true
			break
		}
	}

	if found {
		t.Error("found the deleted user?")
	}
}
