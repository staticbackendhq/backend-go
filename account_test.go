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
			userID = user.ID
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
		if user.ID == userID {
			found = true
			break
		}
	}

	if found {
		t.Error("found the deleted user?")
	}
}

func TestEmailExists(t *testing.T) {
	exists, err := backend.EmailExists("admin@dev.com")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected admin@dev.com to exist")
	}

	missingEmail := fmt.Sprintf("missing_%d@test.com", time.Now().UnixNano())
	exists, err = backend.EmailExists(missingEmail)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("expected %s to not exist", missingEmail)
	}
}

func TestChangeEmail(t *testing.T) {
	oldEmail := fmt.Sprintf("change_email_old_%d@test.com", time.Now().UnixNano())
	newEmail := fmt.Sprintf("change_email_new_%d@test.com", time.Now().UnixNano())
	pass := "passwd1234"

	user, err := backend.AddUser(token, oldEmail, pass)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := backend.RemoveUser(token, user.ID); err != nil {
			t.Fatalf("cleanup remove user: %v", err)
		}
	})

	userToken, err := backend.Login(oldEmail, pass)
	if err != nil {
		t.Fatal(err)
	}

	if err := backend.ChangeEmail(userToken, newEmail); err != nil {
		t.Fatal(err)
	}

	me, err := backend.Me(userToken)
	if err != nil {
		t.Fatal(err)
	}
	if me.Email != newEmail {
		t.Fatalf("expected email %s got %s", newEmail, me.Email)
	}

	if _, err := backend.Login(newEmail, pass); err != nil {
		t.Fatal(err)
	}
}
