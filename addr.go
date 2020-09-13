package netty

import (
	"fmt"
	"strings"

	"github.com/libs4go/errors"
)

// SubAddr .
type SubAddr interface {
	Protocol() string
	Value() string
}

// Protocol .
type Protocol struct {
	Name       string
	HasValue   bool
	CheckValue func(value string) error
}

func (p *Protocol) String() string {
	return p.Name
}

var protocols = make(map[string]*Protocol)

// RegisterProtocol .
func RegisterProtocol(protocol *Protocol) error {

	_, ok := protocols[protocol.String()]

	if ok {
		return errors.Wrap(ErrExists, "protocol %s already register", protocol)
	}

	protocols[protocol.String()] = protocol

	return nil
}

// Addr address support overlay network
type Addr struct {
	subaddrs []SubAddr
}

// MarshalJSON .
func (addr *Addr) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, addr.String())), nil
}

// UnmarshalJSON .
func (addr *Addr) UnmarshalJSON(val []byte) error {

	str := strings.TrimPrefix(string(val), "\"")

	str = strings.TrimSuffix(str, "\"")

	subAddrs, err := newAddr(str)

	if err != nil {
		return err
	}

	addr.subaddrs = subAddrs

	return nil
}

func (addr *Addr) String() string {

	var comps []string

	for _, subAddr := range addr.subaddrs {

		comps = append(comps, subAddr.Protocol())

		if subAddr.Value() != "" {
			comps = append(comps, subAddr.Value())
		}
	}

	return "/" + strings.Join(comps, "/")
}

// SubAddrs .
func (addr *Addr) SubAddrs() []SubAddr {
	return addr.subaddrs
}

// JoinAddr .
func JoinAddr(subAddrs ...SubAddr) *Addr {
	return &Addr{
		subaddrs: subAddrs,
	}
}

// NewAddr create addr from string format
func NewAddr(addr string) (*Addr, error) {
	subAddrs, err := newAddr(addr)

	if err != nil {
		return nil, err
	}

	return &Addr{
		subaddrs: subAddrs,
	}, nil
}

func newAddr(addr string) ([]SubAddr, error) {
	comps := strings.Split(addr, "/")

	var err error

	var subAddr SubAddr

	var subAddrs []SubAddr

	for {

		subAddr, comps, err = newSubAddr(comps)

		if err != nil {
			return nil, err
		}

		subAddrs = append(subAddrs, subAddr)

		if len(comps) == 0 {
			return subAddrs, nil
		}
	}
}

type subAddr struct {
	name  string
	value string
}

func (sa *subAddr) Protocol() string {
	return sa.name
}

func (sa *subAddr) Value() string {
	return sa.value
}

func newSubAddr(comps []string) (SubAddr, []string, error) {

	for len(comps) > 0 {
		if comps[0] == "" {
			comps = comps[1:]
			continue
		}

		break
	}

	if len(comps) == 0 {
		return nil, nil, errors.Wrap(ErrParams, "comps can't be nil")
	}

	name := comps[0]

	protocol, ok := protocols[name]

	if !ok {
		return nil, nil, errors.Wrap(ErrNotFound, "protocol %s not found", name)
	}

	sa := &subAddr{
		name: name,
	}

	comps = comps[1:]

	if protocol.HasValue {
		if len(comps) == 0 {
			return nil, nil, errors.Wrap(ErrParams, "protocol %s require value", name)
		}

		if err := protocol.CheckValue(comps[0]); err != nil {
			return nil, nil, errors.Wrap(err, "protocol %s value check error", name)
		}

		sa.value = comps[0]

		comps = comps[1:]
	}

	return sa, comps, nil
}
