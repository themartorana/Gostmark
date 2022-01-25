package gostmark

import (
	"encoding/json"
	"errors"

	"fmt"

	"github.com/themartorana/Gostmark/v1/raw"
)

type Server struct {
	ID   int
	Name string

	ApiTokens  []string
	ServerLink string
	Color      string

	SmtpApiActivated bool
	RawEmailEnabled  bool

	InboundAddress       string
	InboundHookUrl       string
	InboundDomain        string
	InboundHash          string
	InboundSpamThreshold int

	BounceHookUrl string
	OpenHookUrl   string

	PostFirstOpenOnly bool
	TrackOpens        bool

	client Client
}

func (s Server) Save() (Server, error) {
	if s.client.AccountToken == "" {
		return s, errors.New("accountToken not provided. Please create new servers using Client.NewServer()")
	}

	var body string
	var err error
	if s.ID != 0 {
		body, err = s.saveEdit()
	} else {
		body, err = s.saveNew()
	}

	if err != nil {
		return s, err
	}

	// Parse the response back into a new server object,
	// obstensibly with the things we submitted
	sNew := Server{
		client: s.client,
	}
	err = json.Unmarshal([]byte(body), &sNew)
	return sNew, err
}

// Delete deletes the server from Postmark
// TODO: Implement
func (s Server) Delete() error {
	return errors.New("NOT YET IMPLEMENTED")
}

func (s Server) saveEdit() (string, error) {
	savePacket, err := s.savePacket()
	if err != nil {
		return "", err
	}
	return raw.ResponseFromPostmarkPost(
		s.client.Host,
		fmt.Sprintf(
			"/servers/%s",
			s.ID,
		),
		map[string]string{
			"X-Postmark-Account-Token": s.client.AccountToken,
		},
		savePacket,
	)
}

func (s Server) saveNew() (string, error) {
	savePacket, err := s.savePacket()
	if err != nil {
		return "", err
	}
	return raw.ResponseFromPostmarkPost(
		s.client.Host,
		"/servers",
		map[string]string{
			"X-Postmark-Account-Token": s.client.AccountToken,
		},
		savePacket,
	)
}

// savePacket does error checking and creates an
// appropriate map for sending to the server
func (s Server) savePacket() (map[string]interface{}, error) {
	packet := map[string]interface{}{
		"SmtpApiActivated":     s.SmtpApiActivated,
		"RawEmailEnabled":      s.RawEmailEnabled,
		"PostFirstOpenOnly":    s.PostFirstOpenOnly,
		"TrackOpens":           s.TrackOpens,
		"InboundSpamThreshold": s.InboundSpamThreshold,
	}

	// Name required
	if s.Name == "" {
		return packet, errors.New("Server name required")
	} else {
		packet["Name"] = s.Name
	}

	// Everything else
	if s.Color != "" {
		packet["Color"] = s.Color
	}
	if s.InboundHookUrl != "" {
		packet["InboundHookUrl"] = s.InboundHookUrl
	}
	if s.BounceHookUrl != "" {
		packet["BounceHookUrl"] = s.BounceHookUrl
	}
	if s.OpenHookUrl != "" {
		packet["OpenHookUrl"] = s.OpenHookUrl
	}
	if s.InboundDomain != "" {
		packet["InboundDomain"] = s.InboundDomain
	}

	return packet, nil
}
