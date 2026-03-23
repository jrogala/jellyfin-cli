package client

import "fmt"

// Library represents a Jellyfin library/view.
type Library struct {
	ID             string `json:"Id"`
	Name           string `json:"Name"`
	CollectionType string `json:"CollectionType"`
}

// Item represents a media item (movie, episode, song, etc).
type Item struct {
	ID              string            `json:"Id"`
	Name            string            `json:"Name"`
	Type            string            `json:"Type"`
	Path            string            `json:"Path"`
	ProductionYear  int               `json:"ProductionYear"`
	CommunityRating float64           `json:"CommunityRating"`
	OfficialRating  string            `json:"OfficialRating"`
	Overview        string            `json:"Overview"`
	SeriesName      string            `json:"SeriesName"`
	SeasonName      string            `json:"SeasonName"`
	IndexNumber     int               `json:"IndexNumber"`
	ParentIndexNum  int               `json:"ParentIndexNumber"`
	ProviderIds     map[string]string `json:"ProviderIds"`
	MediaType       string            `json:"MediaType"`
	RunTimeTicks    int64             `json:"RunTimeTicks"`
	HasSubtitles    bool              `json:"HasSubtitles"`
	Container       string            `json:"Container"`
}

// ItemsResult is the response from item queries.
type ItemsResult struct {
	Items            []Item `json:"Items"`
	TotalRecordCount int    `json:"TotalRecordCount"`
}

// SearchHint is a result from the search endpoint.
type SearchHint struct {
	ID             string `json:"Id"`
	Name           string `json:"Name"`
	Type           string `json:"Type"`
	ProductionYear int    `json:"ProductionYear"`
	MatchedTerm    string `json:"MatchedTerm"`
}

// SearchResult is the response from search.
type SearchResult struct {
	SearchHints        []SearchHint `json:"SearchHints"`
	TotalRecordCount   int          `json:"TotalRecordCount"`
}

// Session represents an active playback session.
type SessionInfo struct {
	ID         string `json:"Id"`
	DeviceName string `json:"DeviceName"`
	Client     string `json:"Client"`
	UserName   string `json:"UserName"`
	NowPlaying *Item  `json:"NowPlayingItem"`
}

// FormatRuntime formats ticks (100ns units) to human readable.
func FormatRuntime(ticks int64) string {
	if ticks == 0 {
		return ""
	}
	minutes := ticks / 600000000
	hours := minutes / 60
	mins := minutes % 60
	if hours > 0 {
		return fmt.Sprintf("%dh%02dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}
