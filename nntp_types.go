package nntp

import (
	"fmt"
	"time"
)

type GroupStat struct {
	Count int
	First int
	Last  int
	Group string
}

type GroupListItem struct {
	Group      string
	Last       int
	First      int
	Permission GroupPermission
}

type GroupDescriptionListItem struct {
	Group       string
	Description string
}

type Range struct {
	First int
	Last  int
}

func (r Range) String() string {
	first := r.First
	if first == 0 {
		first = 1
	}
	if r.Last == 0 {
		return fmt.Sprintf("%d-", r.First)
	}
	return fmt.Sprintf("%d-%d", r.First, r.Last)
}

type ArticleOverview struct {
	ArticleNumber int
	Subject       string
	From          string
	Date          Timestamp
	MessageID     MessageID
	References    string
	Bytes         uint64
	Lines         uint64
	ExtraFields   []string
}

func FullMessageID(messageID string) string {
	if messageID[0] != '<' {
		messageID = "<" + messageID
	}
	if messageID[len(messageID)-1] != '>' {
		messageID += ">"
	}
	return messageID
}

type OverviewFieldFormat struct {
	Name string
	Type OverviewFieldType
}

type MessageID string

// Returns the message-id wrapped inside a "<>" pair if not already has.
func (id MessageID) Full() MessageID {
	if len(id) == 0 {
		return id // invalid id already
	}
	if id[0] != '<' {
		id = "<" + id
	}
	if id[len(id)-1] != '>' {
		id += ">"
	}
	return id
}

func (id MessageID) Short() MessageID {
	if id.IsFull() {
		return id[1 : len(id)-1]
	}
	return id
}

func (id MessageID) IsFull() bool {
	return len(id) > 0 && id[0] == '<' && id[len(id)-1] == '>'
}

// ValidateFull valiates the message-id conforms to rfc3977#section-3.6 format:
//
// 1. A message-id MUST begin with "<", end with ">", and MUST NOT contain the latter except at the end.
//
// 2. A message-id MUST be between 3 and 250 octets in length.
//
// 3. A message-id MUST NOT contain octets other than printable US-ASCII characters.
func (id MessageID) ValidateFull() (err error) {
	l := len(id)
	if l < 3 || l > 250 {
		err = fmt.Errorf("full message-id MUST be between 3 and 250 octets in length, instead it has %d bytes: %w", l, err)
		return
	}
	c := id[0]
	if c != '<' {
		err = fmt.Errorf("full message-id MUST begin with '<', instead it has %#x: %w", c, err)
		return
	}
	c = id[l-1]
	if c != '>' {
		err = fmt.Errorf("full message-id MUST end with '>', instead it has %#x: %w", c, err)
		return
	}
	for i := 1; i < l-1; i++ {
		c = id[i]
		if c < ' ' || c > '~' {
			err = fmt.Errorf("message-id MUST NOT contain octets other than printable US-ASCII characters, instead it contains %#x: %w", c, err)
			return
		}
	}
	return
}

// Same as ValidateFull with the relaxation that the message-id doesn't have to begin with "<" and end with ">".
func (id MessageID) Validate() (err error) {
	if id.IsFull() {
		return id.ValidateFull()
	}
	l := len(id)
	var c byte
	for i := 0; i < l; i++ {
		c = id[i]
		if c < ' ' || c > '~' {
			err = fmt.Errorf("message-id MUST NOT contain octets other than printable US-ASCII characters, instead it contains %#x: %w", c, err)
			return
		}
	}
	return
}

type Timestamp string

func (ts *Timestamp) Time() (t time.Time, err error) {
	if t, err = time.Parse(DefaultArticleDateLayout, string(*ts)); err == nil {
		return
	}
	return time.Parse(AlternativeArticleDateLayout, string(*ts))
}
