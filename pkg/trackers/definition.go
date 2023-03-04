package trackers

import "time"

type trackedTimeClient struct {
	id   string
	name string
}

func (t trackedTimeClient) Id() string {
	return t.id
}
func (t trackedTimeClient) Name() string {
	return t.name
}

type trackedTimeProject struct {
	id   string
	name string
}

func (t trackedTimeProject) Id() string {
	return t.id
}
func (t trackedTimeProject) Name() string {
	return t.name
}

type trackedTime struct {
	id          string
	description string
	start       time.Time
	stop        time.Time
	tags        []string
	source      string
	client      Client
	project     Project
}

func (t trackedTime) Tags() []string {
	return t.tags
}

func (t trackedTime) Stop() time.Time {
	return t.stop
}

func (t trackedTime) Start() time.Time {
	return t.start
}

func (t trackedTime) Description() string {
	return t.description
}

func (t trackedTime) Id() string {
	return t.id
}

func (t trackedTime) Source() string {
	return t.source
}

func (t trackedTime) Client() Client {
	return t.client
}

func (t trackedTime) Project() Project {
	return t.project
}
