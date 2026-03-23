package tests

import (
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	"github.com/jrogala/jellyfin-cli/client"
)

func initializeScenario(ctx *godog.ScenarioContext) {
	sc := newScenarioCtx(globalMock)

	ctx.Before(func(ctx context.Context, sc2 *godog.Scenario) (context.Context, error) {
		sc.mock.Reset()
		sc.lastErr = nil
		sc.authToken = ""
		sc.authUserID = ""
		sc.libraries = nil
		sc.itemsResult = nil
		sc.itemInfo = nil
		sc.searchResult = nil
		sc.sessions = nil
		sc.client = nil
		return ctx, nil
	})

	// --- Background ---
	ctx.Step(`^a running Jellyfin mock server$`, sc.aRunningJellyfinMockServer)
	ctx.Step(`^an authenticated client$`, sc.anAuthenticatedClient)

	// --- Auth ---
	ctx.Step(`^the server accepts credentials "([^"]*)" / "([^"]*)"$`, sc.serverAcceptsCredentials)
	ctx.Step(`^the server rejects all credentials$`, sc.serverRejectsAllCredentials)
	ctx.Step(`^I authenticate as "([^"]*)" with password "([^"]*)"$`, sc.iAuthenticateAs)
	ctx.Step(`^authentication should succeed$`, sc.authShouldSucceed)
	ctx.Step(`^authentication should fail with "([^"]*)"$`, sc.authShouldFailWith)
	ctx.Step(`^I should receive a token$`, sc.iShouldReceiveAToken)
	ctx.Step(`^I should receive a user ID$`, sc.iShouldReceiveAUserID)

	// --- Libraries ---
	ctx.Step(`^the server has libraries:$`, sc.serverHasLibraries)
	ctx.Step(`^the server has no libraries$`, sc.serverHasNoLibraries)
	ctx.Step(`^I list libraries$`, sc.iListLibraries)
	ctx.Step(`^I should get (\d+) libraries$`, sc.iShouldGetNLibraries)
	ctx.Step(`^library "([^"]*)" should have type "([^"]*)"$`, sc.libraryShouldHaveType)

	// --- Items ---
	ctx.Step(`^the server has movies:$`, sc.serverHasMovies)
	ctx.Step(`^the server has items of type "([^"]*)":$`, sc.serverHasItemsOfType)
	ctx.Step(`^the server has item "([^"]*)" with name "([^"]*)" and year (\d+)$`, sc.serverHasItem)
	ctx.Step(`^I list movies$`, sc.iListMovies)
	ctx.Step(`^I list items with type "([^"]*)"$`, sc.iListItemsWithType)
	ctx.Step(`^I get info for item "([^"]*)"$`, sc.iGetInfoForItem)
	ctx.Step(`^I should get (\d+) items$`, sc.iShouldGetNItems)
	ctx.Step(`^item "([^"]*)" should have year (\d+)$`, sc.itemShouldHaveYear)
	ctx.Step(`^the item name should be "([^"]*)"$`, sc.itemNameShouldBe)
	ctx.Step(`^the item year should be (\d+)$`, sc.itemYearShouldBe)

	// --- Search ---
	ctx.Step(`^the server returns search results for "([^"]*)":$`, sc.serverReturnsSearchResults)
	ctx.Step(`^the server returns no search results for "([^"]*)"$`, sc.serverReturnsNoSearchResults)
	ctx.Step(`^I search for "([^"]*)" with limit (\d+)$`, sc.iSearchFor)
	ctx.Step(`^I should get (\d+) search results$`, sc.iShouldGetNSearchResults)
	ctx.Step(`^search result "([^"]*)" should have type "([^"]*)"$`, sc.searchResultShouldHaveType)

	// --- Metadata ---
	ctx.Step(`^I update item "([^"]*)" with name "([^"]*)"$`, sc.iUpdateItemWithName)
	ctx.Step(`^the update should succeed$`, sc.operationShouldSucceed)
	ctx.Step(`^the server accepts identify for item "([^"]*)"$`, sc.serverAcceptsIdentify)
	ctx.Step(`^I identify item "([^"]*)" as "([^"]*)" with IMDB "([^"]*)"$`, sc.iIdentifyItem)
	ctx.Step(`^the identify should succeed$`, sc.operationShouldSucceed)
	ctx.Step(`^the server accepts refresh for item "([^"]*)"$`, sc.serverAcceptsRefresh)
	ctx.Step(`^I refresh metadata for item "([^"]*)"$`, sc.iRefreshMetadata)
	ctx.Step(`^the refresh should succeed$`, sc.operationShouldSucceed)
	ctx.Step(`^the server accepts library scan$`, sc.serverAcceptsLibraryScan)
	ctx.Step(`^I scan libraries$`, sc.iScanLibraries)
	ctx.Step(`^the scan should succeed$`, sc.operationShouldSucceed)

	// --- Sessions ---
	ctx.Step(`^the server has active sessions:$`, sc.serverHasActiveSessions)
	ctx.Step(`^the server has no active sessions$`, sc.serverHasNoActiveSessions)
	ctx.Step(`^I list sessions$`, sc.iListSessions)
	ctx.Step(`^I should get (\d+) sessions$`, sc.iShouldGetNSessions)
	ctx.Step(`^session "([^"]*)" should show "([^"]*)" as now playing$`, sc.sessionShouldShowNowPlaying)
}

// --- Background steps ---

func (sc *scenarioCtx) aRunningJellyfinMockServer() error {
	if sc.mock == nil {
		return fmt.Errorf("mock server is not running")
	}
	return nil
}

func (sc *scenarioCtx) anAuthenticatedClient() error {
	sc.client = client.NewWithBaseURL("test-token", "test-user-id", sc.mock.URL())
	return nil
}

// --- Common assertions ---

func (sc *scenarioCtx) operationShouldSucceed() error {
	if sc.lastErr != nil {
		return fmt.Errorf("expected success but got error: %v", sc.lastErr)
	}
	return nil
}

func (sc *scenarioCtx) authShouldFailWith(expected string) error {
	if sc.lastErr == nil {
		return fmt.Errorf("expected an error containing %q but got none", expected)
	}
	if !strings.Contains(sc.lastErr.Error(), expected) {
		return fmt.Errorf("expected error containing %q, got: %s", expected, sc.lastErr.Error())
	}
	return nil
}
