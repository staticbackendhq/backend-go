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

type SMSData struct {
	AccountSID string `json:"accountSID"`
	AuthToken  string `json:"authToken"`
	ToNumber   string `json:"toNumber"`
	FromNumber string `json:"fromNumber"`
	Body       string `json:"body"`
}

// SudoSendSMS sends a text message via the Twilio API. You need a valid Twilio
// AccountSID, AuthToken and phone number.
func SudoSendSMS(token string, data SMSData) error {
	var status bool
	if err := Post(token, "/extra/sms", data, &status); err != nil {
		return err
	}
	return nil
}

// ConvertParam used for the ConvertURLToX request.
type ConvertParam struct {
	// ToPDF indicates if the output is a PDF, otherwise a PNG
	ToPDF bool `json:"toPDF"`
	// URL a publicly available URL
	URL string `json:"url"`
	// FullPage indicates to PNG to screenshot the entire page (still not working)
	FullPage bool `json:"fullpage"`
}

// ConvertURLToX converts a URL (web page) to either a PDF or a PNG. The ID and
// URL of the PDF or PNG is returned.
func ConvertURLToX(token string, data ConvertParam) (res StoreFileResult, err error) {
	err = Post(token, "/extra/htmltox", data, &res)
	return
}
