package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"godep.io/timemate/pkg/sync"
	"godep.io/timemate/pkg/trackers"
	"godep.io/timemate/pkg/utils"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.toggl-sync")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func updateEntries(trackerSvc trackers.TrackerApi, syncSvc sync.SyncApi, projects []interface{}, start time.Time, end time.Time) error {
	log.Printf("Update entries for duration %s\n", end.Sub(start))
	entries, err := trackerSvc.GetTimeEntries(start, end)
	if err != nil {
		return err
	}

	return syncSvc.Sync(sync.GroupByTask(projects, entries))
}

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Description = "Tiny service to sync time entries from toggl to jira"
	app.Name = "toggl-sync"
	app.Copyright = "TimeMate © 2021"

	app.Commands = []cli.Command{
		{
			Name:  "sync",
			Usage: "Sync time entries from timers with task trackers",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "period",
				},
				cli.BoolFlag{
					Name: "service",
				},
			},
			Action: func(c *cli.Context) error {
				service := c.Bool("service")
				t := viper.GetString("period.timeframe")
				p := c.String("period")
				if p != "" {
					t = p
				}
				lookupTimeframe, err := utils.ParseDuration(t)
				if err != nil {
					log.Fatal(err)
				}
				e := viper.GetString("period.every")
				if e == "" {
					e = t
				}
				every, err := utils.ParseDuration(e)
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Sync started with config period: %s, every: %s", lookupTimeframe, every)

				end := time.Now()
				start := end.Add(-lookupTimeframe)

				tr := viper.Get("tracker").([]interface{})
				trackerSvc, err := trackers.NewTogglTracker((tr[0]).(map[interface{}]interface{}))
				if err != nil {
					log.Fatal(err)
				}

				sn := viper.Get("sync").([]interface{})
				jiraConf := (sn[0]).(map[interface{}]interface{})
				syncSvc, err := sync.NewJiraSync(jiraConf)
				if err != nil {
					log.Fatal(err)
				}

				err = updateEntries(trackerSvc, syncSvc, (jiraConf["projects"]).([]interface{}), start, end)
				if err != nil {
					log.Fatal(err)
				}
				lastRun := time.Now().Format("2006-01-02T15:04:05-0700")
				if service {
					ticker := time.NewTicker(time.Second)
					quit := make(chan error, 1)
					go func() {
						for {
							select {
							case <-ticker.C:
								nowFormatted := time.Now().Format("2006-01-02T15:04:05-0700")
								now, err := time.Parse("2006-01-02T15:04:05-0700", nowFormatted)
								last, err := time.Parse("2006-01-02T15:04:05-0700", lastRun)
								if now.Sub(last).Seconds() > every.Seconds() {
									err = updateEntries(trackerSvc, syncSvc, (jiraConf["projects"]).([]interface{}), start, end)
									lastRun = time.Now().Format("2006-01-02T15:04:05-0700")
								}
								if err != nil {
									quit <- err
								}
							}
						}
					}()

					err = <-quit
					log.Printf("%v", err)
					log.Println("Quiting process")
					ticker.Stop()
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
