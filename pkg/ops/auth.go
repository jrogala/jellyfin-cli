package ops

import "github.com/jrogala/jellyfin-cli/client"

// AuthResult holds the result of an authentication attempt.
type AuthResult struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

// Authenticate logs in to Jellyfin and returns the access token and user ID.
func Authenticate(baseURL, username, password string) (*AuthResult, error) {
	token, userID, err := client.Authenticate(baseURL, username, password)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		Token:  token,
		UserID: userID,
	}, nil
}
