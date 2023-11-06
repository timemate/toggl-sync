package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	plugConfig "godep.io/timemate/pkg/config"
	pkgPlugin "godep.io/timemate/pkg/plugin"
	"godep.io/timemate/pkg/time_tracker"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(time_tracker.TimeEntry{})
	gob.Register(time_tracker.Project{})
	gob.Register(time_tracker.Client{})
}

func main() {
	log.SetOutput(io.Discard)

	config, err := plugConfig.ReadConfig()
	if err != nil {
		panic(err)
	}

	timeTracker, client, err := pkgPlugin.GetGRPCClient[time_tracker.ITimeTracker](config, os.Getenv("TOGGL_PLUGIN"), "toggl")
	defer client.Kill()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	entries, err := timeTracker.GetTimeEntries(time.Now().Add(-24*7*time.Hour), time.Now())
	log.Printf("Entries: %v\n", entries)
	log.Printf("Error: %v\n", err)
}
