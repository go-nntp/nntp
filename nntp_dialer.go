package nntp

import (
	"context"
	"crypto/tls"
	"net"
)

type Dialer struct {
	// NetDialer is the optional dialer to use for the underlying TCP connections. A nil NetDialer is equivalent to the
	// net.Dialer zero value.
	NetDialer *net.Dialer

	// Config is the TLS configuration to use for new TLS connections. A nil configuration is equivalent to the zero
	// configuration; see the documentation of tls.Config for the defaults.
	Config *tls.Config
}

func (d *Dialer) Dial(ctx context.Context, network, addr string) (conn *Conn, err error) {
	if d.NetDialer == nil {
		d.NetDialer = &net.Dialer{}
	}
	netconn, err := d.NetDialer.DialContext(ctx, network, addr)
	if err != nil {
		return
	}
	c := NewConn(netconn)
	if err = c.ReadWelcome(); err != nil {
		c.Close()
	}
	conn = c
	return
}

func (d *Dialer) DialTLS(ctx context.Context, network, addr string) (conn *Conn, err error) {
	dialer := tls.Dialer{
		NetDialer: d.NetDialer,
		Config:    d.Config,
	}
	netconn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return
	}
	c := NewConn(netconn)
	if err = c.ReadWelcome(); err != nil {
		c.Close()
	}
	conn = c
	return
}
