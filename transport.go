package onet

import (
	"context"
	"fmt"
)

// Transport .
type Transport interface {
	fmt.Stringer
	Protocol() string
}

// NativeTransport .
type NativeTransport interface {
	Transport
	Listen(laddr *Addr, config *Config) (Listener, error)
	Dial(ctx context.Context, raddr *Addr, config *Config) (Conn, error)
}

// OverlayTransport .
type OverlayTransport interface {
	Transport
	Client(conn Conn, raddr *Addr, config *Config) (Conn, error)
	Server(conn Conn, laddr *Addr, config *Config) (Conn, error)
}

// MuxTransport .
type MuxTransport interface {
	OpenStream(ctx context.Context, raddr *Addr, config *Config) (Conn, error)
	AcceptStream() (Conn, error)
}

var transports = make(map[string]Transport)
