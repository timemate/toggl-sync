# toggl-sync: Synchronize time reports with ease

![Latest GitHub release](https://img.shields.io/github/release/timemate/toggl-sync.svg)
![GitHub stars](https://img.shields.io/github/stars/timemate/toggl-sync.svg?label=github%20stars)
![Homebrew downloads](https://img.shields.io/homebrew/installs/dy/toggl-sync?label=macOS%20installs)
[![Go implementation (CI)](https://github.com/timemate/toggl-sync/workflows/Go%20implementation%20(CI)/badge.svg)](https://github.com/timemate/toggl-sync/actions?query=workflow%3A"Go+implementation+(CI)")

**toggl-sync** is a tiny CLI command/daemon for syncing [Toggl](https://toggl.com/) time entries with [Jira](https://www.atlassian.com/software/jira) written in [Go](https://go.dev/) and distributed with [HomeBrew](https://brew.sh/).

This is the first implementation of the bigger concept of time-tracking, syncing, and reporting among third-party platforms.

Find the full concept overview in [TimeMate docs](https://github.com/timemate).

## Installation

Install the app via [brew](https://brew.sh/) package manager.

```shell
brew tap timemate/tap
brew install toggl-sync
```

Re-install newer version of the library
```shell
brew update && \
  brew reinstall toggl-sync && \
  brew services restart toggl-sync
```

## Configuration

### Toggl
1. Login to your toggl account
2. Visit https://track.toggl.com/profile
3. Find API token

### Jira
1. Login to your atlassian account
2. Visit https://id.atlassian.com/manage-profile/security
3. Create new API token with a name: toggl-sync

### Config
Create a file in `~/.toggl-sync/config.yaml`
```yaml
tracker:
  - type: toggl
    token: "token-from-toggl"
#    projects:  # <-- use this section to filter projects
#      - 174942904

period:
  timeframe: 2w
  every: 1d

# places to sync time entries with
sync:
  - type: jira
    url: https://customer-host.atlassian.net
    login: login@email.com
    token: "token-from-jira"
    projects:
      - DO
      - DEV
```

- `timeframe` - look for timeframe in the past for new entries
- `every` - for `--service` mode, how often to repeat sync
- `projects` it's a list of project keys in Jira.

## Run the app

### Brew background service

```
brew services start toggl-sync
```
The program will start on system start up.

Find logs in here `tail -f $(brew --prefix)/var/log/toggl-sync/toggl-sync.log`

### One time run

```
toggl-sync sync -period 2w
```

Where `-period` specifies the period of time from current moment to sync. In this example `2w` - 2 past weeks.

### As a service

```
toggl-sync sync -period 1d --service
```

Where `--service` allows program to work as a daemon, it will repeatedly sync time entries every 1 day (`1d`).

## Development

### Build proto

```shell
protoc -I proto/ proto/*.proto --go_out=plugins=grpc:.
```
