package backend_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/staticbackendhq/backend-go"
)

var token string

func init() {
	backend.PublicKey = "dev-memory-pk"

	// this user is created by default in dev mod as an admin
	t, err := backend.Login("admin@dev.com", "devpw1234")
	if err != nil {
		log.Fatal("unable to login to the backend: ", err)
	}

	token = t
}

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func TestGetRequest(t *testing.T) {
	// we add a tmp doc
	var insertedTask Task
	if err := backend.Create(token, "pub_test", Task{Name: "test"}, &insertedTask); err != nil {
		t.Error(err)
	}

	var s Task
	u := fmt.Sprintf("/db/pub_test/%s", insertedTask.ID)
	if err := backend.Get(token, u, &s); err != nil {
		t.Error(err)
	} else if insertedTask.ID != s.ID {
		t.Errorf("expected inserted to be %s got %s", insertedTask.ID, s.ID)
	}
}
