package onet

import (
	"github.com/libs4go/errors"
	"github.com/libs4go/slf4go"
)

// OverlayNetwork .
type OverlayNetwork struct {
	slf4go.Logger
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

	var result = &OverlayNetwork{
		Logger: slf4go.Get("overlay-network"),
		Addr:   addr,
	}

	config := NewConfig()

	for _, option := range options {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	result.Config = config

	subAddrs := addr.SubAddrs()

	count := len(subAddrs)

	for i := 1; i < count; i++ {

		current := subAddrs[count-i]

		result.D("parse sub protocol {@p}", current.Protocol())

		transport, ok := transports[current.Protocol()]

		if !ok {
			return nil, errors.Wrap(ErrNotFound, "transport support protocol %s not found", current.Protocol())
		}

		// warning !!! Never change the match order
		// because MuxTransport mixin NativeTransport and OverlayTransport that must be first test
		switch t := transport.(type) {
		case MuxTransport:
			result.D("transport {@t} is mux transport", transport)
			result.MuxTransports = append(result.MuxTransports, t)
			result.MuxAddrs = append(result.MuxAddrs, JoinAddr(current))
			result.OverlayTransports = append(result.OverlayTransports, t)
			result.OverlayAddrs = append(result.OverlayAddrs, JoinAddr(current))
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
		}
	}

	return nil, errors.Wrap(ErrNotFound, "expect native transport")
}
