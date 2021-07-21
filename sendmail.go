package backend

// EmailData used to request the send email process
type EmailData struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	ReplyTo string `json:"replyTo"`
}

func SendMail(token, from, to, subject, body, replyTo string) (ok bool, err error) {
	data := EmailData{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
		ReplyTo: replyTo,
	}
	err = Post(token, "/sudo/sendmail", data, &ok)
	return
}
