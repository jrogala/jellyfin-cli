package ops

import "github.com/jrogala/jellyfin-cli/client"

// SessionEntry represents an active playback session.
type SessionEntry struct {
	ID         string `json:"id"`
	DeviceName string `json:"device_name"`
	Client     string `json:"client"`
	UserName   string `json:"user_name"`
	NowPlaying string `json:"now_playing,omitempty"`
}

// ListSessions returns all active playback sessions.
func ListSessions(c *client.Client) ([]SessionEntry, error) {
	sessions, err := c.GetSessions()
	if err != nil {
		return nil, err
	}
	entries := make([]SessionEntry, 0, len(sessions))
	for _, s := range sessions {
		np := ""
		if s.NowPlaying != nil {
			np = s.NowPlaying.Name
		}
		entries = append(entries, SessionEntry{
			ID:         s.ID,
			DeviceName: s.DeviceName,
			Client:     s.Client,
			UserName:   s.UserName,
			NowPlaying: np,
		})
	}
	return entries, nil
}
