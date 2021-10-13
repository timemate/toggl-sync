# toggl-sync
A tool for syncing toggl time entries with Jira

## Toggl
1. Login to your toggl account
2. Visit https://track.toggl.com/profile
3. Find API token

## Jira
1. Login to your atlassian account
2. Visit https://id.atlassian.com/manage-profile/security
3. Create new API token with a name: toggl-sync

Create a file in `~/toggl-sync/config.yaml`
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
      - DEV
```
