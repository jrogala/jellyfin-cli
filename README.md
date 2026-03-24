# jellyfin-cli

CLI for Jellyfin media server management and browsing.

## Install

Download a binary from the [latest release](https://github.com/jrogala/jellyfin-cli/releases/latest), or install with Go:

```bash
go install github.com/jrogala/jellyfin-cli@latest
```

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
| `movies` | List all movies |
| `search` | Search across all media |
| `info` | Show detailed item info |
| `update` | Update item metadata |
| `identify` | Identify/match an item to an external provider |
| `refresh` | Refresh item metadata from providers |
| `scan` | Scan all libraries for changes |
| `sessions` | Show active playback sessions |
| `items` | List or filter items |

## Examples

```bash
$ jellyfin libraries
ID                                NAME    TYPE
507f1f77bcf86cd79943b735          Movies  movies
507f1f77bcf86cd79943b736          Shows   tvshows

$ jellyfin search "Interstellar"
Results: 2

ID                                TYPE   NAME          YEAR
507f1f77bcf86cd79943b735          Movie  Interstellar  2014
507f1f77bcf86cd79943b736          Movie  Interstellar  2024

$ jellyfin scan
Library scan started

$ jellyfin sessions
ID            CLIENT          DEVICE    USER   NOW PLAYING
abc123def4    Jellyfin Web    Chrome    admin  Interstellar
```

## JSON Output

All commands support `--json` for machine-readable output.
