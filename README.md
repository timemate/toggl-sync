# toggl-sync
A tiny cli command/daemon for syncing toggl time entries with Jira

## Installation

Install the app via [brew](https://brew.sh/) package manager.

```shell
brew tap timemate/tap
brew install toggl-sync
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

- `projects` it's a list of project keys in Jira.

## Run the app

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
