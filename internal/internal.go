package internal

import (
	"errors"

	"fmt"

	"github.com/franela/goreq"
)

type errorInfo struct {
	ErrorCode int
	Message   string
}

func GetRawResponseFromPostmark(url string, headers map[string]string, body interface{}) (string, error) {
	req := goreq.Request{
		Uri:         url,
		Accept:      "application/json",
		ContentType: "application/json",
	}

	// Add headers
	for key, val := range headers {
		req.AddHeader(key, val)
	}

	// Body?
	if body != nil {
		req.Body = body
		req.Method = "POST"
	}

	// Send
	resp, err := req.Do()
	if err != nil {
		return "", err
	}

	// Return body
	respBody, err := resp.Body.ToString()
	if err != nil {
		return "", err
	}

	// Check response code
	switch resp.StatusCode {
	case 200:
		return respBody, nil
	case 401:
		return "", errors.New("Missing or incorrect API token in header")
	case 422:
		var errInfo errorInfo
		err = resp.Body.FromJsonTo(&errInfo)
		if err == nil {
			err = errors.New(
				fmt.Sprintf(
					"API error %d: %s",
					errInfo.ErrorCode,
					errInfo.Message,
				),
			)
		}
		return "", err
	case 500:
		return "", errors.New("Internal Server Error")
	case 503:
		return "", errors.New("Postmark Servers Temporarilty Unavailable")
	default:
		return "", errors.New(
			fmt.Sprintf(
				"Unrecognized error %d: %s",
				resp.StatusCode,
				respBody,
			),
		)
	}
}
