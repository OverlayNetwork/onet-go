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
		"backend":"console",
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
		Name: "tls",
	},
	{
		Name:   "kcp",
		Native: true,
	},
}

type mockTransport struct {
	name string
}

func (mock *mockTransport) String() string {
	return mock.Protocol()
}

func (mock *mockTransport) Protocol() string {
	return mock.name
}

func (mock *mockTransport) Client(ctx context.Context, onet *OverlayNetwork, addr *Addr, next Next) (Conn, error) {
	return nil, nil
}

func (mock *mockTransport) Server(ctx context.Context, onet *OverlayNetwork, addr *Addr, next Next) (Conn, error) {
	return nil, nil
}

func init() {
	if err := RegisterProtocols(mockProtocols...); err != nil {
		panic(err)
	}

	if err := RegisterTransports(&mockTransport{name: "kcp"}, &mockTransport{name: "tls"}, &mockTransport{name: "mux"}); err != nil {
		panic(err)
	}
}

func TestParseOverlayNetwork(t *testing.T) {

	raddr, err := NewAddr("/ip/127.0.0.1/udp/1812/kcp/tls/mux")

	require.NoError(t, err)

	require.NotNil(t, raddr)

	network, err := ParseOverlayNetwork(raddr)

	require.NoError(t, err)

	require.NotNil(t, network)

	require.Equal(t, len(network.Addrs), 3)

	require.Equal(t, network.Addrs[0].String(), "/ip/127.0.0.1/udp/1812/kcp/tls/mux")
	require.Equal(t, network.Addrs[1].String(), "/ip/127.0.0.1/udp/1812/kcp/tls")
	require.Equal(t, network.Addrs[2].String(), "/ip/127.0.0.1/udp/1812/kcp")

}

func BenchmarkParseOverlayNetwork(t *testing.B) {
	for i := 0; i < t.N; i++ {
		raddr, err := NewAddr("/ip/127.0.0.1/udp/1812/kcp/mux")

		require.NoError(t, err)

		require.NotNil(t, raddr)

		network, err := ParseOverlayNetwork(raddr)

		require.NoError(t, err)

		require.NotNil(t, network)
	}
}
