package time_tracker

import "time"

type IClient interface {
	GetId() string
	GetName() string
}

type IProject interface {
	GetId() string
	GetName() string
}

type ITimeEntry interface {
	GetId() string
	GetDescription() string
	GetStart() time.Time
	GetStop() time.Time
	GetTags() []string
	GetSource() string
	GetClient() IClient
	GetProject() IProject
}

type ITimeTracker interface {
	GetTimeEntries(start time.Time, end time.Time) ([]ITimeEntry, error)
}
