package backend

import (
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
