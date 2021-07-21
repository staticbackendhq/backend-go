package backend

import "encoding/json"

func CacheGet(token, key string, v interface{}) error {
	var s string
	if err := Get(token, "/sudo/cache?key="+key, &s); err != nil {
		return err
	}
	return json.Unmarshal([]byte(s), v)
}

func CacheSet(token, key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data := new(struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	})
	data.Key = key
	data.Value = string(b)
	var ok bool
	return Post(token, "/sudo/cache", data, &ok)
}
