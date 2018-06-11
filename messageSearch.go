package gostmark

import (
	"net/url"
	"strconv"
	"time"
)

type MessageSearchPacket struct {
	Recipient EmailAddress
	FromEmail EmailAddress
	Tag       string
	Status    MessageStatus
	ToDate    time.Time
	FromDate  time.Time
	Subject   string

	Count  int
	Offset int
}

type SearchResults struct {
	TotalCount int            `json:"TotalCount"`
	Messages   []SearchResult `json:"Messages"`
}

type SearchResult struct {
	Tag        string         `json:"Tag"`
	MessageID  string         `json:"MessageID"`
	To         []EmailAddress `json:"To"`
	Cc         []EmailAddress `json:"Cc"`
	Bcc        []EmailAddress `json:"Bcc"`
	Recipients []string       `json:"Recipients"`
	ReceivedAt time.Time      `json:"ReceivedAt"`
	From       string         `json:"From"`
	Subject    string         `json:"Subject"`
	Status     string         `json:"Status"`
	TrackOpens bool           `json:"TrackOpens"`
	TrackLinks string         `json:"TrackLinks"`
}

// MessageSearchPacket returns the search packet as url.Values
func (msp MessageSearchPacket) AsValues() url.Values {
	vals := make(url.Values)
	vals.Add("count", strconv.Itoa(msp.Count))
	vals.Add("offset", strconv.Itoa(msp.Count))

	if str, err := msp.Recipient.String(); err == nil {
		vals.Add("recipient", str)
	}
	if str, err := msp.FromEmail.String(); err == nil {
		vals.Add("fromemail", str)
	}
	if msp.Tag != "" {
		vals.Add("tag", msp.Tag)
	}
	if msp.Status != "" {
		vals.Add("status", string(msp.Status))
	}
	if !msp.FromDate.IsZero() {
		vals.Add("fromdate", msp.FromDate.Format(time.RFC3339))
	}
	if !msp.ToDate.IsZero() {
		vals.Add("todate", msp.ToDate.Format(time.RFC3339))
	}
	if msp.Subject != "" {
		vals.Add("subject", msp.Subject)
	}

	return vals
}
