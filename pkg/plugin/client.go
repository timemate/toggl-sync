package plugin

import (
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"godep.io/timemate/pkg/config"
)

func getRPCClient(
	conf config.IConfig,
	pluginPath string,
	pluginName string,
	isGRPC bool,
) (interface{}, *plugin.Client, error) {
	var r interface{}
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

	return raw, client, nil
}

func GetRPCClient(
	conf config.IConfig,
	pluginPath string,
	pluginName string,
) (interface{}, *plugin.Client, error) {
	return getRPCClient(conf, pluginPath, pluginName, false)
}

func GetGRPCClient(
	conf config.IConfig,
	pluginPath string,
	pluginName string,
) (interface{}, *plugin.Client, error) {
	return getRPCClient(conf, pluginPath, pluginName, true)
}
