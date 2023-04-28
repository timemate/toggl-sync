package shared

import (
	"time"

	"godep.io/timemate/pkg/time_tracker"
	"godep.io/timemate/proto"

	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.TimeTrackerClient }

func (m *GRPCClient) GetTimeEntries(start time.Time, stop time.Time) ([]time_tracker.TimeEntry, error) {
	resp, err := m.client.GetTimeEntries(context.Background(), &proto.GetTimeEntriesRequest{
		Start: start.Format(time.RFC3339),
		Stop:  stop.Format(time.RFC3339),
	})

	entries := make([]time_tracker.TimeEntry, 0)
	for _, e := range resp.Entries {
		start, err = time.Parse(time.RFC3339, e.Start)
		if err != nil {
			return nil, err
		}
		stop, err = time.Parse(time.RFC3339, e.Stop)
		if err != nil {
			return nil, err
		}
		var client *time_tracker.Client
		if e.Client != nil {
			client = &time_tracker.Client{
				Id:   e.Client.Id,
				Name: e.Client.Name,
			}
		}
		var project *time_tracker.Project
		if e.Project != nil {
			project = &time_tracker.Project{
				Id:   e.Project.Id,
				Name: e.Project.Name,
			}
		}

		entries = append(entries, time_tracker.TimeEntry{
			Id:          e.Id,
			Description: e.Description,
			Start:       start,
			Stop:        stop,
			Tags:        e.Tags,
			Source:      e.Source,
			Client:      client,
			Project:     project,
		})
	}

	return entries, err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl time_tracker.ITimeTracker
}

func (m *GRPCServer) GetTimeEntries(
	ctx context.Context,
	req *proto.GetTimeEntriesRequest) (*proto.GetTimeEntriesResponse, error) {

	start, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		return nil, err
	}
	stop, err := time.Parse(time.RFC3339, req.Stop)
	if err != nil {
		return nil, err
	}

	entries, err := m.Impl.GetTimeEntries(start, stop)

	resp := make([]*proto.TimeEntry, 0)
	for _, e := range entries {
		var client *proto.TimeEntry_Client
		if e.GetClient() != nil {
			client = &proto.TimeEntry_Client{
				Id:   e.GetClient().GetId(),
				Name: e.GetClient().GetName(),
			}
		}
		var project *proto.TimeEntry_Project
		if e.GetProject() != nil {
			project = &proto.TimeEntry_Project{
				Id:   e.GetProject().GetId(),
				Name: e.GetProject().GetName(),
			}
		}

		resp = append(resp, &proto.TimeEntry{
			Id:          e.GetId(),
			Description: e.GetDescription(),
			Start:       e.GetStart().Format(time.RFC3339),
			Stop:        e.GetStop().Format(time.RFC3339),
			Tags:        e.GetTags(),
			Source:      e.GetSource(),
			Client:      client,
			Project:     project,
		})
	}

	return &proto.GetTimeEntriesResponse{
		Entries: resp,
	}, err
}
