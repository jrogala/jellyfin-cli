# jellyfin-cli

CLI for Jellyfin media server management and browsing.

## Setup

Set `JELLYFIN_URL` env var or use config file, then authenticate:

```bash
export JELLYFIN_URL=http://jellyfin.local:8096
jellyfin login -u admin -p password
```

## Commands

| Command | Description |
|---|---|
| `login` | Authenticate with Jellyfin server |
| `libraries` | List media libraries |
| `movies` | List movies in a library |
| `search` | Search across all media |
| `info` | Show detailed item info |
| `update` | Update item metadata |
| `identify` | Identify/match an item to an external provider |
| `refresh` | Refresh item metadata from providers |
| `scan` | Trigger a library scan |
| `sessions` | List active sessions |
| `items` | List or filter items |

## Examples

```bash
# Login and browse libraries
jellyfin login -u admin -p password
jellyfin libraries

# Search for a movie
jellyfin search "Interstellar"

# Trigger a full library scan
jellyfin scan --library "Movies"

# Show active streaming sessions
jellyfin sessions
```

## JSON Output

All commands support `--json` for machine-readable output.
