package postmark

import (
	"net/mail"
	"time"
)

type Message struct {
	From    mail.Address
	ReplyTo mail.Address

	To  []mail.Address
	Cc  []mail.Address
	Bcc []mail.Address

	Subject  string
	HtmlBody string
	TextBody string

	Headers []Header

	Tag        string
	TrackOpens bool

	Attachments []Attachment
}

type Response struct {
	To        string    `json:"To"`
	Submitted time.Time `json:"SubmittedAt"`
	ErrorCode int       `json:"ErrorCode"`

	MessageID string `json:"MessageID"`
	Message   string `json:"Message"`
}
