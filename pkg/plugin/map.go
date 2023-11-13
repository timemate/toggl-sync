package plugin

import (
	"errors"

	"godep.io/timemate/pkg/config"
	"godep.io/timemate/pkg/task_tracker"
	sharedTask "godep.io/timemate/pkg/task_tracker/shared"
	"godep.io/timemate/pkg/time_tracker"
	sharedTime "godep.io/timemate/pkg/time_tracker/shared"

	"github.com/hashicorp/go-plugin"
)

func GetPluginMap(config config.IConfig, pluginName string, isGRPC bool, args ...interface{}) (map[string]plugin.Plugin, error) {
	plugins := make(map[string]plugin.Plugin, 0)

	conf := config.FindPlugin(pluginName)
	if conf == nil {
		return nil, errors.New("no plugin configuration found")
	}
	switch conf.Type {
	case PluginTypeTime:
		var impl time_tracker.ITimeTracker
		if args != nil && len(args) > 0 {
			impl = args[0].(time_tracker.ITimeTracker)
		}
		if isGRPC {
			plugins[PluginWithGrpc(conf.Name)] = &sharedTime.TimeTrackerGRPCPlugin{
				Impl: impl,
			}
		} else {
			plugins[conf.Name] = &sharedTime.TimeTrackerPlugin{
				Impl: impl,
			}
		}
	case PluginTypeTask:
		var impl task_tracker.ITaskTracker
		if args != nil && len(args) > 0 {
			impl = args[0].(task_tracker.ITaskTracker)
		}
		if isGRPC {
			plugins[PluginWithGrpc(conf.Name)] = &sharedTask.TaskTrackerGRPCPlugin{
				Impl: impl,
			}
		} else {
			plugins[conf.Name] = &sharedTask.TaskTrackerPlugin{
				Impl: impl,
			}
		}
	}
	return plugins, nil
}
