package backend

import (
	"encoding/json"
	"log"
	"time"
)

// CacheGet retrieves a cache value
func CacheGet(token, key string, v interface{}) error {
	var s string
	if err := Get(token, "/sudo/cache?key="+key, &s); err != nil {
		return err
	}
	return json.Unmarshal([]byte(s), v)
}

// CacheSet sets a cache value
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

// WorkerTask is the function type needed for work queue action
type WorkerTask func(val string)

// WorkerQueue monitors a work queue each 5 seconds.
// If there's new work available it will call the WorkerTask function back.
// This function should be ran concurrently i.e. go backend.WorkerQueue()
func WorkerQueue(token, key string, worker WorkerTask) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			go checkQueue(token, key, worker)
		}
	}
}

func checkQueue(token, key string, worker WorkerTask) {
	var s string
	if err := Get(token, "/sudo/cache?type=queue&key="+key, &s); err != nil {
		log.Println("error while checking for worker queue: ", err)
	}

	if len(s) > 0 {
		worker(s)
	}
}

// QueueWork adds a work queue value that will be dequeue via WorkerQueue
func QueueWork(token, key, value string) error {
	data := new(struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	})
	data.Key = key
	data.Value = value
	data.Type = "queue"
	var ok bool
	return Post(token, "/sudo/cache", data, &ok)
}
