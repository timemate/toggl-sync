package trackers

import "time"

type TimeEntry interface {
	Id() string
	Description() string
	Start() time.Time
	Stop() time.Time
	Tags() []string
	Source() string
}

type TrackerApi interface {
	GetTimeEntries(start time.Time, end time.Time) ([]TimeEntry, error)
}
