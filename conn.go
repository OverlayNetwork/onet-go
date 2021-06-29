package onet

import (
	"context"
	"io"
	"net"
	"time"
)

// Conn overlay network conn object equal to net.Conn
type Conn interface {
	io.ReadWriteCloser

	LocalAddr() *Addr

	RemoteAddr() *Addr

	SetDeadline(t time.Time) error

	SetReadDeadline(t time.Time) error

	SetWriteDeadline(t time.Time) error

	ONet() *OverlayNetwork
}

func callClientNext(ctx context.Context, network *OverlayNetwork, i int) (Conn, error) {
	if i < len(network.Transports) {
		return network.Transports[i].Client(ctx, network, network.Addrs[i], func() (Conn, error) {
			return callClientNext(ctx, network, i+1)
		})
	}

	return nil, nil
}

// Dial dial to the remote overlay address
func Dial(ctx context.Context, raddr *Addr, options ...Option) (Conn, error) {
	network, err := ParseOverlayNetwork(raddr, options...)

	if err != nil {
		return nil, err
	}

	return callClientNext(ctx, network, 0)
}

type netConnWrapper struct {
	net.Conn
	onet    *OverlayNetwork
	laddr   *Addr
	raddr   *Addr
	closing func() error
}

func (conn *netConnWrapper) LocalAddr() *Addr {
	return conn.laddr
}

func (conn *netConnWrapper) RemoteAddr() *Addr {
	return conn.raddr
}

func (conn *netConnWrapper) ONet() *OverlayNetwork {
	return conn.onet
}

func (conn *netConnWrapper) Close() error {
	if err := conn.Conn.Close(); err != nil {
		return err
	}

	if conn.closing != nil {
		return conn.closing()
	}

	return nil
}

// ToOnetConn .
func ToOnetConn(conn net.Conn, onet *OverlayNetwork, addr *Addr) (Conn, error) {

	_, relativeAddr, err := addr.ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	laddr, err := FromNetAddr(conn.LocalAddr())

	if err != nil {
		return nil, err
	}

	raddr, err := FromNetAddr(conn.RemoteAddr())

	if err != nil {
		return nil, err
	}

	laddr = laddr.Join(relativeAddr.SubAddrs()...)

	raddr = raddr.Join(relativeAddr.SubAddrs()...)

	return &netConnWrapper{
		Conn:  conn,
		onet:  onet,
		laddr: laddr,
		raddr: raddr,
	}, nil
}

type onetConnWrapper struct {
	Conn
	laddr net.Addr
	raddr net.Addr
}

func (conn *onetConnWrapper) LocalAddr() net.Addr {
	return conn.laddr
}

func (conn *onetConnWrapper) RemoteAddr() net.Addr {
	return conn.raddr
}

// FromOnetConn .
func FromOnetConn(conn Conn) (net.Conn, error) {

	laddr, _, err := conn.LocalAddr().ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	raddr, _, err := conn.RemoteAddr().ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	return &onetConnWrapper{
		Conn:  conn,
		laddr: laddr,
		raddr: raddr,
	}, nil
}
