package main

import (
	"encoding/gob"
	"time"

	plugConfig "godep.io/timemate/pkg/config"
	pkgPlugin "godep.io/timemate/pkg/plugin"
	"godep.io/timemate/pkg/task_tracker"

	"github.com/hashicorp/go-plugin"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(task_tracker.Task{})
	gob.Register(task_tracker.Project{})
}

const pluginName = "jira"

func main() {
	config, err := plugConfig.ReadConfig()
	if err != nil {
		panic(err)
	}
	pluginConfig := config.FindPlugin(pluginName)
	if pluginConfig == nil {
		panic("No plugin configuration found")
	}
	impl, err := NewJiraTracker(*pluginConfig)
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
