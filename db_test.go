package backend_test

import (
	"testing"

	"github.com/staticbackendhq/backend-go"
)

func TestList(t *testing.T) {
	// we add a tmp doc
	var insertedTask Task
	if err := backend.Create(token, "tasks", Task{Name: "test"}, &insertedTask); err != nil {
		t.Error(err)
	}

	var results []Task
	meta, err := backend.List(token, "tasks", &results, nil)
	if err != nil {
		t.Error(err)
	} else if len(results) < 1 {
		t.Errorf("expected tasks repo to have 1 doc got %d", len(results))
	} else if meta.Total < 1 {
		t.Errorf("expected total to be >= 1 goc %d", meta.Total)
	} else if results[len(results)-1].ID != insertedTask.ID {
		t.Errorf("expected last task ID to be %s got %s", insertedTask.ID, results[len(results)-1].ID)
	}
}
