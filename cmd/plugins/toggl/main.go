package main

import (
	"encoding/gob"
	"time"

	plugConfig "godep.io/timemate/pkg/config"
	pkgPlugin "godep.io/timemate/pkg/plugin"
	"godep.io/timemate/pkg/time_tracker"

	"github.com/hashicorp/go-plugin"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(time_tracker.TimeEntry{})
	gob.Register(time_tracker.Project{})
	gob.Register(time_tracker.Client{})
}

const pluginName = "toggl"

func main() {
	config, err := plugConfig.ReadConfig()
	if err != nil {
		panic(err)
	}
	pluginConfig := config.FindPlugin(pluginName)
	if pluginConfig == nil {
		panic("No plugin configuration found")
	}
	impl, err := NewTogglTracker(*pluginConfig)
	if err != nil {
		panic(err)
	}
	m, err := pkgPlugin.GetPluginMap(config, pluginName, true, impl)
	if err != nil {
		panic("No plugin configuration found")
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pkgPlugin.Handshake,
		Plugins:         m,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
