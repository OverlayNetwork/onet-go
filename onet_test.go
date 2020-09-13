package onet

import (
	"context"
	"testing"

	"github.com/libs4go/scf4go"
	_ "github.com/libs4go/scf4go/codec" //
	"github.com/libs4go/scf4go/reader/memory"
	"github.com/libs4go/slf4go"
	_ "github.com/libs4go/slf4go/backend/console" //
	"github.com/stretchr/testify/require"
)

var loggerjson = `
{
	"default":{
		"backend":"null",
		"level":"debug"
	},
	"backend":{
		"console":{
			"formatter":{
				"output": "@t @l @s @m"
			}
		}
	}
}
`

func init() {
	config := scf4go.New()

	err := config.Load(memory.New(memory.Data(loggerjson, "json")))

	if err != nil {
		panic(err)
	}

	err = slf4go.Config(config)

	if err != nil {
		panic(err)
	}
}

var mockProtocols = []*Protocol{
	{
		Name: "mux",
	},
	{
		Name: "kcp",
	},
}

type mockMux struct{}

func (mock *mockMux) String() string {
	return mock.Protocol()
}

func (mock *mockMux) Protocol() string {
	return "mux"
}

func (mock *mockMux) Listen(onet *OverlayNetwork, laddr *Addr, config *Config) (Listener, error) {
	return nil, nil
}

func (mock *mockMux) Dial(ctx context.Context, onet *OverlayNetwork, raddr *Addr, config *Config) (Conn, error) {
	return nil, nil
}

func (mock *mockMux) Client(onet *OverlayNetwork, conn Conn, raddr *Addr, config *Config) (Conn, error) {
	return nil, nil
}

func (mock *mockMux) Server(onet *OverlayNetwork, conn Conn, laddr *Addr, config *Config) (Conn, error) {
	return nil, nil
}

type mockKCP struct{}

func (mock *mockKCP) String() string {
	return mock.Protocol()
}

func (mock *mockKCP) Protocol() string {
	return "kcp"
}

func (mock *mockKCP) Client(onet *OverlayNetwork, conn Conn, raddr *Addr, config *Config) (Conn, error) {
	return nil, nil
}

func (mock *mockKCP) Server(onet *OverlayNetwork, conn Conn, laddr *Addr, config *Config) (Conn, error) {
	return nil, nil
}

type mockUDP struct{}

func (mock *mockUDP) String() string {
	return mock.Protocol()
}

func (mock *mockUDP) Protocol() string {
	return "udp"
}

func (mock *mockUDP) Listen(onet *OverlayNetwork, laddr *Addr, config *Config) (Listener, error) {
	return nil, nil
}

func (mock *mockUDP) Dial(ctx context.Context, onet *OverlayNetwork, raddr *Addr, config *Config) (Conn, error) {
	return nil, nil
}

func init() {
	if err := RegisterProtocols(mockProtocols...); err != nil {
		panic(err)
	}

	if err := RegisterTransports(&mockKCP{}, &mockUDP{}, &mockMux{}); err != nil {
		panic(err)
	}
}

func TestParseOverlayNetwork(t *testing.T) {

	raddr, err := NewAddr("/ip4/127.0.0.1/udp/1812/kcp/mux")

	require.NoError(t, err)

	require.NotNil(t, raddr)

	network, err := ParseOverlayNetwork(raddr)

	require.NoError(t, err)

	require.NotNil(t, network)

	require.Equal(t, 2, len(network.OverlayTransports))

	require.Equal(t, len(network.OverlayTransports), len(network.OverlayAddrs))

	require.Equal(t, 1, len(network.MuxTransports), 1)

	require.Equal(t, len(network.MuxTransports), len(network.MuxAddrs))

	require.NotNil(t, network.NativeTransport)

	require.NotNil(t, network.NavtiveAddr)

	require.Equal(t, network.Addr.String(), raddr.String())

	require.Equal(t, network.MuxAddrs[0].String(), "/mux")

	require.Equal(t, network.OverlayAddrs[0].String(), "/kcp")
	require.Equal(t, network.OverlayAddrs[1].String(), "/mux")

	require.Equal(t, network.OverlayTransports[1], network.MuxTransports[0])

	require.Equal(t, network.NavtiveAddr.String(), "/ip4/127.0.0.1/udp/1812")
}

func BenchmarkParseOverlayNetwork(t *testing.B) {
	for i := 0; i < t.N; i++ {
		raddr, err := NewAddr("/ip4/127.0.0.1/udp/1812/kcp/mux")

		require.NoError(t, err)

		require.NotNil(t, raddr)

		network, err := ParseOverlayNetwork(raddr)

		require.NoError(t, err)

		require.NotNil(t, network)
	}
}
