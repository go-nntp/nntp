package nntp

type Client struct {
	*Conn
}

func New(conn *Conn) *Client {
	return &Client{conn}
}

/**
 * Authenticate.
 *
 * <b>Non-standard!</b><br>
 * This method uses non-standard commands, which is not part
 * of the original RFC977, but has been formalized in RFC2890.
 *
 * @param string	$user	The username
 * @param string	$pass	The password
 */
func (client *Client) Authenticate(user, pass string) (err error) {
	return client.CmdAuthinfo(user, pass)
}
