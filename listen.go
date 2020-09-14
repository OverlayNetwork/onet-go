package onet

import (
	"net"

	"github.com/libs4go/errors"
)

// Listener overlay network listener equal to net.Listener
type Listener interface {
	Accept() (Conn, error)

	Close() error

	// Addr returns the listener's network address.
	Addr() *Addr
}

// Listen listen on the local overlay address
func Listen(addr *Addr, options ...Option) (Listener, error) {
	network, err := ParseOverlayNetwork(addr, options...)

	if err != nil {
		return nil, err
	}

	return network.Listen()
}

// Listen listen on the local overlay address with config
func (network *OverlayNetwork) Listen() (Listener, error) {
	return newOverlayNetworkListener(network)
}

type acceptConn struct {
	conn Conn
	err  error
}

type overlayNetworkListener struct {
	network        *OverlayNetwork
	acceptConns    chan *acceptConn
	nativeListener Listener
	muxListeners   []Listener
}

func newOverlayNetworkListener(network *OverlayNetwork) (Listener, error) {
	nativeListener, err := network.NativeTransport.Listen(network)

	if err != nil {
		return nil, errors.Wrap(err, "call %s transport listen error", network.NativeTransport)
	}

	onetListener := &overlayNetworkListener{
		network:        network,
		nativeListener: nativeListener,
		acceptConns:    make(chan *acceptConn),
	}

	for i, mux := range network.MuxTransports {
		muxlistener, err := mux.Listen(network, i)

		if err != nil {
			onetListener.Close()
			return nil, errors.Wrap(err, "call %s transport listen error", mux)
		}

		onetListener.muxListeners = append(onetListener.muxListeners, muxlistener)
	}

	onetListener.startLoop()

	return onetListener, nil
}

func (listener *overlayNetworkListener) startLoop() {

	go listener.doAccept(listener.nativeListener)

	for _, l := range listener.muxListeners {
		go listener.doAccept(l)
	}
}

func (listener *overlayNetworkListener) doAccept(l Listener) {

	defer recover()

	for {
		conn, err := l.Accept()

		if err != nil {
			if errors.Unwrap(err) == ErrClosed {
				return
			}

			listener.acceptConns <- &acceptConn{conn: conn}

			continue
		}

		listener.wrapConn(l, conn)
	}
}

func (listener *overlayNetworkListener) wrapConn(l Listener, conn Conn) {

	var overlayTransports []OverlayTransport

	if listener.nativeListener == l {
		overlayTransports = listener.network.OverlayTransports
	} else {
		for i, tl := range listener.muxListeners {
			if tl == l {
				muxTransport := listener.network.MuxTransports[i]

				for j, transport := range listener.network.OverlayTransports {
					if transport == muxTransport {
						overlayTransports = listener.network.OverlayTransports[j+1:]

						break
					}
				}

				break
			}
		}
	}

	var err error

	for i, transport := range overlayTransports {
		conn, err = transport.Server(listener.network, conn, i)

		if err != nil {
			listener.acceptConns <- &acceptConn{
				err: errors.Wrap(err, "call transport %s Server error", transport),
			}

			return
		}
	}

	listener.acceptConns <- &acceptConn{conn: conn}
}

func (listener *overlayNetworkListener) Accept() (Conn, error) {
	conn, ok := <-listener.acceptConns

	if !ok {
		return nil, errors.Wrap(ErrClosed, "listener %s closed", listener.network.Addr)
	}

	return conn.conn, conn.err
}

func (listener *overlayNetworkListener) Close() error {

	listener.nativeListener.Close()

	for _, muxListener := range listener.muxListeners {
		muxListener.Close()
	}

	close(listener.acceptConns)

	return nil
}

func (listener *overlayNetworkListener) Addr() *Addr {
	return listener.network.Addr
}

type netListenerWrapper struct {
	listener net.Listener
	laddr    *Addr
	onet     *OverlayNetwork
}

func (listener *netListenerWrapper) Accept() (Conn, error) {

	conn, err := listener.listener.Accept()

	if err != nil {
		return nil, err
	}

	return ToOnetConn(conn, listener.onet)
}
func (listener *netListenerWrapper) Close() error {
	return listener.listener.Close()
}
func (listener *netListenerWrapper) Addr() *Addr {
	return listener.laddr
}

// ToOnetListener .
func ToOnetListener(listener net.Listener, onet *OverlayNetwork) (Listener, error) {

	laddr, err := FromNetAddr(listener.Addr())

	if err != nil {
		return nil, err
	}

	return &netListenerWrapper{
		listener: listener,
		laddr:    laddr,
		onet:     onet,
	}, nil
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
	laddr, err := listener.Addr().ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	return &onetListenerWrapper{
		listener: listener,
		laddr:    laddr,
	}, nil

}
