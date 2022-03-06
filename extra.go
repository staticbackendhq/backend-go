package backend

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

// ResizeImage upload and resize and image based on the max width allowed.
// The input image must be a PNG or JPG and the end result will always be a JPG
func ResizeImage(token, filename string, file io.ReadSeeker, maxWidth float64) (StoreFileResult, error) {
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

	ww, err := w.CreateFormField("width")
	if err != nil {
		return res, fmt.Errorf("error creating the width field: %v", err)
	}

	if _, err := ww.Write([]byte(fmt.Sprintf("%f", maxWidth))); err != nil {
		return res, fmt.Errorf("error writing the max width parameter: %v", err)
	}

	w.Close()

	if err := request(token, "POST", "/extra/resizeimg", w.FormDataContentType(), &buf, &res); err != nil {
		return res, fmt.Errorf("error while uploading file: %v", err)
	}

	return res, nil
}
