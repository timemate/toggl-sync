package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	plugConfig "godep.io/timemate/pkg/config"
	pkgPlugin "godep.io/timemate/pkg/plugin"
	"godep.io/timemate/pkg/task_tracker"
	"godep.io/timemate/pkg/time_tracker"

	"github.com/urfave/cli"
)

var (
	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

func init() {
	gob.Register(time.Time{})
	gob.Register(time_tracker.TimeEntry{})
	gob.Register(time_tracker.Project{})
	gob.Register(time_tracker.Client{})

	gob.Register(task_tracker.Project{})
	gob.Register(task_tracker.Task{})
}

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s (%s; %s)", version, commit, date)
	app.Description = "Tiny service to sync time entries from toggl to jira"
	app.Name = "timemate"
	app.Copyright = "TimeMate © 2021-2023"
	//log.SetOutput(io.Discard)

	config, err := plugConfig.ReadConfig()
	if err != nil {
		panic(err)
	}

	var timeTracker time_tracker.ITimeTracker
	if impl, togglClient, err := pkgPlugin.GetGRPCClient(config, os.Getenv("TOGGL_PLUGIN"), "toggl"); err == nil {
		defer togglClient.Kill()
		timeTracker = impl.(time_tracker.ITimeTracker)
	} else {
		panic("Unexpected error")
	}

	var taskTracker task_tracker.ITaskTracker
	if impl, jiraClient, err := pkgPlugin.GetGRPCClient(config, os.Getenv("JIRA_PLUGIN"), "jira"); err == nil {
		defer jiraClient.Kill()
		taskTracker = impl.(task_tracker.ITaskTracker)
	} else {
		panic("Unexpected error")
	}

	entries, err := timeTracker.GetTimeEntries(time.Now().Add(-24*7*time.Hour), time.Now())
	log.Printf("Entries: %v\n", len(entries))
	log.Printf("Error: %v\n", err)

	var v = make([]task_tracker.ITask, 0)
	err = taskTracker.UpdateTasks(v)
	log.Printf("Error: %v\n", err)
}
