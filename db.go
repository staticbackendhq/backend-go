package backend

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Create adds a new document to a repository and returns the created document.
func Create(token, repo string, body interface{}, v interface{}) error {
	return Post(token, fmt.Sprintf("/db/%s", repo), body, v)
}

// ListParams are used to page results and sort.
type ListParams struct {
	Page       int
	Size       int
	Descending bool
}

// ListResult is used for list document
type ListResult struct {
	Page     int         `json:"page"`
	PageSize int         `json:"size"`
	Total    int         `json:"total"`
	Results  interface{} `json:"results"`
}

// List returns a list of document in a specific repository.
func List(token, repo string, v interface{}, params *ListParams) (meta ListResult, err error) {
	meta.Results = v

	qs := url.Values{}
	if params != nil {
		qs.Add("page", strconv.Itoa(params.Page))
		qs.Add("size", strconv.Itoa(params.Size))
		if params.Descending {
			qs.Add("desc", "true")
		}
	}

	err = Get(token, fmt.Sprintf("/db/%s?%s", repo, qs.Encode()), &meta)
	return
}

// GetByID returns a specific document.
func GetByID(token, repo, id string, v interface{}) error {
	return Get(token, fmt.Sprintf("/db/%s/%s", repo, id), v)
}

type QueryItem struct {
	Field string
	Op    QueryOperator
	Value interface{}
}

type QueryOperator string

const (
	QueryEqual            QueryOperator = "=="
	QueryNotEqual         QueryOperator = "!="
	QueryLowerThan        QueryOperator = "<"
	QueryLowerThanEqual   QueryOperator = "<="
	QueryGreaterThan      QueryOperator = ">"
	QueryGreaterThanEqual QueryOperator = ">="
	QueryIn               QueryOperator = "in"
	QueryNotIn            QueryOperator = "!in"
)

var (
	ErrNoDocument       = errors.New("no document found")
	ErrMultipleDocument = errors.New("multiple documents found")
)

func FindOne(token, repo string, filters []QueryItem, v interface{}) error {
	meta, err := Find(token, repo, filters, v, nil)
	if err != nil {
		return err
	} else if meta.Total == 0 {
		return ErrNoDocument
	} else if meta.Total > 1 {
		return ErrMultipleDocument
	}

	return nil
}

// Find returns a slice of matching documents.
func Find(token, repo string, filters []QueryItem, v interface{}, params *ListParams) (meta ListResult, err error) {
	meta.Results = v

	qs := url.Values{}
	if params != nil {
		qs.Add("page", strconv.Itoa(params.Page))
		qs.Add("size", strconv.Itoa(params.Size))
		if params.Descending {
			qs.Add("desc", "true")
		}
	}

	var body [][]interface{}
	for _, f := range filters {
		body = append(body, []interface{}{
			f.Field,
			f.Op,
			f.Value,
		})
	}

	u := fmt.Sprintf("/query/%s?%s", repo, qs.Encode())
	err = Post(token, u, body, &meta)
	if err != nil {
		return meta, fmt.Errorf("error executing the querying: %v", err)
	}

	return
}

func Update(token, repo, id string, body interface{}) (status bool, err error) {
	err = Put(token, fmt.Sprintf("/db/%s/%s", repo, id), body, &status)
	return
}

// SudoUpdate perform an update if a "root" token is specified
// This call cannot be done from JavaScript, only from a backend HTTP call.
//
// You can obtain this token via the CLI or web interface.
func SudoUpdate(token, repo, id string, body interface{}) (status bool, err error) {
	err = Put(token, fmt.Sprintf("/sudo/%s/%s", repo, id), body, &status)
	return
}
