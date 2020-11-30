package backend_test

import (
	"testing"

	"github.com/staticbackendhq/backend-go"
)

func TestRegisterAndLogin(t *testing.T) {
	email, pass := "unit", "unit"

	tok, err := backend.Register(email, pass)
	if err != nil {
		t.Error(err)
	}

	tok2, err := backend.Login(email, pass)
	if err != nil {
		t.Error(err)
	} else if tok != tok2 {
		t.Errorf("register/login tokens does not match, expected %s got %s", tok, tok2)
	}
}
