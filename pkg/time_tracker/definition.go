package time_tracker

import "time"

type Client struct {
	Id   string
	Name string
}

func (t Client) GetId() string {
	return t.Id
}
func (t Client) GetName() string {
	return t.Name
}

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

type TimeEntry struct {
	Id          string
	Description string
	Start       time.Time
	Stop        time.Time
	Tags        []string
	Source      string
	Client      IClient
	Project     IProject
}

func (t TimeEntry) GetTags() []string {
	return t.Tags
}

func (t TimeEntry) GetStop() time.Time {
	return t.Stop
}

func (t TimeEntry) GetStart() time.Time {
	return t.Start
}

func (t TimeEntry) GetDescription() string {
	return t.Description
}

func (t TimeEntry) GetId() string {
	return t.Id
}

func (t TimeEntry) GetSource() string {
	return t.Source
}

func (t TimeEntry) GetClient() IClient {
	return t.Client
}

func (t TimeEntry) GetProject() IProject {
	return t.Project
}
