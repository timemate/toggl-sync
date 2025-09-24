package trackers

import "time"

type Client interface {
	Id() string
	Name() string
}

type Project interface {
	Id() string
	ClientId() string
	Name() string
}

type TimeEntry interface {
	Id() string
	Description() string
	Start() time.Time
	Stop() time.Time
	Tags() []string
	Source() string
	Client() Client
	Project() Project
}

type TrackerApi interface {
	GetTimeEntries(start time.Time, end time.Time) ([]TimeEntry, error)
}
