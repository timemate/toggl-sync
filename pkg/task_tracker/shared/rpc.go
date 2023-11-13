package shared

import (
	"godep.io/timemate/pkg/task_tracker"
	"net/rpc"
)

// RPCClient is an implementation of KV that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) GetTasks(ids []string) ([]task_tracker.ITask, error) {
	var resp []task_tracker.ITask
	err := m.client.Call("Plugin.GetTasks", map[string]interface{}{
		"ids": ids,
	}, &resp)
	return resp, err
}

func (m *RPCClient) UpdateTasks(tasks []task_tracker.ITask) error {
	var resp []task_tracker.ITask
	err := m.client.Call("Plugin.UpdateTasks", map[string]interface{}{
		"tasks": tasks,
	}, &resp)
	return err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl task_tracker.ITaskTracker
}

func (m *RPCServer) UpdateTasks(args map[string]interface{}) error {
	err := m.Impl.UpdateTasks(args["tasks"].([]task_tracker.ITask))
	return err
}

func (m *RPCServer) GetTasks(args map[string]interface{}, resp *[]task_tracker.ITask) error {
	v, err := m.Impl.GetTasks(args["ids"].([]string))
	*resp = v
	return err
}
