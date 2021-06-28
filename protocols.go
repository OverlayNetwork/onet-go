package onet

import (
	"net"
	"strconv"

	"github.com/libs4go/errors"
)

func checkPort(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return errors.Wrap(ErrProtocolValue, "invalid ip4 address %s", value)
	}
	if i >= 65536 {
		return errors.Wrap(ErrProtocolValue, "failed to parse port addr: %s", "greater than 65536")
	}

	return nil
}

func checkIP(value string) error {
	parsed := net.ParseIP(value)

	if parsed == nil {
		return errors.Wrap(ErrProtocolValue, "invalid ip4 address %s", value)
	}

	return nil
}

var builtinProtocols = []*Protocol{
	{
		Name:       "ip",
		CheckValue: checkIP,
	},
	{
		Name:       "udp",
		CheckValue: checkPort,
	},
	{
		Name:       "tcp",
		CheckValue: checkPort,
	},
}

func init() {
	if err := RegisterProtocols(builtinProtocols...); err != nil {
		panic(err)
	}
}
