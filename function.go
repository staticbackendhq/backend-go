package backend

import "time"

// Function represents a server-side function
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

// RunHistory represents a function execution result
type RunHistory struct {
	ID        string    `json:"id"`
	Version   int       `json:"version"`
	Started   time.Time `json:"started"`
	Completed time.Time `json:"completed"`
	Success   bool      `json:"success"`
	Output    []string  `json:"output"`
}

// AddFunction adds a server-side function to an account
func AddFunction(token string, fn Function) error {
	return Post(token, "/fn/add", fn, nil)
}

// ListFunctions lists all functions for this account
func ListFunctions(token string) (results []Function, err error) {
	err = Get(token, "/fn", &results)
	return
}

// UpdateFunction updates a function
func UpdateFunction(token string, fn Function) error {
	return Post(token, "/fn/update", fn, nil)
}

// DeleteFunction removes a function
func DeleteFunction(token, name string) error {
	return Get(token, "/fn/del/"+name, nil)
}

// FunctionInfo gets info on a function, including execution histories
func FunctionInfo(token, name string) (fn Function, err error) {
	err = Get(token, "/fn/info/"+name, &fn)
	return
}
