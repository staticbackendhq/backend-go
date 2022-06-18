package backend_test

import (
	"fmt"
	"log"

	"github.com/staticbackendhq/backend-go"
)

func Example() {
	// it's required to set a PublicKey see: https://staticbackend.com/getting-started/local-dev
	backend.PublicKey = "Get from CLI, self-hosted or hosted svc email"

	// we're setting the Region to "dev" to use with the CLI
	// For self-hosted you may set your backend API URL directly:
	// backend.Region = "https://your-server.com/""
	backend.Region = backend.RegionLocalDev

	// we register a new user
	// token is the new user's session token, we need it for all subsequent
	// requests
	token, err := backend.Register("go@example.com", "my-safe-pw")
	if err != nil {
		log.Fatal(err)
	}

	// let's create a task in a tasks repository
	//
	// you'd typically create your struct normally, not like this.
	task := new(struct {
		ID        string `json:"id"`
		AccountID string `json:"accountId"`
		Title     string `json:"title"`
		Done      bool   `json:"done"`
	})

	task.Title = "My first to-do item"

	if err := backend.Create(token, "tasks", task, &task); err != nil {
		log.Fatal(err)
	}

	// task.ID and task.AccountID has been automatically set
	fmt.Println("ID", task.ID, "AccountID", task.AccountID)

	// get the same document to showcase the GetByID function
	if err := backend.GetByID(token, "tasks", task.ID, &task); err != nil {
		log.Fatal(err)
	}

	fmt.Println(task.Title)
}
