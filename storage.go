package backend

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// StoreFileResult incluses the file id and url. The ID is required
// when deleting file
type StoreFileResult struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// File describes a stored file entry returned by the storage API.
type File struct {
	ID        string    `json:"id"`
	AccountID string    `json:"accountId"`
	Key       string    `json:"key"`
	URL       string    `json:"url"`
	Size      int64     `json:"size"`
	Uploaded  time.Time `json:"uploaded"`
}

// FileUsage describes the current account storage usage.
type FileUsage struct {
	Bytes int64   `json:"bytes"`
	GB    float64 `json:"gb"`
}

// FileListResult contains paged file results from the storage API.
type FileListResult struct {
	Page    int64  `json:"page"`
	Size    int64  `json:"size"`
	Total   int64  `json:"total"`
	Results []File `json:"results"`
}

// ListFilesParams configures paging and sorting for file listing.
//
// The current API supports page selection and an optional "size" sort.
type ListFilesParams struct {
	Page   int
	SortBy string
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

// StorageUsage returns the storage usage for the authenticated account.
func StorageUsage(token string) (usage FileUsage, err error) {
	err = Get(token, "/storage/usage", &usage)
	return
}

// ListFiles lists files for the authenticated account.
func ListFiles(token string, params *ListFilesParams) (result FileListResult, err error) {
	qs := url.Values{}
	if params != nil {
		if params.Page > 0 {
			qs.Add("page", strconv.Itoa(params.Page))
		}
		if params.SortBy != "" {
			qs.Add("sort", params.SortBy)
		}
	}

	u := "/storage/files"
	if enc := qs.Encode(); enc != "" {
		u = fmt.Sprintf("%s?%s", u, enc)
	}

	err = Get(token, u, &result)
	return
}

// DeleteFile deletes the file from storage and remove from space used for
// this account
func DeleteFile(token, id string) (ok bool, err error) {
	err = Get(token, fmt.Sprintf("/sudostorage/delete?id=%s", id), &ok)
	return
}
