package gostmark

import (
	"encoding/json"
	"errors"

	"sync"
	"time"
)

type Message struct {
	From    EmailAddress
	ReplyTo EmailAddress
	To      EmailAddress

	Cc  []EmailAddress
	Bcc []EmailAddress

	Subject  string
	HtmlBody string
	TextBody string

	Headers []Header

	Tag        string
	TrackOpens bool

	Attachments []*Attachment

	// Template stuff, to incorporate eventually
	// TemplateId    int
	// TemplateModel interface{}
	// InlineCSS 	 bool

	sync.Mutex
}

type MessageSendResponse struct {
	To          string
	SubmittedAt time.Time
	ErrorCode   int

	MessageID string
	Message   string
}

type BatchMessageSendResponse struct {
	TotalCount int
	Messages   []MessageSendResponse
}

func (m *Message) AddAttachment(attachment *Attachment) {
	m.Mutex.Lock()
	m.Attachments = append(m.Attachments, attachment)
	m.Mutex.Unlock()
}

func (m *Message) AddHeader(header Header) {
	m.Mutex.Lock()
	m.Headers = append(m.Headers, header)
	m.Mutex.Unlock()
}

func (m *Message) MarshalJSON() ([]byte, error) {
	if m.To.Email == "" {
		return []byte{}, errors.New("To EmailAddress required")
	}
	if m.From.Email == "" {
		return []byte{}, errors.New("From EmailAddress required")
	}
	if len(m.Cc) > 50 {
		return []byte{}, errors.New("Cc field cannot contain more than 50 entries")
	}
	if len(m.Bcc) > 50 {
		return []byte{}, errors.New("Bcc field cannot contain more than 50 entries")
	}
	if m.HtmlBody == "" && m.TextBody == "" {
		return []byte{}, errors.New("HtmlBody and TextBody cannot both be blank")
	}

	// Render
	packet, err := m.toSendPacket()
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(packet)
}

func (m *Message) toSendPacket() (map[string]interface{}, error) {
	// Slow, but flexible
	packet := map[string]interface{}{
		"From": m.From,
		"To":   m.To,
	}

	// Optional fields
	if len(m.Cc) > 0 {
		str, err := m.ccAsString()
		if err != nil {
			return packet, err
		}
		packet["Cc"] = str
	}
	if len(m.Bcc) > 0 {
		str, err := m.bccAsString()
		if err != nil {
			return packet, err
		}
		packet["Bcc"] = str
	}
	if m.Subject != "" {
		packet["Subject"] = m.Subject
	}
	if m.Tag != "" {
		packet["Tag"] = m.Tag
	}
	if m.HtmlBody != "" {
		packet["HtmlBody"] = m.HtmlBody
	}
	if m.TextBody != "" {
		packet["TextBody"] = m.TextBody
	}
	if m.ReplyTo.Email != "" {
		packet["ReplyTo"] = m.ReplyTo
	}
	if len(m.Headers) != 0 {
		packet["Headers"] = m.Headers
	}
	if m.TrackOpens {
		packet["TrackOpens"] = true
	}

	// Attachments marshal themselves
	if len(m.Attachments) > 0 {
		packet["Attachments"] = m.Attachments
	}

	return packet, nil
}

func (m *Message) ccAsString() (string, error) {
	return joinEmailAddresses(m.Cc)
}

func (m *Message) bccAsString() (string, error) {
	return joinEmailAddresses(m.Bcc)
}
