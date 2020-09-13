package netty

import (
	"github.com/libs4go/errors"
)

// OverlayNetwork .
type OverlayNetwork struct {
	navtiveAddr       *Addr
	nativeTransport   NativeTransport
	muxAddrs          []SubAddr
	muxTransports     []MuxTransport
	overlayAddrs      []SubAddr
	overlayTransports []OverlayTransport
	config            *Config
}

// ParseOverlayNetwork parse addr to generate overlay network config
func ParseOverlayNetwork(addr *Addr, options ...Option) (*OverlayNetwork, error) {

	config := NewConfig()

	for _, option := range options {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	subAddrs := addr.SubAddrs()

	count := len(subAddrs)

	var result = &OverlayNetwork{
		config: config,
	}

	for i := 1; i < count; i++ {
		current := subAddrs[count-i]

		transport, ok := transports[current.Protocol()]

		if !ok {
			return nil, errors.Wrap(ErrNotFound, "transport support protocol %s not found", current.Protocol())
		}

		switch t := transport.(type) {
		case NativeTransport:
			result.navtiveAddr = JoinAddr(subAddrs[0 : count-i+1]...)

			result.nativeTransport = t

			for i, j := 0, len(result.muxAddrs)-1; i < j; i, j = i+1, j-1 {
				result.muxAddrs[i], result.muxAddrs[j] = result.muxAddrs[j], result.muxAddrs[i]
			}

			for i, j := 0, len(result.muxTransports)-1; i < j; i, j = i+1, j-1 {
				result.muxTransports[i], result.muxTransports[j] = result.muxTransports[j], result.muxTransports[i]
			}

			for i, j := 0, len(result.overlayAddrs)-1; i < j; i, j = i+1, j-1 {
				result.overlayAddrs[i], result.overlayAddrs[j] = result.overlayAddrs[j], result.overlayAddrs[i]
			}

			for i, j := 0, len(result.overlayTransports)-1; i < j; i, j = i+1, j-1 {
				result.overlayTransports[i], result.overlayTransports[j] = result.overlayTransports[j], result.overlayTransports[i]
			}

			return result, nil
		case OverlayTransport:
			result.overlayTransports = append(result.overlayTransports, t)
			result.overlayAddrs = append(result.overlayAddrs, current)
		case MuxTransport:
			result.muxTransports = append(result.muxTransports, t)
			result.muxAddrs = append(result.muxAddrs, current)
		}
	}

	return nil, errors.Wrap(ErrNotFound, "expect native transport")
}

// Listen listen on the local overlay address with config
func (network *OverlayNetwork) Listen() (Listener, error) {
	return nil, nil
}

// Dial dial to the remote overlay address with config
func (network *OverlayNetwork) Dial() (Conn, error) {
	return nil, nil
}
