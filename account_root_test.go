package backend_test

import (
	"fmt"
	"testing"
	"time"

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

func TestSetRole(t *testing.T) {
	me, err := currentAdminUser()
	if err != nil {
		t.Fatal(err)
	}

	email := fmt.Sprintf("setrole_%d@test.com", time.Now().UnixNano())
	user, err := backend.AddUser(token, email, "passwd1234")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := backend.RemoveUser(token, user.ID); err != nil {
			t.Fatalf("cleanup remove user: %v", err)
		}
	})

	const newRole = 25
	if err := backend.SetRole(token, me.AccountID, email, newRole); err != nil {
		t.Fatal(err)
	}

	users, err := backend.Users(token)
	if err != nil {
		t.Fatal(err)
	}

	for _, got := range users {
		if got.ID != user.ID {
			continue
		}
		if got.Role != newRole {
			t.Fatalf("expected role %d got %d", newRole, got.Role)
		}
		return
	}

	t.Fatalf("expected user %s to be present in users list", user.ID)
}

func currentAdminUser() (backend.CurrentUser, error) {
	users, err := backend.Users(token)
	if err != nil {
		return backend.CurrentUser{}, err
	}

	for _, user := range users {
		if user.Email == "admin@dev.com" {
			return backend.CurrentUser{
				AccountID: user.AccountID,
				UserID:    user.ID,
				Email:     user.Email,
				Role:      user.Role,
			}, nil
		}
	}

	return backend.CurrentUser{}, backend.ErrNoDocument
}
