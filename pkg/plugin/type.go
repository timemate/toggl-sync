package plugin

import "fmt"

const ( // iota is reset to 0
	PluginTypeTime = "time-tracker"
	PluginTypeTask = "task-tracker"
)

func PluginWithGrpc(name string) string {
	return fmt.Sprintf("%s_grpc", name)
}
