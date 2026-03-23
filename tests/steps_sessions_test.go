package tests

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/jrogala/jellyfin-cli/pkg/ops"
)

func (sc *scenarioCtx) serverHasActiveSessions(table *godog.Table) error {
	var sessions []map[string]any
	for i, row := range table.Rows {
		if i == 0 {
			continue
		}
		session := map[string]any{
			"Id":         row.Cells[0].Value,
			"DeviceName": row.Cells[1].Value,
			"Client":     row.Cells[2].Value,
			"UserName":   row.Cells[3].Value,
		}
		if row.Cells[4].Value != "" {
			session["NowPlayingItem"] = map[string]string{
				"Name": row.Cells[4].Value,
			}
		}
		sessions = append(sessions, session)
	}
	sc.mock.On("GET", "/Sessions", 200, sessions)
	return nil
}

func (sc *scenarioCtx) serverHasNoActiveSessions() error {
	sc.mock.On("GET", "/Sessions", 200, []any{})
	return nil
}

func (sc *scenarioCtx) iListSessions() error {
	sessions, err := ops.ListSessions(sc.client)
	sc.lastErr = err
	sc.sessions = sessions
	return nil
}

func (sc *scenarioCtx) iShouldGetNSessions(n int) error {
	sessions, ok := sc.sessions.([]ops.SessionEntry)
	if !ok {
		if n == 0 && sc.sessions == nil {
			return nil
		}
		return fmt.Errorf("no session results available")
	}
	if len(sessions) != n {
		return fmt.Errorf("expected %d sessions, got %d", n, len(sessions))
	}
	return nil
}

func (sc *scenarioCtx) sessionShouldShowNowPlaying(sessionID, expected string) error {
	sessions, ok := sc.sessions.([]ops.SessionEntry)
	if !ok {
		return fmt.Errorf("no session results available")
	}
	for _, s := range sessions {
		if s.ID == sessionID {
			if s.NowPlaying == expected {
				return nil
			}
			return fmt.Errorf("session %q now playing %q, expected %q", sessionID, s.NowPlaying, expected)
		}
	}
	return fmt.Errorf("session %q not found", sessionID)
}
