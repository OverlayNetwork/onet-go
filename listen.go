package onet

import (
	"context"
	"net"
)

// Listener overlay network listener equal to net.Listener
type Listener interface {
	Accept() (Conn, error)

	Close() error

	Addr() *Addr
}

type listenerImpl struct {
	network *OverlayNetwork
}

func acceptNextTransport(ctx context.Context, network *OverlayNetwork, i int) (Conn, error) {
	if i < len(network.Transports) {
		return network.Transports[i].Server(ctx, network, network.Addrs[i], func() (Conn, error) {
			return acceptNextTransport(ctx, network, i+1)
		})
	}

	return nil, nil
}

func closeNextTransport(network *OverlayNetwork, i int) error {
	if i < len(network.Transports) {
		return network.Transports[i].Close(network, network.Addrs[i], func() error {
			return closeNextTransport(network, i+1)
		})
	}

	return nil
}

func (l *listenerImpl) Accept() (Conn, error) {
	return acceptNextTransport(context.Background(), l.network, 0)
}

func (l *listenerImpl) Close() error {
	return closeNextTransport(l.network, 0)
}

func (l *listenerImpl) Addr() *Addr {
	return l.network.Addr
}

// Listen listen on the local overlay address
func Listen(addr *Addr, options ...Option) (Listener, error) {
	network, err := ParseOverlayNetwork(addr, options...)

	if err != nil {
		return nil, err
	}

	return &listenerImpl{network: network}, nil
}

type onetListenerWrapper struct {
	listener Listener
	laddr    net.Addr
}

func (listener *onetListenerWrapper) Accept() (net.Conn, error) {

	conn, err := listener.listener.Accept()

	if err != nil {
		return nil, err
	}

	return FromOnetConn(conn)

}

func (listener *onetListenerWrapper) Close() error {
	return listener.listener.Close()
}

func (listener *onetListenerWrapper) Addr() net.Addr {
	return listener.laddr
}

// FromOnetListener .
func FromOnetListener(listener Listener) (net.Listener, error) {
	laddr, _, err := listener.Addr().ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	return &onetListenerWrapper{
		listener: listener,
		laddr:    laddr,
	}, nil
}
