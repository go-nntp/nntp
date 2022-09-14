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
	Date          time.Time
	MessageID     string
	References    string
	Bytes         uint64
	Lines         uint64
	ExtraFields   []string
}

// Make sure the message-id is wrapped inside a "<>" pair.
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

// Validate the message-id conforms to rfc3977#section-3.6 format:
//
// 1. A message-id MUST begin with "<", end with ">", and MUST NOT contain the latter except at the end.
//
// 2. A message-id MUST be between 3 and 250 octets in length.
//
// 3. A message-id MUST NOT contain octets other than printable US-ASCII characters.
func ValidateMessageID(messageID string) (err error) {
	l := len(messageID)
	if l < 3 || l > 250 {
		err = fmt.Errorf("message-id MUST be between 3 and 250 octets in length, instead it has %d bytes: %w", l, err)
		return
	}
	c := messageID[0]
	if c != '<' {
		err = fmt.Errorf("message-id MUST begin with '<', instead it has %#x: %w", c, err)
		return
	}
	c = messageID[l-1]
	if c != '>' {
		err = fmt.Errorf("message-id MUST end with '>', instead it has %#x: %w", c, err)
		return
	}
	for i := 1; i < l-1; i++ {
		c = messageID[i]
		if c < ' ' || c > '~' {
			err = fmt.Errorf("message-id MUST NOT contain octets other than printable US-ASCII characters, instead it contains %#x: %w", c, err)
			return
		}
	}
	return
}
