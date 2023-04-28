package task_tracker

import (
	"godep.io/timemate/pkg/time_tracker"
)

type IProject interface {
	GetId() string
	GetName() string
}

type ITask interface {
	GetId() string
	GetEntries() []time_tracker.ITimeEntry
	GetProject() IProject
}

type ITaskTracker interface {
	GetTasks(ids []string) ([]ITask, error)
	UpdateTasks(tasks []ITask) error
}
