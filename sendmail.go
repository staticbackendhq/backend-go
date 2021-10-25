package backend

// EmailData used to request the send email process
type EmailData struct {
	FromName string `json:"fromName"`
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	ReplyTo  string `json:"replyTo"`
}

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
