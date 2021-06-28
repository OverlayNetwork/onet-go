package onet

import (
	"github.com/libs4go/errors"
	"github.com/libs4go/slf4go"
)

// OverlayNetwork .
type OverlayNetwork struct {
	slf4go.Logger
	Addr       *Addr
	Addrs      []*Addr
	Transports []Transport
	Config     *Config
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

		result.Addrs = append(result.Addrs, JoinAddr(subAddrs[0:count-i+1]...))

		result.Transports = append(result.Transports, transport)

		if current.Native() {
			return result, nil
		}

	}

	return nil, errors.Wrap(ErrNotFound, "expect native transport")
}

// FindTransports find transports by protocol name
func (network *OverlayNetwork) FindTransports(protocol string) []Transport {
	var transports []Transport

	for _, transport := range network.Transports {
		if transport.Protocol() == protocol {
			transports = append(transports, transport)
		}
	}

	return transports
}
