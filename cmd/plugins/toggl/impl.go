package main

import (
	"strconv"
	"time"

	"godep.io/timemate/pkg/config"
	"godep.io/timemate/pkg/time_tracker"
	"godep.io/timemate/pkg/utils"

	"github.com/andreaskoch/togglapi"
	"github.com/andreaskoch/togglapi/model"
)

type TogglTracker struct {
	model.TimeEntryAPI
	model.ProjectAPI
	model.ClientAPI
	projects []string
}

const baseHost = "https://api.track.toggl.com/api/v8"

func NewTogglTracker(config config.PluginConfig) (*TogglTracker, error) {
	token := config.Config["token"]
	var projects []string
	//p, ok := (params["projects"]).([]interface{})
	//if ok {
	//	projects = make([]string, 0)
	//	for _, pp := range p {
	//		switch v := pp.(type) {
	//		case int:
	//			projects = append(projects, strconv.Itoa(v))
	//		case string:
	//			projects = append(projects, v)
	//		}
	//	}
	//}
	return &TogglTracker{
		TimeEntryAPI: togglapi.NewTimeEntryAPI(baseHost, token),
		ProjectAPI:   togglapi.NewProjectAPI(baseHost, token),
		ClientAPI:    togglapi.NewClientAPI(baseHost, token),
		projects:     projects,
	}, nil
}

func getProjectClientByTimeEntry(clients map[string]model.Client, projects map[string]model.Project, entry model.TimeEntry) (client time_tracker.IClient, project time_tracker.IProject) {
	p, ok := projects[strconv.Itoa(entry.Pid)]
	if !ok {
		return client, project
	}
	project = &time_tracker.Project{
		Id:   strconv.Itoa(p.ID),
		Name: p.Name,
	}
	c, ok := clients[strconv.Itoa(p.ClientID)]
	if !ok {
		return client, project
	}
	client = &time_tracker.Client{
		Id:   strconv.Itoa(c.ID),
		Name: c.Name,
	}
	return client, project
}

func (tg *TogglTracker) getClientsAndProjects() (clientsMap map[string]model.Client, projectsMap map[string]model.Project, err error) {
	clientsMap = make(map[string]model.Client)
	projectsMap = make(map[string]model.Project)
	clients, err := tg.ClientAPI.GetClients()
	if err != nil {
		return clientsMap, projectsMap, err
	}
	workspacesRefs := make([]int, 0)
	for _, c := range clients {
		clientsMap[strconv.Itoa(c.ID)] = c
		if !utils.InArray[int](workspacesRefs, c.WorkspaceID) {
			workspacesRefs = append(workspacesRefs, c.WorkspaceID)
		}
	}

	for _, ref := range workspacesRefs {
		projects, err := tg.ProjectAPI.GetProjects(ref)
		for _, p := range projects {
			projectsMap[strconv.Itoa(p.ID)] = p
		}
		if err != nil {
			return clientsMap, projectsMap, err
		}
	}
	return clientsMap, projectsMap, err
}

func (tg *TogglTracker) GetTimeEntries(start time.Time, end time.Time) ([]time_tracker.ITimeEntry, error) {
	entries, err := tg.TimeEntryAPI.GetTimeEntries(start, end)
	if err != nil {
		return nil, err
	}

	clientsMap, projectsMap, err := tg.getClientsAndProjects()
	if err != nil {
		return nil, err
	}

	r := make([]time_tracker.ITimeEntry, 0)
	for _, e := range entries {
		client, project := getProjectClientByTimeEntry(clientsMap, projectsMap, e)
		if len(tg.projects) > 0 && !utils.InArray[string](tg.projects, project.GetId()) {
			continue
		}
		r = append(r, time_tracker.TimeEntry{
			Id:          strconv.Itoa(e.ID),
			Description: e.Description,
			Start:       e.Start,
			Stop:        e.Stop,
			Tags:        e.Tags,
			Source:      "toggl",
			Client:      client,
			Project:     project,
		})
	}
	return r, nil
}
