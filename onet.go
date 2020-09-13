package onet

import (
	"github.com/libs4go/errors"
)

// OverlayNetwork .
type OverlayNetwork struct {
	Addr              *Addr
	NavtiveAddr       *Addr
	NativeTransport   NativeTransport
	MuxAddrs          []*Addr
	MuxTransports     []MuxTransport
	OverlayAddrs      []*Addr
	OverlayTransports []OverlayTransport
	Config            *Config
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
		Addr:   addr,
		Config: config,
	}

	for i := 1; i < count; i++ {
		current := subAddrs[count-i]

		transport, ok := transports[current.Protocol()]

		if !ok {
			return nil, errors.Wrap(ErrNotFound, "transport support protocol %s not found", current.Protocol())
		}

		switch t := transport.(type) {
		case NativeTransport:
			result.NavtiveAddr = JoinAddr(subAddrs[0 : count-i+1]...)

			result.NativeTransport = t

			for i, j := 0, len(result.MuxAddrs)-1; i < j; i, j = i+1, j-1 {
				result.MuxAddrs[i], result.MuxAddrs[j] = result.MuxAddrs[j], result.MuxAddrs[i]
			}

			for i, j := 0, len(result.MuxTransports)-1; i < j; i, j = i+1, j-1 {
				result.MuxTransports[i], result.MuxTransports[j] = result.MuxTransports[j], result.MuxTransports[i]
			}

			for i, j := 0, len(result.OverlayAddrs)-1; i < j; i, j = i+1, j-1 {
				result.OverlayAddrs[i], result.OverlayAddrs[j] = result.OverlayAddrs[j], result.OverlayAddrs[i]
			}

			for i, j := 0, len(result.OverlayTransports)-1; i < j; i, j = i+1, j-1 {
				result.OverlayTransports[i], result.OverlayTransports[j] = result.OverlayTransports[j], result.OverlayTransports[i]
			}

			return result, nil
		case OverlayTransport:
			result.OverlayTransports = append(result.OverlayTransports, t)
			result.OverlayAddrs = append(result.OverlayAddrs, JoinAddr(current))
		case MuxTransport:
			result.MuxTransports = append(result.MuxTransports, t)
			result.MuxAddrs = append(result.MuxAddrs, JoinAddr(current))
			result.OverlayTransports = append(result.OverlayTransports, t)
			result.OverlayAddrs = append(result.OverlayAddrs, JoinAddr(current))
		}
	}

	return nil, errors.Wrap(ErrNotFound, "expect native transport")
}
