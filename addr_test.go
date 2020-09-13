package netty

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type testJSON struct {
	Name Addr
}

func TestAddr(t *testing.T) {
	addr, err := NewAddr("/ip4/127.0.0.1/udp/1812")

	require.NoError(t, err)

	require.NotNil(t, addr)

	buff, err := json.Marshal(addr)

	require.NoError(t, err)

	require.Equal(t, `"/ip4/127.0.0.1/udp/1812"`, string(buff))

	var j *testJSON

	err = json.Unmarshal([]byte(`
	{
		"name": "/ip4/127.0.0.1/tcp/1812"
	}
	`), &j)

	require.NoError(t, err)

	require.Equal(t, j.Name.String(), "/ip4/127.0.0.1/tcp/1812")
}
