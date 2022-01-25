package gostmark

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// Address represents an email
// structured address string
type EmailAddress struct {
	Name  string
	Email string
}

// Returns an Address struct initialized with
// an email address only
func EmailAddressForEmail(email string) EmailAddress {
	return EmailAddress{
		Email: email,
	}
}

func (e EmailAddress) String() (string, error) {
	if e.Email == "" {
		return "", errors.New("email cannot be empty")
	}

	// Steal already considered functionality, why not?
	ma := &mail.Address{
		Name:    e.Name,
		Address: e.Email,
	}

	return ma.String(), nil
}

// MarshalJSON returns the Address as a proper email string.
func (e EmailAddress) MarshalJSON() ([]byte, error) {
	s, err := e.String()
	if err != nil {
		return []byte(""), err
	}
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = fmt.Sprintf("\"%s\"", s)
	return []byte(s), nil
}

// joinEmailAddresses is a convenience function to return a
// comma delimited list of email addresses from a []EmailAddress
func joinEmailAddresses(addresses []EmailAddress) (string, error) {
	addressStrings := make([]string, 0, len(addresses))
	for _, emailAddress := range addresses {
		str, err := emailAddress.String()
		if err != nil {
			return "", err
		}
		addressStrings = append(addressStrings, str)
	}

	return strings.Join(addressStrings, ","), nil
}
