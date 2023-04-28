package main

import (
	plugConfig "godep.io/timemate/pkg/config"
	"godep.io/timemate/pkg/time_tracker/shared"

	"github.com/hashicorp/go-plugin"
)

func main() {
	config, err := plugConfig.ReadConfig()
	if err != nil {
		panic(err)
	}
	pluginConfig := config.FindPlugin("toggl")
	if pluginConfig == nil {
		panic("No plugin configuration found")
	}
	impl, err := NewTogglTracker(*pluginConfig)
	if err != nil {
		panic(err)
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"toggl": &shared.TimeTrackerPlugin{Impl: impl},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
