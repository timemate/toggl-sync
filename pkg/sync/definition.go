package sync

import (
	"strings"

	"godep.io/timemate/pkg/trackers"
)

type Task struct {
	Id      string
	entries []trackers.TimeEntry
}

func GroupByTask(projects []interface{}, entries []trackers.TimeEntry) []Task {
	taskMap := map[string][]trackers.TimeEntry{}

	for _, e := range entries {
		var taskId string
		tags := e.Tags()
	ex:
		for _, pp := range projects {
			p := pp.(string)
			if tags != nil {
				for _, t := range tags {
					if strings.Index(t, p) == 0 {
						taskId = t
						break ex
					}
				}
			}
			if strings.Index(e.Description(), p) == 0 {
				if ss := strings.Split(e.Description(), " "); len(ss) > 0 {
					taskId = ss[0]
					break ex
				}
			}
		}

		if taskId == "" {
			continue
		}

		if taskMap[taskId] == nil {
			taskMap[taskId] = make([]trackers.TimeEntry, 0)
		}
		taskMap[taskId] = append(taskMap[taskId], e)
	}

	res := make([]Task, len(taskMap), len(taskMap))
	var i int
	for taskId, entries := range taskMap {
		res[i] = Task{
			Id:      taskId,
			entries: entries,
		}
		i++
	}

	return res
}
