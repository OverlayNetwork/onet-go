package onet

import (
	"context"
	"io"
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
