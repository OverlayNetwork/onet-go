package onet

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/libs4go/errors"
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

// Dial dial to the remote overlay address
func Dial(ctx context.Context, raddr *Addr, options ...Option) (Conn, error) {
	network, err := ParseOverlayNetwork(raddr, options...)

	if err != nil {
		return nil, err
	}

	return network.Dial(ctx)
}

// Dial dial to the remote overlay address with config
func (network *OverlayNetwork) Dial(ctx context.Context) (Conn, error) {

	var conn Conn

	var err error

	var muxTransport MuxTransport

	for i, transport := range network.MuxTransports {
		conn, err = transport.Dial(ctx, network, i)

		if err != nil {
			if errors.Unwrap(err) != ErrMuxNotFound {
				return nil, errors.Wrap(err, "call mux %s Dial error", transport)
			}

			continue
		}

		muxTransport = transport

		break
	}

	var overlayTransports []OverlayTransport

	if muxTransport == nil {
		conn, err = network.NativeTransport.Dial(ctx, network)

		if err != nil {
			return nil, errors.Wrap(err, "call transport %s Dial error", network.NativeTransport)
		}

		overlayTransports = network.OverlayTransports

	} else {
		for i, t := range network.OverlayTransports {
			if t == muxTransport {
				overlayTransports = network.OverlayTransports[i+1:]
				break
			}
		}
	}

	for i, t := range overlayTransports {

		conn, err = t.Client(network, conn, i)

		if err != nil {
			return nil, errors.Wrap(err, "call transport %s Dial error", network.NativeTransport)
		}
	}

	return conn, nil
}

type netConnWrapper struct {
	net.Conn
	onet  *OverlayNetwork
	laddr *Addr
	raddr *Addr
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

// ToOnetConn .
func ToOnetConn(conn net.Conn, onet *OverlayNetwork) (Conn, error) {

	laddr, err := FromNetAddr(conn.LocalAddr())

	if err != nil {
		return nil, err
	}

	raddr, err := FromNetAddr(conn.RemoteAddr())

	if err != nil {
		return nil, err
	}

	return &netConnWrapper{
		Conn:  conn,
		onet:  onet,
		laddr: laddr,
		raddr: raddr,
	}, nil
}

// ToOnetConnWithAddr .
func ToOnetConnWithAddr(conn net.Conn, onet *OverlayNetwork, laddr, raddr *Addr) (Conn, error) {

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

	// wrapper, ok := conn.(*netConnWrapper)

	// if ok {
	// 	return wrapper.Conn, nil
	// }

	laddr, err := conn.LocalAddr().ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	raddr, err := conn.RemoteAddr().ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	return &onetConnWrapper{
		Conn:  conn,
		laddr: laddr,
		raddr: raddr,
	}, nil
}
