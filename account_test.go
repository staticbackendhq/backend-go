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

	if _, err := backend.Login(email, pass); err != nil {
		t.Error(err)
	}
}
