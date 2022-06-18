package backend

import (
	"fmt"
	"net/url"
)

// ListForm returns submissions for all or a specific form.
func ListForm(token, name string) (data []map[string]interface{}, err error) {
	qs := url.Values{}
	if len(name) > 0 {
		qs.Add("name", name)
	}

	err = Get(token, fmt.Sprintf("/form?%s", qs.Encode()), &data)
	return
}
