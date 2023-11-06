package plugin

import (
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"godep.io/timemate/pkg/config"
	"godep.io/timemate/pkg/time_tracker"
)

func getRPCClient[T time_tracker.ITimeTracker](
	conf config.IConfig,
	pluginPath string,
	pluginName string,
	isGRPC bool,
) (T, *plugin.Client, error) {
	var r T
	m, err := GetPluginMap(conf, pluginName, isGRPC)
	if err != nil {
		return r, nil, err
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         m,
		Cmd:             exec.Command("sh", "-c", pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	})

	rpcClient, err := client.Client()
	if err != nil {
		return r, nil, err
	}

	pgName := pluginName
	if isGRPC {
		pgName = PluginWithGrpc(pgName)
	}
	raw, err := rpcClient.Dispense(pgName)
	if err != nil {
		return r, nil, err
	}

	return raw.(T), client, nil
}

func GetRPCClient[T time_tracker.ITimeTracker](
	conf config.IConfig,
	pluginPath string,
	pluginName string,
) (T, *plugin.Client, error) {
	return getRPCClient[T](conf, pluginPath, pluginName, false)
}

func GetGRPCClient[T time_tracker.ITimeTracker](
	conf config.IConfig,
	pluginPath string,
	pluginName string,
) (T, *plugin.Client, error) {
	return getRPCClient[T](conf, pluginPath, pluginName, true)
}
