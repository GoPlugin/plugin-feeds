package main

import (
	"github.com/hashicorp/go-plugin"

	"github.com/goplugin/plugin-common/pkg/loop"

	"github.com/goplugin/plugin-feeds/median"
)

const (
	loggerName = "PluginMedian"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	p := median.NewPlugin(s.Logger)
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p)

	stop := make(chan struct{})
	defer close(stop)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginMedianHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginMedianName: &loop.GRPCPluginMedian{
				PluginServer: p,
				BrokerConfig: loop.BrokerConfig{
					StopCh:   stop,
					Logger:   s.Logger,
					GRPCOpts: s.GRPCOpts,
				},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}
