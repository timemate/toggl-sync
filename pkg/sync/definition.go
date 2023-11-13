package sync

import (
	"strings"

	"godep.io/timemate/pkg/time_tracker"
)

type Task struct {
	Id      string
	entries []time_tracker.ITimeEntry
}

func GroupByTask(projects []interface{}, entries []time_tracker.ITimeEntry) []Task {
	taskMap := map[string][]time_tracker.ITimeEntry{}

	for _, e := range entries {
		var taskId string
		tags := e.GetTags()
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
			if strings.Index(e.GetDescription(), p) == 0 {
				if ss := strings.Split(e.GetDescription(), " "); len(ss) > 0 {
					taskId = ss[0]
					break ex
				}
			}
		}

		if taskId == "" {
			continue
		}

		if taskMap[taskId] == nil {
			taskMap[taskId] = make([]time_tracker.ITimeEntry, 0)
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
