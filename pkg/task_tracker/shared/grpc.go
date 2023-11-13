package shared

import (
	"time"

	"godep.io/timemate/pkg/task_tracker"
	"godep.io/timemate/pkg/time_tracker"
	"godep.io/timemate/proto"

	"golang.org/x/net/context"
)

func getProtoFromEntity(t task_tracker.ITask) (*proto.Task, error) {
	entries := make([]*proto.TimeEntry, 0)
	for _, e := range t.GetEntries() {
		entries = append(entries, &proto.TimeEntry{
			Id:          e.GetId(),
			Description: e.GetDescription(),
			Start:       e.GetStart().Format(time.RFC3339),
			Stop:        e.GetStop().Format(time.RFC3339),
			Tags:        e.GetTags(),
			Source:      e.GetSource(),
		})
	}
	var project *proto.Task_Project
	if t.GetProject() != nil {
		project = &proto.Task_Project{
			Id:   t.GetProject().GetId(),
			Name: t.GetProject().GetName(),
		}
	}
	return &proto.Task{
		Id:      t.GetId(),
		Entries: entries,
		Project: project,
	}, nil
}

func getProtoListFromEntityList(tasks []task_tracker.ITask) ([]*proto.Task, error) {
	protoTasks := make([]*proto.Task, 0)
	for _, t := range tasks {
		protoTask, err := getProtoFromEntity(t)
		if err != nil {
			return nil, err
		}
		protoTasks = append(protoTasks, protoTask)
	}
	return protoTasks, nil
}

func getEntityFromProto(p *proto.Task) (task_tracker.ITask, error) {
	entries := make([]time_tracker.ITimeEntry, 0)
	for _, e := range p.Entries {
		start, err := time.Parse(time.RFC3339, e.Start)
		if err != nil {
			return nil, err
		}
		stop, err := time.Parse(time.RFC3339, e.Stop)
		if err != nil {
			return nil, err
		}
		entries = append(entries, time_tracker.TimeEntry{
			Id:          e.Id,
			Description: e.Description,
			Start:       start,
			Stop:        stop,
			Tags:        e.Tags,
			Source:      e.Source,
		})
	}
	var project *time_tracker.Project
	if p.Project != nil {
		project = &time_tracker.Project{
			Id:   p.Project.Id,
			Name: p.Project.Name,
		}
	}
	return task_tracker.Task{
		Id:      p.Id,
		Entries: entries,
		Project: project,
	}, nil
}

func getEntityListFromProtoList(p []*proto.Task) ([]task_tracker.ITask, error) {
	tasks := make([]task_tracker.ITask, 0)
	for _, e := range p {
		task, err := getEntityFromProto(e)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.TaskTrackerClient }

func (m *GRPCClient) GetTasks(ids []string) ([]task_tracker.ITask, error) {
	resp, err := m.client.GetTasks(context.Background(), &proto.GetTasksRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	return getEntityListFromProtoList(resp.GetTasks())
}

func (m *GRPCClient) UpdateTasks(tasks []task_tracker.ITask) error {
	in, err := getProtoListFromEntityList(tasks)
	if err != nil {
		return err
	}
	_, err = m.client.UpdateTasks(context.Background(), &proto.UpdateTasksRequest{
		Tasks: in,
	})
	if err != nil {
		return err
	}
	return nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl task_tracker.ITaskTracker
}

func (m *GRPCServer) GetTasks(
	ctx context.Context,
	req *proto.GetTasksRequest) (*proto.GetTasksResponse, error) {

	tasks, err := m.Impl.GetTasks(req.Ids)
	if err != nil {
		return nil, err
	}
	resp, err := getProtoListFromEntityList(tasks)
	if err != nil {
		return nil, err
	}

	return &proto.GetTasksResponse{
		Tasks: resp,
	}, nil
}

func (m *GRPCServer) UpdateTasks(
	ctx context.Context,
	req *proto.UpdateTasksRequest) (*proto.UpdateTasksResponse, error) {

	tasks, err := getEntityListFromProtoList(req.Tasks)
	if err != nil {
		return nil, err
	}

	err = m.Impl.UpdateTasks(tasks)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateTasksResponse{}, nil
}
