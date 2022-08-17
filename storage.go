package backend

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// StoreFileResult incluses the file id and url. The ID is required
// when deleting file
type StoreFileResult struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// StoreFile uploads a new file and returns its public URL using SB CDN.
func StoreFile(token, filename string, file io.ReadSeeker) (StoreFileResult, error) {
	var res StoreFileResult

	// multipart form data
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return res, fmt.Errorf("error creating form field: %v", err)
	}

	if _, err := io.Copy(fw, file); err != nil {
		return res, fmt.Errorf("error copying file data to form field: %v", err)
	}

	w.Close()

	if err := request(token, "POST", "/storage/upload", w.FormDataContentType(), &buf, &res); err != nil {
		return res, fmt.Errorf("error while uploading file: %v", err)
	}

	return res, nil
}

// DownloadFile retrieves the file content as []byte
func DownloadFile(token, fileURL string) (buf []byte, err error) {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf, err = io.ReadAll(resp.Body)
	return
}

// DeleteFile deletes the file from storage and remove from space used for
// this account
func DeleteFile(token, id string) (ok bool, err error) {
	err = Get(token, fmt.Sprintf("/sudostorage/delete?id=%s", id), &ok)
	return
}
