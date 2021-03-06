package backend_test

import (
	"testing"
	"time"

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

func TestFindOne(t *testing.T) {
	// we add a few docs
	var insertedTask Task
	if err := backend.Create(token, "tasks", Task{Name: "task 1"}, &insertedTask); err != nil {
		t.Error(err)
	}

	var insertedTask2 Task
	if err := backend.Create(token, "tasks", Task{Name: "task 2"}, &insertedTask2); err != nil {
		t.Error(err)
	}

	var filters []backend.QueryItem
	filters = append(filters, backend.QueryItem{
		Field: "name",
		Op:    backend.QueryEqual,
		Value: "task 2",
	})
	var doc []Task
	if err := backend.FindOne(token, "tasks", filters, &doc); err != nil {
		t.Error(err)
	} else if len(doc) != 1 {
		t.Errorf("expected 1 doc got %d", len(doc))
	} else if insertedTask2.ID != doc[0].ID {
		t.Errorf("found doc's id differ expected %s got %s", insertedTask2.ID, doc[0].ID)
	}
}

type ProductTypes string

const (
	ProductTypeDownload     ProductTypes = "download"
	ProductTypeCourse       ProductTypes = "course"
	ProductTypeSubscription ProductTypes = "subscription"
)

type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusPublished ProductStatus = "published"
	ProductStatusArchived  ProductStatus = "archived"
)

type Product struct {
	ID          string        `json:"id"`
	Hunter      Person        `json:"hunter"`
	Creators    []Person      `json:"creators"`
	Name        string        `json:"name"`
	Description string        `json:"desc"`
	Type        ProductTypes  `json:"type"`
	Category    string        `json:"category"`
	Topics      []string      ` json:"topics"`
	Status      ProductStatus ` json:"status"`
	ListImage   string        ` json:"listImage"`
	Likes       int           ` json:"likes"`
	Link        string        `json:"link"`
	Detail      string        ` json:"detail"`
	Slug        string        `json:"slug"`
	CreatedAt   time.Time     `json:"created"`
	UpdatedAt   time.Time     ` json:"updated"`
}

type Person struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Voter struct {
	ProductID string `json:"productId"`
	UserID    string `json:"userId"`
}

func TestWTH(t *testing.T) {
	backend.PublicKey = "5fc640ff4eea8b872b907f8d"

	products := make([]Product, 0)
	meta, err := backend.List("", "pub_posts_744_", &products, nil)
	if err != nil {
		t.Fatal(err)
	} else if meta.Total == 0 {
		t.Errorf("total is %d expected > 0", meta.Total)
	} else if meta.Results == nil {
		t.Error("Results is nil")
	}
}
