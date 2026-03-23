package tests

import "github.com/jrogala/jellyfin-cli/client"

// scenarioCtx holds per-scenario state.
type scenarioCtx struct {
	mock   *MockServer
	client *client.Client

	// results from the latest "When" step
	lastErr       error
	authToken     string
	authUserID    string
	libraries     any
	itemsResult   any
	itemInfo      any
	searchResult  any
	sessions      any
}

func newScenarioCtx(mock *MockServer) *scenarioCtx {
	return &scenarioCtx{
		mock: mock,
	}
}
