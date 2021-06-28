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
	Client(ctx context.Context, onet *OverlayNetwork, addr *Addr, next Next) (Conn, error)
	Server(ctx context.Context, onet *OverlayNetwork, addr *Addr, next Next) (Conn, error)
}

// Next .
type Next func() (Conn, error)

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
