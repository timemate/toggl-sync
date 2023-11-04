package main

import (
	"encoding/gob"
	"time"

	"github.com/hashicorp/go-plugin"
	plugConfig "godep.io/timemate/pkg/config"
	"godep.io/timemate/pkg/time_tracker"
	"godep.io/timemate/pkg/time_tracker/shared"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(time_tracker.TimeEntry{})
	gob.Register(time_tracker.Project{})
	gob.Register(time_tracker.Client{})
}

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
			//"toggl":      &shared.TimeTrackerPlugin{Impl: impl},
			"toggl_grpc": &shared.TimeTrackerGRPCPlugin{Impl: impl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
