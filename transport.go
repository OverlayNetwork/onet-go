package onet

import (
	"context"
	"fmt"

	"github.com/libs4go/errors"
)

// Transport .
type Transport interface {
	fmt.Stringer
	Protocol() string
}

// NativeTransport .
type NativeTransport interface {
	Transport
	Listen(onet *OverlayNetwork) (Listener, error)
	Dial(ctx context.Context, onet *OverlayNetwork) (Conn, error)
}

// OverlayTransport .
type OverlayTransport interface {
	Transport
	Client(onet *OverlayNetwork, conn Conn, chainOffset int) (Conn, error)
	Server(onet *OverlayNetwork, conn Conn, chainOffset int) (Conn, error)
}

// MuxTransport .
type MuxTransport interface {
	Transport
	Listen(onet *OverlayNetwork, chainOffset int) (Listener, error)
	Dial(ctx context.Context, onet *OverlayNetwork, chainOffset int) (Conn, error)
	OverlayTransport
}

var transports = make(map[string]Transport)

// RegisterTransport .
func RegisterTransport(transport Transport) error {
	if _, ok := transports[transport.Protocol()]; ok {
		return errors.Wrap(ErrExists, "transport for protocol %s already register", transport.Protocol())
	}

	transports[transport.Protocol()] = transport

	return nil
}

// RegisterTransports .
func RegisterTransports(transports ...Transport) error {
	for _, transport := range transports {
		if err := RegisterTransport(transport); err != nil {
			return err
		}
	}

	return nil
}
