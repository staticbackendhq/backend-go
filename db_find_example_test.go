package backend_test

import (
	"fmt"
	"log"

	"github.com/staticbackendhq/backend-go"
)

type TaskItem struct {
	ID        string `json:"id"`
	AccountID string `json:"accountId"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
}

func ExampleFind() {
	var filters []backend.QueryItem

	// you may add multiple criteria, they're AND combined
	// No OR available at this time.
	filters = append(filters, backend.QueryItem{
		Field: "done",
		Op:    backend.QueryEqual,
		Value: true,
	})

	params := &backend.ListParams{
		Page:       1,    // current result page
		Size:       5,    // number of item per page
		Descending: true, // sort by created time descending
	}

	var matches []TaskItem
	result, err := backend.Find("session-token", "tasks", filters, &matches, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matches", result.Total, "Data", matches)
}
