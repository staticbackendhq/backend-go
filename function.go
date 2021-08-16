package backend

import "time"

type Function struct {
	ID           string       `json:"id"`
	FunctionName string       `json:"name"`
	TriggerTopic string       `json:"trigger"`
	Code         string       `json:"code"`
	Version      int          `json:"version"`
	LastUpdated  time.Time    `json:"lastUpdated"`
	LastRun      time.Time    `json:"lastRun"`
	History      []RunHistory `json:"history"`
}

type RunHistory struct {
	ID        string    `json:"id"`
	Version   int       `json:"version"`
	Started   time.Time `json:"started"`
	Completed time.Time `json:"completed"`
	Success   bool      `json:"success"`
	Output    []string  `json:"output"`
}

func AddFunction(token string, fn Function) error {
	return Post(token, "/fn/add", fn, nil)
}

func ListFunctions(token string) (results []Function, err error) {
	err = Get(token, "/fn", &results)
	return
}

func UpdateFunction(token string, fn Function) error {
	return Post(token, "/fn/update", fn, nil)
}

func DeleteFunction(token, name string) error {
	return Get(token, "/fn/del/"+name, nil)
}

func FunctionInfo(token, name string) (fn Function, err error) {
	err = Get(token, "/fn/info/"+name, &fn)
	return
}
