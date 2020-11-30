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

func request(token, method, url, ct string, body io.Reader, v interface{}) error {
	host := "http://localhost:8099"

	start := time.Now()

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", host, url), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", ct)

	// We provide authentication
	req.Header.Set("SB-PUBLIC-KEY", PublicKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if Verbose {
		fmt.Printf("%d\t%s\t%v\t%s bytes\n", res.StatusCode, url, time.Since(start), res.Header.Get("Content-Length"))
	}

	// Did we got an error
	if res.StatusCode > 299 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}

		return fmt.Errorf("error returned by the backend: %s", string(b))
	}

	if res.Header.Get("Content-Type") == "application/json" {
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

// Get sends an HTTP GET request to the backend
func Get(token, url string, v interface{}) error {
	return request(token, http.MethodGet, url, "application/json", nil, v)
}

func Post(token, url string, body interface{}, v interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error while encoding the body to JSON: %v", err)
	}

	buf := bytes.NewReader(b)
	return request(token, "POST", url, "application/json", buf, v)
}