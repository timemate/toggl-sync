package shared

import (
	"context"
	"net/rpc"

	"godep.io/timemate/pkg/time_tracker"
	"godep.io/timemate/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type TimeTrackerPlugin struct {
	Impl time_tracker.ITimeTracker
}

func (p *TimeTrackerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*TimeTrackerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

type TimeTrackerGRPCPlugin struct {
	plugin.Plugin
	Impl time_tracker.ITimeTracker
}

func (p *TimeTrackerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterTimeTrackerServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *TimeTrackerGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewTimeTrackerClient(c)}, nil
}
