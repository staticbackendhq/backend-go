// Package backend is a client wrapper for [StaticBackend] API.
//
// Before using its functions you need to supply a Region and PublicKey.
//
// [StaticBackend]: https://staticbackend.com
package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// PublicKey is required for all HTTP requests.
var (
	PublicKey string
	Verbose   bool
	Region    string
)

const (
	RegionNorthAmerica1 = "na1" // managed hosting region of North-America
	RegionLocalDev      = "dev" // local dev region default to http://localhost:8099
)

func getHost() (host string) {
	switch Region {
	case RegionNorthAmerica1:
		host = "https://na1.staticbackend.com"
	case RegionLocalDev, "":
		host = "http://localhost:8099"
	default:
		// for self-hosted instance
		host = Region
	}
	return
}

func request(token, method, url, ct string, body io.Reader, v interface{}) error {
	host := getHost()

	start := time.Now()

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", host, url), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", ct)

	// We provide authentication
	req.Header.Set("SB-PUBLIC-KEY", PublicKey)
	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if Verbose {
		fmt.Printf("%s\t%d\t%s\t%v\t%s bytes\n", method, res.StatusCode, url, time.Since(start), res.Header.Get("Content-Length"))
	}

	// Did we got an error
	if res.StatusCode > 299 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}

		return fmt.Errorf("error returned by the backend: %s", string(b))
	}

	if res.Header.Get("Content-Type") == "application/json" && v != nil && res.Body != nil {
		if err := json.NewDecoder(res.Body).Decode(v); err != nil {
			return fmt.Errorf("unable to decode the response body: %v", err)
		}

		return nil
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %v", err)
	}

	v = b
	return nil
}

// Get sends an HTTP GET request to the backend API
func Get(token, url string, v interface{}) error {
	return request(token, http.MethodGet, url, "application/json", nil, v)
}

// Post sends an HTTP POST request to the backend API
func Post(token, url string, body interface{}, v interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error while encoding the body to JSON: %v", err)
	}

	buf := bytes.NewReader(b)
	return request(token, "POST", url, "application/json", buf, v)
}

// Put sends an HTTP POST request to the backend API
func Put(token, url string, body interface{}, v interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error while encoding the body to JSON: %v", err)
	}

	buf := bytes.NewReader(b)
	return request(token, "PUT", url, "application/json", buf, v)
}

// Delete sends an HTTP DELETE request to the backend API
func Del(token, url string) error {
	var v interface{}
	return request(token, http.MethodDelete, url, "application/json", nil, v)
}
