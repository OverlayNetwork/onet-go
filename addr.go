package onet

import (
	"fmt"
	"net"
	"strings"

	"github.com/libs4go/errors"
)

// SubAddr .
type SubAddr interface {
	Protocol() string
	Value() string
	Native() bool
}

// Protocol .
type Protocol struct {
	Name string
	// HasValue   bool
	CheckValue func(value string) error
	// Indicate native protocol
	Native bool
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

// RegisterProtocols .
func RegisterProtocols(protocols ...*Protocol) error {
	for _, protocol := range protocols {
		if err := RegisterProtocol(protocol); err != nil {
			return err
		}
	}

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

// Join .
func (addr *Addr) Join(subAddrs ...SubAddr) *Addr {
	return &Addr{append(addr.subaddrs, subAddrs...)}

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
	name   string
	value  string
	native bool
}

func (sa *subAddr) Protocol() string {
	return sa.name
}

func (sa *subAddr) Value() string {
	return sa.value
}

func (sa *subAddr) Native() bool {
	return sa.native
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
		name:   name,
		native: protocol.Native,
	}

	comps = comps[1:]

	if protocol.CheckValue != nil {
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

// ResolveNetAddr .
func (addr *Addr) ResolveNetAddr() (net.Addr, *Addr, error) {
	subAddrs := addr.SubAddrs()

	if len(subAddrs) < 2 {
		return nil, nil, errors.Wrap(ErrAddr, "can't parse addr %s as net addr", addr.String())
	}

	if subAddrs[0].Protocol() != "ip" {
		return nil, nil, errors.Wrap(ErrAddr, "can't parse addr %s as net addr", addr.String())
	}

	if subAddrs[1].Protocol() != "tcp" && subAddrs[1].Protocol() != "udp" {
		return nil, nil, errors.Wrap(ErrAddr, "can't parse addr %s as net addr", addr.String())
	}

	str := fmt.Sprintf("%s:%s", subAddrs[0].Value(), subAddrs[1].Value())

	relativeAddr := JoinAddr()

	if len(subAddrs) > 2 {
		relativeAddr = JoinAddr(subAddrs[2:]...)
	}

	if subAddrs[1].Protocol() == "tcp" {
		tcpAddr, err := net.ResolveTCPAddr("tcp", str)

		if err != nil {
			return nil, nil, errors.Wrap(err, "resolve net addr %s error", str)
		}

		return tcpAddr, relativeAddr, nil
	}

	udpAddr, err := net.ResolveUDPAddr("udp", str)

	if err != nil {
		return nil, nil, errors.Wrap(err, "resolve net addr %s error", str)
	}

	return udpAddr, relativeAddr, nil
}

// FromNetAddr create addr from net.Addr
func FromNetAddr(addr net.Addr) (*Addr, error) {

	switch addr.Network() {
	case "tcp", "tcp6", "tcp4":
		tcpAddr := addr.(*net.TCPAddr)
		if tcpAddr.Zone != "" {
			return NewAddr(fmt.Sprintf("/ip/%s%%%s/tcp/%d", tcpAddr.IP.String(), tcpAddr.Zone, tcpAddr.Port))
		}

		return NewAddr(fmt.Sprintf("/ip/%s/tcp/%d", tcpAddr.IP.String(), tcpAddr.Port))

	case "udp", "udp4", "udp6":
		udpAddr := addr.(*net.UDPAddr)
		if udpAddr.Zone != "" {
			return NewAddr(fmt.Sprintf("/ip/%s%%%s/udp/%d", udpAddr.IP.String(), udpAddr.Zone, udpAddr.Port))
		}

		return NewAddr(fmt.Sprintf("/ip/%s/udp/%d", udpAddr.IP.String(), udpAddr.Port))

	default:
		return nil, errors.Wrap(ErrAddr, "unsupport net.Addr %s %s", addr.Network(), addr.String())
	}
}
