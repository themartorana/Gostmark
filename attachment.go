package gostmark

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"
)

type Attachment struct {
	Name        string
	ContentType string
	Reader      io.Reader
	ContentID   string

	contents string

	sync.Mutex
}

// Return an attachment
func NewAttachment(name, contentType string, reader io.Reader) *Attachment {
	return &Attachment{
		Name:        name,
		ContentType: contentType,
		Reader:      reader,
	}
}

// Contents returns the file contents. Currently
// non-streaming, and memory-caching, so hardly
// efficient.
func (a *Attachment) Contents() (string, error) {
	a.Mutex.Lock()
	if len(a.contents) == 0 {
		b, e := ioutil.ReadAll(a.Reader)
		if e != nil {
			return "", e
		} else {
			b64 := base64.StdEncoding.EncodeToString(b)
			a.contents = b64
		}
	}
	a.Mutex.Unlock()
	return a.contents, nil
}

// MarshalJSON exports the attachment as JSON
// for sending to the server
func (a *Attachment) MarshalJSON() ([]byte, error) {
	fileContents, err := a.Contents()
	if err != nil {
		return []byte{}, err
	}

	packet := struct {
		Name        string
		ContentType string
		Content     string
		ContentID   *string `json:",omitempty"`
	}{
		Name:        a.Name,
		ContentType: a.ContentType,
		Content:     string(fileContents),
	}

	// Content ID - omit if empty
	if a.ContentID != "" {
		// Copy
		cid := a.ContentID
		packet.ContentID = &cid
	}

	return json.Marshal(&packet)
}
