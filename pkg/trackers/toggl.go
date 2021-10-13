package trackers

import (
	"strconv"
	"time"

	"github.com/andreaskoch/togglapi"
	"github.com/andreaskoch/togglapi/model"
)

type TogglTracker struct {
	api model.TimeEntryAPI
}

func NewTogglTracker(params map[interface{}]interface{}) (*TogglTracker, error) {
	return &TogglTracker{
		api: togglapi.NewTimeEntryAPI("https://api.track.toggl.com/api/v8", (params["token"]).(string)),
	}, nil
}

func (tg *TogglTracker) GetTimeEntries(start time.Time, end time.Time) ([]TimeEntry, error) {
	entries, err := tg.api.GetTimeEntries(start, end)
	if err != nil {
		return nil, err
	}
	r := make([]TimeEntry, len(entries), len(entries))
	for i, e := range entries {
		r[i] = trackedTime{
			id:          strconv.Itoa(e.ID),
			description: e.Description,
			start:       e.Start,
			stop:        e.Stop,
			tags:        e.Tags,
			source:      "toggl",
		}
	}
	return r, nil
}
