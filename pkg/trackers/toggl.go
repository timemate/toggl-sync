package trackers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"godep.io/timemate/pkg/utils"
)

const baseHost = "https://api.track.toggl.com/api/v9"

type WorkspaceImpl struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ClientImpl struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	WorkspaceID int    `json:"workspace_id"`
}

type ProjectImpl struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ClientID    int    `json:"client_id"`
	WorkspaceID int    `json:"workspace_id"`
}

type TimeEntryImpl struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	Start       string   `json:"start"`
	Stop        string   `json:"stop"`
	Tags        []string `json:"tags"`
	Pid         int      `json:"project_id"`
}

type TogglTracker struct {
	token    string
	projects []string
}

func NewTogglTracker(params map[interface{}]interface{}) (*TogglTracker, error) {
	token := (params["token"]).(string)
	var projects []string
	p, ok := (params["projects"]).([]interface{})
	if ok {
		projects = make([]string, 0)
		for _, pp := range p {
			switch v := pp.(type) {
			case int:
				projects = append(projects, strconv.Itoa(v))
			case string:
				projects = append(projects, v)
			}
		}
	}
	return &TogglTracker{
		token:    token,
		projects: projects,
	}, nil
}

func (tg *TogglTracker) getRequest(endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseHost+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(tg.token, "api_token")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var body []byte
		if resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
		}
		return body, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (tg *TogglTracker) getClientsAndProjects() (clientsMap map[string]Client, projectsMap map[string]Project, err error) {
	clientsMap = make(map[string]Client)
	projectsMap = make(map[string]Project)
	// Get workspaces
	wsData, err := tg.getRequest("/workspaces")
	if err != nil {
		return clientsMap, projectsMap, err
	}
	var workspaces []WorkspaceImpl
	err = json.Unmarshal(wsData, &workspaces)
	if err != nil {
		return clientsMap, projectsMap, err
	}
	// Get clients and projects for each workspace
	for _, ws := range workspaces {
		clData, err := tg.getRequest(fmt.Sprintf("/workspaces/%d/clients", ws.ID))
		if err != nil {
			return clientsMap, projectsMap, err
		}
		var clients []ClientImpl
		err = json.Unmarshal(clData, &clients)
		if err != nil {
			return clientsMap, projectsMap, err
		}
		for _, c := range clients {
			idStr := strconv.Itoa(c.ID)
			clientsMap[idStr] = trackedTimeClient{
				id:   idStr,
				name: c.Name,
			}
		}
		prData, err := tg.getRequest(fmt.Sprintf("/workspaces/%d/projects", ws.ID))
		if err != nil {
			return clientsMap, projectsMap, err
		}
		var projects []ProjectImpl
		err = json.Unmarshal(prData, &projects)
		if err != nil {
			return clientsMap, projectsMap, err
		}
		for _, p := range projects {
			idStr := strconv.Itoa(p.ID)
			projectsMap[idStr] = trackedTimeProject{
				id:       idStr,
				name:     p.Name,
				clientId: strconv.Itoa(p.ClientID),
			}
		}
	}
	return clientsMap, projectsMap, nil
}

func getProjectClientByTimeEntry(clients map[string]Client, projects map[string]Project, entry TimeEntryImpl) (client Client, project Project) {
	p, ok := projects[strconv.Itoa(entry.Pid)]
	if !ok {
		return client, project
	}
	c, ok := clients[p.ClientId()]
	if !ok {
		return client, project
	}
	return c, p
}

func (tg *TogglTracker) GetTimeEntries(start time.Time, end time.Time) ([]TimeEntry, error) {
	// Get time entries
	startStr := url.QueryEscape(start.Format(time.RFC3339))
	endStr := url.QueryEscape(end.Format(time.RFC3339))
	endpoint := fmt.Sprintf("/me/time_entries?start_date=%s&end_date=%s", startStr, endStr)
	teData, err := tg.getRequest(endpoint)
	if err != nil {
		return nil, err
	}
	var entries []TimeEntryImpl
	err = json.Unmarshal(teData, &entries)
	if err != nil {
		return nil, err
	}
	clientsMap, projectsMap, err := tg.getClientsAndProjects()
	if err != nil {
		return nil, err
	}
	r := make([]TimeEntry, 0)
	for _, e := range entries {
		client, project := getProjectClientByTimeEntry(clientsMap, projectsMap, e)
		if len(tg.projects) > 0 && !utils.InArray[string](tg.projects, project.Id()) {
			continue
		}
		parsedStart, _ := time.Parse(time.RFC3339, e.Start)
		parsedStop, _ := time.Parse(time.RFC3339, e.Stop)
		r = append(r, trackedTime{
			id:          strconv.Itoa(e.ID),
			description: e.Description,
			start:       parsedStart,
			stop:        parsedStop,
			tags:        e.Tags,
			source:      "toggl",
			client:      client,
			project:     project,
		})
	}
	return r, nil
}
