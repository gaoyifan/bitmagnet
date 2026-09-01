package server

import (
	"context"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewDefaultConfigLocalAddress(t *testing.T) {
	require.Equal(t, "0.0.0.0", NewDefaultConfig().LocalAddress)
}

func TestNewBindsConfiguredLocalAddress(t *testing.T) {
	result := New(Params{
		Config: Config{
			LocalAddress: "127.0.0.1",
			Port:         0,
		},
		Logger: zap.NewNop().Sugar(),
	})

	actualServer, err := result.Server.Get()
	require.NoError(t, err)
	limiter := actualServer.(queryLimiter)
	prometheusWrapper := limiter.server.(prometheusServerWrapper)
	healthWrapper := prometheusWrapper.server.(healthCollector)
	baseServer := healthWrapper.baseServer.(*server)
	require.Equal(t, netip.MustParseAddrPort("127.0.0.1:0"), baseServer.localAddr)
	require.NoError(t, result.AppHook.OnStop(context.Background()))
}

func TestNewRejectsInvalidLocalAddress(t *testing.T) {
	for _, localAddress := range []string{"invalid", "::1"} {
		t.Run(localAddress, func(t *testing.T) {
			result := New(Params{
				Config: Config{LocalAddress: localAddress},
				Logger: zap.NewNop().Sugar(),
			})

			_, err := result.Server.Get()
			require.ErrorContains(t, err, "invalid IPv4 DHT server local address")
		})
	}
}
