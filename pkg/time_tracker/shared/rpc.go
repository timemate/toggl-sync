package shared

import (
	"net/rpc"
	"time"

	"godep.io/timemate/pkg/time_tracker"
)

type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) GetTimeEntries(start time.Time, stop time.Time) ([]time_tracker.ITimeEntry, error) {
	var resp []time_tracker.ITimeEntry
	err := m.client.Call("Plugin.GetTimeEntries", map[string]interface{}{
		"start": start,
		"stop":  stop,
	}, &resp)
	return resp, err
}

type RPCServer struct {
	Impl time_tracker.ITimeTracker
}

func (m *RPCServer) GetTimeEntries(args map[string]interface{}, resp *[]time_tracker.ITimeEntry) error {
	v, err := m.Impl.GetTimeEntries(args["start"].(time.Time), args["stop"].(time.Time))
	*resp = v
	return err
}
