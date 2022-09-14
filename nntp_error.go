package nntp

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func (err Error) Error() string {
	return fmt.Sprintf("NNTP Response Code: %d, %s", err.Code, err.Message)
}

var ErrorInvalidParams = errors.New("invalid parameters")
var ErrorInvalidMessageID = errors.New("invalid message-id format")
var ErrorParsingResponse = errors.New("cannot parse response")
