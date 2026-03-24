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
ID                                NAME          TYPE
f137a2dd21bbc1b99aa5c0f6bf02a805  Movies        movies
7e64e319657a9516ec78490da03edccb  Music         music
b26028d1e9bd4b29c82f1bfc13f523d2  TV en direct  livetv

$ jellyfin search "Dexter"
Results: 16

ID                                TYPE   NAME                              YEAR
f92febffef45ee1015ec6878292413aa  Movie  Dexter S01E11 - Truth Be Told     2006
58c757e688f2bb832f34343acc919cf0  Movie  Dexter S01E12 - Born Free         2006
816babf4d9a5f41778d33f096cf21803  Movie  Dexter S02E01 - It's Alive!       2007

$ jellyfin sessions
ID                                CLIENT        DEVICE  USER    NOW PLAYING
6e75dad0e95c74a50807ab20de3aff04  jellyfin-cli  cli     claude
```

## JSON Output

All commands support `--json` for machine-readable output.
