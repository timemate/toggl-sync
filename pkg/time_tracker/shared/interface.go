package shared

import (
	"context"
	"net/rpc"

	"godep.io/timemate/pkg/time_tracker"
	"godep.io/timemate/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"timetracker_grpc": &TimeTrackerGRPCPlugin{},
	"timetracker":      &TimeTrackerPlugin{},
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type TimeTrackerPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl time_tracker.ITimeTracker
}

func (p *TimeTrackerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*TimeTrackerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type TimeTrackerGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl time_tracker.ITimeTracker
}

func (p *TimeTrackerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterTimeTrackerServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *TimeTrackerGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewTimeTrackerClient(c)}, nil
}
