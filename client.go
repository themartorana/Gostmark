package gostmark

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/themartorana/Gostmark/internal"
)

type Client struct {
	Host string

	AccountToken string
	ServerToken  string
}

const defaultHost string = "https://api.postmarkapp.com"

// ClientForAPIKey returns a new client intialized
// to the provided API key
func ClientForAccountToken(accountToken string) Client {
	return Client{
		AccountToken: accountToken,
		Host:         defaultHost,
	}
}

// ClientForServerToken returns a new client intialized
// to the provided Server API key
func ClientForServerToken(serverToken string) Client {
	return Client{
		ServerToken: serverToken,
		Host:        defaultHost,
	}
}

// servers is an internal container for
// unmarshalling the Postmark server response
type servers struct {
	TotalCount int
	Servers    []Server
}

func (c Client) NewServer() Server {
	return Server{
		client: c,
	}
}

// GetServerForToken retreives a server struct for
// the server token supplied
func (c Client) GetServerByToken(serverToken string) (Server, error) {
	body, err := internal.GetRawResponseFromPostmark(
		c.Host,
		"/server",
		map[string]string{
			"X-Postmark-Server-Token": serverToken,
		},
		nil,
	)
	if err != nil {
		return Server{}, err
	}

	var s Server
	err = json.Unmarshal([]byte(body), &s)
	if err == nil {
		s.client = c
	}
	return s, err
}

// GetServerForToken retreives a server struct for
// the server token supplied
func (c Client) GetServerByID(serverID string) (Server, error) {
	body, err := internal.GetRawResponseFromPostmark(
		c.Host,
		fmt.Sprintf(
			"/servers/%s",
			serverID,
		),
		map[string]string{
			"X-Postmark-Account-Token": c.AccountToken,
		},
		nil,
	)
	if err != nil {
		return Server{}, err
	}

	var s Server
	err = json.Unmarshal([]byte(body), &s)
	if err == nil {
		s.client = c
	}
	return s, err
}

func (c Client) GetAllServers() ([]Server, error) {
	body, err := internal.GetRawResponseFromPostmark(
		c.Host,
		"/servers?",
		map[string]string{
			"X-Postmark-Account-Token": c.AccountToken,
		},
		nil,
	)
	if err != nil {
		return []Server{}, err
	}

	var serversResponse servers
	err = json.Unmarshal([]byte(body), &serversResponse)
	if err != nil {
		return serversResponse.Servers, err
	}

	// Associate the account token
	returnServers := make([]Server, 0, len(serversResponse.Servers))
	for _, server := range serversResponse.Servers {
		server.client = c
		returnServers = append(returnServers, server)
	}

	return returnServers, nil
}

// SendMessage sends a single message through Postmark
func (c Client) SendMessage(message Message) (MessageSendResponse, error) {
	if c.ServerToken == "" {
		return MessageSendResponse{}, errors.New("ServerToken must be set in Client")
	}

	// Get the internal JSON
	bytes, err := json.Marshal(message)
	if err != nil {
		return MessageSendResponse{}, err
	}

	// Check the size against the 10MB limit
	size := len(bytes)
	if size > 1024*10 {
		return MessageSendResponse{}, errors.New(
			fmt.Sprintf(
				"Message + attachments cannot excede 10MB. Current size: %d Bytes",
				size,
			),
		)
	}

	// Post and get the response
	body, err := internal.GetRawResponseFromPostmark(
		c.Host,
		"/email",
		map[string]string{
			"X-Postmark-Server-Token": c.ServerToken,
		},
		string(bytes),
	)
	if err != nil {
		return MessageSendResponse{}, err
	}

	// Unmarshal response
	var msr MessageSendResponse
	err = json.Unmarshal([]byte(body), &msr)
	return msr, err
}

// SendMessages batch-sends Messages
func (c Client) SendMessages(messages []Message) (BatchMessageSendResponse, error) {
	if len(messages) > 500 {
		return BatchMessageSendResponse{}, errors.New("Cannot send over 500 messages in a single batch")
	}

	// Get the internal JSON. The
	// array should marshal properly
	bytes, err := json.Marshal(messages)
	if err != nil {
		return BatchMessageSendResponse{}, err
	}

	// Post and get the response
	body, err := internal.GetRawResponseFromPostmark(
		c.Host,
		"/email/batch",
		map[string]string{
			"X-Postmark-Account-Token": c.AccountToken,
		},
		string(bytes),
	)
	if err != nil {
		return BatchMessageSendResponse{}, err
	}

	// Unmarshal the response
	var responses BatchMessageSendResponse
	err = json.Unmarshal([]byte(body), &responses)
	if err != nil {
		return BatchMessageSendResponse{}, err
	}

	return responses, nil
}
