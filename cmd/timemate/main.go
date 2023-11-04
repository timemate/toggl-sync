package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"godep.io/timemate/pkg/time_tracker"
	sharedTime "godep.io/timemate/pkg/time_tracker/shared"

	"github.com/hashicorp/go-plugin"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(time_tracker.TimeEntry{})
	gob.Register(time_tracker.Project{})
	gob.Register(time_tracker.Client{})
}

func main() {
	log.SetOutput(io.Discard)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: sharedTime.Handshake,
		Plugins:         sharedTime.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("TOGGL_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	raw, err := rpcClient.Dispense("toggl_grpc")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	timeTracker := raw.(time_tracker.ITimeTracker)
	entries, err := timeTracker.GetTimeEntries(time.Now().Add(-24*7*time.Hour), time.Now())
	log.Printf("Entries: %v\n", entries)
	log.Printf("Error: %v\n", err)
}
