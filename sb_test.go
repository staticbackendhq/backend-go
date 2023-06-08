package backend_test

import (
	"fmt"
	"testing"

	"github.com/staticbackendhq/backend-go"
)

func TestNewSysAccountBypassStripe(t *testing.T) {
	data, err := backend.NewSystemAccountBypassStripe("nosripe@d.com", "no-stripe-test-flag")
	if err != nil {
		t.Fatal(err)
	} else if len(data.AdminPassword) != 6 {
		t.Errorf("expected 6 pw length, got %d", len(data.AdminPassword))
	}
	fmt.Println(data)
}
