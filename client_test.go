package backend_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/staticbackendhq/backend-go"
)

var token string

func init() {
	backend.PublicKey = "dev_memory_pk"

	// this user is created by default in dev mod as an admin in the CLI
	t, err := backend.Login("admin@dev.com", "devpw1234")
	if err != nil {
		log.Fatal("unable to login to the backend: ", err)
	}

	token = t

	if err := cleanUp(); err != nil {
		// TODO: this could be better
		// if the tasks collection did not existed, skip
		if strings.Index(err.Error(), "collection not found") == -1 {
			log.Fatal(err)
		}
	}
}

func cleanUp() error {
	lp := &backend.ListParams{Size: 500, Page: 1}

	var tasks []Task
	if _, err := backend.List(token, "tasks", &tasks, lp); err != nil {
		return err
	}

	for _, task := range tasks {
		if err := backend.Delete(token, "tasks", task.ID); err != nil {
			return err
		}
	}

	return nil
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
