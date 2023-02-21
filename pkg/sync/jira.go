package sync

import (
	"log"

	"github.com/andygrunwald/go-jira"
	"godep.io/timemate/pkg/trackers"
)

type JiraSync struct {
	api   *jira.Client
	login string
}

func NewJiraSync(params map[interface{}]interface{}) (*JiraSync, error) {
	login := (params["login"]).(string)
	token := (params["token"]).(string)
	url := (params["url"]).(string)
	tp := jira.BasicAuthTransport{
		Username: login,
		Password: token,
	}
	client, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil, err
	}
	return &JiraSync{
		api:   client,
		login: login,
	}, nil
}

func findWorklog(worklogs []jira.WorklogRecord, entry trackers.TimeEntry) *jira.WorklogRecord {
	if worklogs == nil || len(worklogs) == 0 {
		return nil
	}
	for _, w := range worklogs {
		if w.Properties == nil || len(w.Properties) == 0 {
			continue
		}
		for _, p := range w.Properties {
			if v, ok := (p.Value).(map[string]interface{}); ok && p.Key == entry.Source() && v[entry.Source()] == entry.Id() {
				return &w
			}
		}
	}
	return nil
}

func (ji *JiraSync) Sync(tasks []Task) (err error) {
	for _, t := range tasks {
		log.Printf("Processing toggl task %s\n", t.Id)
		worklog, _, err := ji.api.Issue.GetWorklogs(t.Id, jira.WithQueryOptions(&jira.AddWorklogQueryOptions{
			Expand: "properties",
		}))
		if err != nil {
			log.Printf("Error occured: %s\n", err)
			continue
		}
		for _, e := range t.entries {
			w := findWorklog(worklog.Worklogs, e)
			st := jira.Time(e.Start())
			diff := e.Stop().Sub(e.Start())
			// jira allows to save the minimum of 1m
			secondsToSave := int(diff.Seconds()) - (int(diff.Seconds()) % 60)
			// do not perform update if we have the same values for time/description
			if w != nil && w.Comment == e.Description() && secondsToSave == w.TimeSpentSeconds {
				log.Printf("Time entry \"%s\" of %s for task %s is unchanged. Skipping update...\n", e.Description(), diff, t.Id)
				continue
			}
			record := &jira.WorklogRecord{
				Comment:          e.Description(),
				Started:          &st,
				TimeSpentSeconds: secondsToSave,
				Properties: []jira.EntityProperty{
					{
						Key: e.Source(),
						Value: map[string]interface{}{
							e.Source(): e.Id(),
						},
					},
				},
			}
			if w != nil {
				_, _, err = ji.api.Issue.UpdateWorklogRecord(t.Id, w.ID, record)
			} else {
				_, _, err = ji.api.Issue.AddWorklogRecord(t.Id, record)
			}
			log.Printf("Synchronized time entry \"%s\" of %s for task %s\n", e.Description(), diff, t.Id)
		}
	}
	return err
}
