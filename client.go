package postmark

import (
	"encoding/json"
	"fmt"

	"github.com/themartorana/Gostmark/internal"
)

type Client struct {
	AccountToken string
	Host         string
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

// servers is used for requesting
// multiple servers at the same time
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
