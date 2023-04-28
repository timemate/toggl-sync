package task_tracker

import (
	"godep.io/timemate/pkg/time_tracker"
)

type Project struct {
	Id   string
	Name string
}

func (t Project) GetId() string {
	return t.Id
}
func (t Project) GetName() string {
	return t.Name
}

type Task struct {
	Id      string
	Entries []time_tracker.ITimeEntry
	Project IProject
}

func (t Task) GetId() string {
	return t.Id
}
func (t Task) GetProject() IProject {
	return t.Project
}
func (t Task) GetEntries() []time_tracker.ITimeEntry {
	return t.Entries
}
