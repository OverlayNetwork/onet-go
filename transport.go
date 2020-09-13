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
	Listen(onet *OverlayNetwork, laddr *Addr, config *Config) (Listener, error)
	Dial(ctx context.Context, onet *OverlayNetwork, raddr *Addr, config *Config) (Conn, error)
}

// OverlayTransport .
type OverlayTransport interface {
	Transport
	Client(onet *OverlayNetwork, conn Conn, raddr *Addr, config *Config) (Conn, error)
	Server(onet *OverlayNetwork, conn Conn, laddr *Addr, config *Config) (Conn, error)
}

// MuxTransport .
type MuxTransport interface {
	NativeTransport
	OverlayTransport
}

var transports = make(map[string]Transport)
