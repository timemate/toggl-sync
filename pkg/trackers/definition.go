package trackers

import "time"

type trackedTime struct {
	id          string
	description string
	start       time.Time
	stop        time.Time
	tags        []string
	source      string
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
