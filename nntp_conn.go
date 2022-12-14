package nntp

import (
	"fmt"
	"io"

	"gopkg.in/textproto.v0"
)

type Conn struct {
	*textproto.Conn
}

func (conn *Conn) Close() error {
	return conn.Conn.Close()
}

func NewConn(conn io.ReadWriteCloser) *Conn {
	return &Conn{textproto.NewConn(conn)}
}

func (conn *Conn) ReadWelcome() (err error) {
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.readWelcome] failed to read Welcome message: %w", err)
		return
	}
	switch ResponseCode(code) {
	case ResponseCodeReadyPostingAllowed, ResponseCodeReadyPostingProhibited: // 200 || 201
		err = nil
	default:
		err = fmt.Errorf("[nntp.readWelcome] unexpected response: %w", &Error{ResponseCode(code), msg})
	}
	return
}
