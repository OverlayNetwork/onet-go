package onet

// Listener overlay network listener equal to net.Listener
type Listener interface {
	Accept() (Conn, error)

	Close() error

	// Addr returns the listener's network address.
	Addr() Addr
}

// Listen listen on the local overlay address
func Listen(addr *Addr, options ...Option) (Listener, error) {
	network, err := ParseOverlayNetwork(addr, options...)

	if err != nil {
		return nil, err
	}

	return network.Listen()
}
