package nntp

import (
	"net/textproto"
)

type Conn struct {
	*textproto.Conn
}

func (conn *Conn) Close() error {
	return conn.Conn.Close()
}
