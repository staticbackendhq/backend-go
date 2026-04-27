package backend_test

import (
	"testing"

	"github.com/staticbackendhq/backend-go"
)

const rootToken = "safe-to-use-in-dev-root-token"

func TestSudoGetAuthTokenByUserID(t *testing.T) {
	me, err := currentAdminUser()
	if err != nil {
		t.Fatal(err)
	}

	authToken, err := backend.SudoGetAuthTokenByUserID(rootToken, me.AccountID, me.UserID)
	if err != nil {
		t.Fatal(err)
	}
	if authToken == "" {
		t.Fatal("expected auth token to be returned")
	}
}

func TestSudoGetUserByID(t *testing.T) {
	me, err := currentAdminUser()
	if err != nil {
		t.Fatal(err)
	}

	user, err := backend.SudoGetUserByID(rootToken, me.AccountID, me.UserID)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != me.UserID {
		t.Fatalf("expected user id %s got %s", me.UserID, user.ID)
	}
	if user.AccountID != me.AccountID {
		t.Fatalf("expected account id %s got %s", me.AccountID, user.AccountID)
	}
	if user.Email != me.Email {
		t.Fatalf("expected email %s got %s", me.Email, user.Email)
	}
}

func currentAdminUser() (backend.CurrentUser, error) {
	users, err := backend.Users(token)
	if err != nil {
		return backend.CurrentUser{}, err
	}

	for _, user := range users {
		if user.Email == "admin@dev.com" {
			return user, nil
		}
	}

	return backend.CurrentUser{}, backend.ErrNoDocument
}
