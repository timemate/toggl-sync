package shared

import (
	"context"
	"net/rpc"

	"godep.io/timemate/pkg/task_tracker"
	"godep.io/timemate/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// This is the implementation of plugin.Plugin so we can serve/consume this.
type TaskTrackerPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl task_tracker.ITaskTracker
}

func (p *TaskTrackerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*TaskTrackerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type TaskTrackerGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl task_tracker.ITaskTracker
}

func (p *TaskTrackerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterTaskTrackerServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *TaskTrackerGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewTaskTrackerClient(c)}, nil
}
