package shared

import (
	"net/rpc"
	"time"

	"godep.io/timemate/pkg/time_tracker"
)

// RPCClient is an implementation of KV that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) GetTimeEntries(start time.Time, stop time.Time) ([]time_tracker.ITimeEntry, error) {
	var resp []time_tracker.ITimeEntry
	err := m.client.Call("Plugin.GetTimeEntries", map[string]interface{}{
		"start": start,
		"stop":  stop,
	}, &resp)
	return resp, err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl time_tracker.ITimeTracker
}

func (m *RPCServer) GetTimeEntries(args map[string]interface{}, resp *[]time_tracker.ITimeEntry) error {
	v, err := m.Impl.GetTimeEntries(args["start"].(time.Time), args["stop"].(time.Time))
	*resp = v
	return err
}
