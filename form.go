package backend

import (
	"fmt"
	"net/url"
)

func ListForm(token, name string) (data []map[string]interface{}, err error) {
	qs := url.Values{}
	if len(name) > 0 {
		qs.Add("name", name)
	}

	err = Get(token, fmt.Sprintf("/form?%s", qs.Encode()), &data)
	return
}
