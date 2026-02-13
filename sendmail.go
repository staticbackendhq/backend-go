package backend

// Attachment represents an email attachment
type Attachment struct {
	URL         string `json:"url"`         // URL where the attachment was fetched from
	Body        []byte `json:"body"`        // Raw bytes of the attachment
	ContentType string `json:"contentType"` // MIME type of the attachment
	Filename    string `json:"filename"`    // Name of the file
}

// EmailData used to request the send email process
type EmailData struct {
	FromName string `json:"fromName"`
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	ReplyTo  string `json:"replyTo"`

	Attachments []Attachment `json:"attachments"`
}

// SendMail sends an email
func SendMail(token, from, fromName, to, subject, body, replyTo string) (ok bool, err error) {
	data := EmailData{
		FromName: fromName,
		From:     from,
		To:       to,
		Subject:  subject,
		Body:     body,
		ReplyTo:  replyTo,
	}
	err = Post(token, "/sudo/sendmail", data, &ok)
	return
}

// SendMailWithAttachments sends an email with attachments. The attachments
// either are downloaded via passing the URL or directly filling the Body,
// ContentType, and Filename. If using only URL, ensure the URL are publicly
// available.
func SendMailWithAttachments(token string, email EmailData) (ok bool, err error) {
	err = Post(token, "/sudo/sendmail", email, &ok)
	return
}
