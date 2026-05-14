package backend_test

import (
	"testing"

	"github.com/staticbackendhq/backend-go"
)

func TestPublishMessage(t *testing.T) {
	fakeTask := Task{
		Name: "not real",
		Done: true,
	}

	if err := backend.Publish(rootToken, "test-channel", "test-type", fakeTask); err != nil {
		t.Fatal(err)
	}
}
