package shared

import (
	"godep.io/timemate/pkg/task_tracker"
	"net/rpc"
)

// RPCClient is an implementation of KV that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) GetTimeEntries(ids []string) ([]task_tracker.ITask, error) {
	var resp []task_tracker.ITask
	err := m.client.Call("Plugin.GetTasks", map[string]interface{}{
		"ids": ids,
	}, &resp)
	return resp, err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl task_tracker.ITaskTracker
}

func (m *RPCServer) GetTasks(args map[string]interface{}, resp *[]task_tracker.ITask) error {
	v, err := m.Impl.GetTasks(args["ids"].([]string))
	*resp = v
	return err
}
