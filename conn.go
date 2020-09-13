package netty

import (
	"io"
	"time"
)

// Conn overlay network conn object equal to net.Conn
type Conn interface {
	io.ReadWriteCloser

	LocalAddr() Addr

	RemoteAddr() Addr

	SetDeadline(t time.Time) error

	SetReadDeadline(t time.Time) error

	SetWriteDeadline(t time.Time) error
}

// Dial dial to the remote overlay address
func Dial(addr *Addr, options ...Option) (Conn, error) {
	network, err := ParseOverlayNetwork(addr, options...)

	if err != nil {
		return nil, err
	}

	return network.Dial()
}
