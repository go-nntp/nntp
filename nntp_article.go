package nntp

import (
	"io"
	"net/textproto"
)

type Article struct {
	// The article number in the currently selected group on the NNTP server. The article number is set to 0 if no
	// current group is selected, or the article is used as an argument to posting commands.
	ArticleNumber int

	// The message-id string that globally identifies the article across all groups and all peered NNTP servers. Can be
	// empty string if the article is used as an argument to posting commands.
	MessageID string

	// The parsed MIME header. Only not nil if the article is returned from ARTICLE or HEAD commands, or used as an
	// argument to posting commands.
	Header textproto.MIMEHeader

	// Reader of the contents after the article's MIME header and the double CRLF separator. Contents are dot encodong
	// decoded, that is, leading double dots are unescaped to single dot, and the final ".\r\n" sequence is dropped.
	// Clients consuming the article should finish consuming the Body content before issuing any other NNTP command on
	// the connection object where it gets the article, since the connection manager will automatically drain the
	// article if another NNTP command is issued.
	//
	// Body is only not nil if the article is returned from ARTICLE or BODY commands, or used as an argument to posting
	// commands.
	Body io.Reader
}
