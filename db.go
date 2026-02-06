package backend

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// Create adds a new document to a repository and returns the created document.
func Create(token, repo string, body interface{}, v interface{}) error {
	return Post(token, fmt.Sprintf("/db/%s", repo), body, v)
}

// CreateBulk creates multiple documents, useful when importing data.
func CreateBulk(token, repo string, body interface{}) (bool, error) {
	var status bool
	if err := Post(token, fmt.Sprintf("/db/%s?bulk=1", repo), body, &status); err != nil {
		return false, err
	}
	return status, nil
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

// GetByIDs returns matching docs by ids. Useful to prevent n+1 query
func GetByIDs(token, repo string, ids []string, v interface{}) error {
	return Post(token, fmt.Sprintf("/db/%s?ids=true", repo), ids, v)
}

// QueryItem used to perform query
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

// FindOne returns one document if it's found
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

	body := toFilterSlice(filters)

	u := fmt.Sprintf("/query/%s?%s", repo, qs.Encode())
	err = Post(token, u, body, &meta)
	if err != nil {
		return meta, fmt.Errorf("error executing the querying: %v", err)
	}

	return
}

// Update updates a document. Can be just a subset of the fields.
func Update(token, repo, id string, body interface{}, v interface{}) error {
	return Put(token, fmt.Sprintf("/db/%s/%s", repo, id), body, v)
}

// UpdateBulk updates multiple documents based on filter clauses
func UpdateBulk(token, repo string, filters []QueryItem, body interface{}) (n int, err error) {
	var data struct {
		Update  interface{}     `json:"update"`
		Clauses [][]interface{} `json:"clauses"`
	}
	data.Update = body
	data.Clauses = toFilterSlice(filters)

	err = Put(token, fmt.Sprintf("/db/%s?bulk=1", repo), data, &n)
	return
}

// Delete permanently delets a document
func Delete(token, repo, id string) error {
	return Del(token, fmt.Sprintf("/db/%s/%s", repo, id))
}

// DeleteBulk permanently deletes multiple documents matching filters
func DeleteBulk(token, repo string, filters []QueryItem) error {
	data := toFilterSlice(filters)

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	x := base64.StdEncoding.EncodeToString(b)
	uri := fmt.Sprintf("/db/%s?bulk=1&x=%s", repo, x)
	return Del(token, uri)
}

// SudoCreate adds a new document to a repository and returns the created document.
func SudoCreate(token, repo string, body interface{}, v interface{}) error {
	return Post(token, fmt.Sprintf("/sudo/%s", repo), body, v)
}

// SudoList returns a list of documents in a specific repository if a "root token" is used.
func SudoList(token, repo string, v interface{}, params *ListParams) (meta ListResult, err error) {
	meta.Results = v

	qs := url.Values{}
	if params != nil {
		qs.Add("page", strconv.Itoa(params.Page))
		qs.Add("size", strconv.Itoa(params.Size))
		if params.Descending {
			qs.Add("desc", "true")
		}
	}

	err = Get(token, fmt.Sprintf("/sudo/%s?%s", repo, qs.Encode()), &meta)
	return
}

// SudoSudoGetByID returns a specific document if a "root token" is provided.
func SudoGetByID(token, repo, id string, v interface{}) error {
	return Get(token, fmt.Sprintf("/sudo/%s/%s", repo, id), v)
}

// SudoGetByIDs returns matching documents by ids if root token is provided
func SudoGetByIDs(token, repo string, ids []string, v interface{}) error {
	return post(token, fmt.Sprintf("/sudo/%s?ids=true", repo, ids), v)
}

// SudoUpdate perform an update if a "root" token is specified
// This call cannot be done from JavaScript, only from a backend HTTP call.
//
// You can obtain this token via the CLI or web interface.
func SudoUpdate(token, repo, id string, body interface{}, v interface{}) error {
	return Put(token, fmt.Sprintf("/sudo/%s/%s", repo, id), body, v)
}

// SudoFind returns a slice of matching documents if a "root token" is provided.
func SudoFind(token, repo string, filters []QueryItem, v interface{}, params *ListParams) (meta ListResult, err error) {
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

	u := fmt.Sprintf("/sudoquery/%s?%s", repo, qs.Encode())
	err = Post(token, u, body, &meta)
	if err != nil {
		return meta, fmt.Errorf("error executing the querying: %v", err)
	}

	return
}

// SudoFindOne returns one document if it's found with a root token
func SudoFindOne(token, repo string, filters []QueryItem, v interface{}) error {
	meta, err := SudoFind(token, repo, filters, v, nil)
	if err != nil {
		return err
	} else if meta.Total == 0 {
		return ErrNoDocument
	} else if meta.Total > 1 {
		return ErrMultipleDocument
	}

	return nil
}

// Delete permanently delete a document via a root token
func SudoDelete(token, repo, id string) error {
	return Del(token, fmt.Sprintf("/sudo/%s/%s", repo, id))
}

// SudoListRepositories lists all database repositories if a "root token" is provided.
func SudoListRepositories(token string) ([]string, error) {
	var names []string
	if err := Get(token, "/sudolistall", &names); err != nil {
		return nil, fmt.Errorf("error getting repositories list: %v", err)
	}

	return names, nil
}

// Increase increases or decreases a field in a document based on n signed
func Increase(token, repo, id, field string, n int) error {
	var body = new(struct {
		Field string `json:"field"`
		Range int    `json:"range"`
	})
	body.Field = field
	body.Range = n

	var status bool
	return Put(token, fmt.Sprintf("/inc/%s/%s", repo, id), body, &status)
}

// SudoAddIndex creates a new database index on a specific field
func SudoAddIndex(token, repo, field string) error {
	var status bool

	url := fmt.Sprintf("/db/index?col=%s&field=%s", repo, field)
	return Post(token, url, nil, &status)
}

func toFilterSlice(filters []QueryItem) [][]interface{} {
	var body [][]interface{}
	for _, f := range filters {
		body = append(body, []interface{}{
			f.Field,
			f.Op,
			f.Value,
		})
	}
	return body
}

// Count returns the number of document in a repo matching the optional
// filters, which are same query filter as Query.
func Count(token, repo string, filters []QueryItem) (n int64, err error) {
	data := new(struct {
		Count int64 `json:"count"`
	})

	payload := toFilterSlice(filters)
	err = Post(token, fmt.Sprintf("/db/count/%s", repo), payload, &data)
	n = data.Count
	return
}

// Search returns the matching document of "repo" based on keywords
func Search(token, repo, keywords string, v interface{}) error {
	data := new(struct {
		Col      string `json:"col"`
		Keywords string `json:"keywords"`
	})
	data.Col = repo
	data.Keywords = keywords

	return Post(token, "/search", data, v)
}
