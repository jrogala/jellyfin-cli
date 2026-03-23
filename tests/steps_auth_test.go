package tests

import (
	"fmt"

	"github.com/jrogala/jellyfin-cli/pkg/ops"
)

func (sc *scenarioCtx) serverAcceptsCredentials(username, password string) error {
	sc.mock.On("POST", "/Users/AuthenticateByName", 200, map[string]any{
		"AccessToken": "mock-token-abc",
		"User":        map[string]string{"Id": "mock-user-123"},
	})
	return nil
}

func (sc *scenarioCtx) serverRejectsAllCredentials() error {
	sc.mock.On("POST", "/Users/AuthenticateByName", 401, map[string]string{
		"error": "Invalid credentials",
	})
	return nil
}

func (sc *scenarioCtx) iAuthenticateAs(username, password string) error {
	result, err := ops.Authenticate(sc.mock.URL(), username, password)
	sc.lastErr = err
	if result != nil {
		sc.authToken = result.Token
		sc.authUserID = result.UserID
	}
	return nil
}

func (sc *scenarioCtx) authShouldSucceed() error {
	if sc.lastErr != nil {
		return fmt.Errorf("expected auth to succeed, got: %v", sc.lastErr)
	}
	return nil
}

func (sc *scenarioCtx) iShouldReceiveAToken() error {
	if sc.authToken == "" {
		return fmt.Errorf("expected a token but got empty string")
	}
	return nil
}

func (sc *scenarioCtx) iShouldReceiveAUserID() error {
	if sc.authUserID == "" {
		return fmt.Errorf("expected a user ID but got empty string")
	}
	return nil
}
