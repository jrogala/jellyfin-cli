package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client communicates with the Jellyfin API.
type Client struct {
	baseURL    string
	token      string
	userID     string
	httpClient *http.Client
}

// New creates a new Jellyfin client.
func New(baseURL, token, userID string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		userID:     userID,
		httpClient: http.DefaultClient,
	}
}

// NewWithBaseURL creates a client pointing at a custom base URL (for testing).
func NewWithBaseURL(token, userID, baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		userID:     userID,
		httpClient: http.DefaultClient,
	}
}

// SetHTTPClient overrides the HTTP client used for requests.
func (c *Client) SetHTTPClient(hc *http.Client) {
	c.httpClient = hc
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

const authHeader = `MediaBrowser Client="jellyfin-cli", Device="cli", DeviceId="jellyfin-cli-1", Version="1.0"`

func (c *Client) do(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("X-Emby-Authorization", authHeader+`, Token="`+c.token+`"`)
	} else {
		req.Header.Set("X-Emby-Authorization", authHeader)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jellyfin error %d: %s", resp.StatusCode, string(data))
	}
	return resp, nil
}

func (c *Client) doJSON(method, path string, body io.Reader, result any) error {
	resp, err := c.do(method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if result != nil && resp.StatusCode != http.StatusNoContent {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// Authenticate logs in and returns token + user ID.
func Authenticate(baseURL, username, password string) (token, userID string, err error) {
	return AuthenticateWithClient(baseURL, username, password, http.DefaultClient)
}

// AuthenticateWithClient logs in using a custom HTTP client (for testing).
func AuthenticateWithClient(baseURL, username, password string, hc *http.Client) (token, userID string, err error) {
	body := fmt.Sprintf(`{"Username":%q,"Pw":%q}`, username, password)
	req, err := http.NewRequest("POST", baseURL+"/Users/AuthenticateByName", strings.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Emby-Authorization", authHeader)

	resp, err := hc.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("authentication failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("authentication failed (%d): %s", resp.StatusCode, string(data))
	}

	var result struct {
		AccessToken string `json:"AccessToken"`
		User        struct {
			ID string `json:"Id"`
		} `json:"User"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}
	return result.AccessToken, result.User.ID, nil
}

// --- Libraries ---

func (c *Client) GetLibraries() ([]Library, error) {
	var result struct {
		Items []Library `json:"Items"`
	}
	if err := c.doJSON("GET", "/Users/"+c.userID+"/Views", nil, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// --- Items ---

func (c *Client) GetItems(parentID string, itemType string, sortBy string, limit int) (*ItemsResult, error) {
	params := url.Values{
		"Recursive": {"true"},
		"Fields":    {"Path,Overview,CommunityRating,ProductionYear,ProviderIds"},
		"SortBy":    {sortBy},
		"SortOrder": {"Ascending"},
	}
	if parentID != "" {
		params.Set("ParentId", parentID)
	}
	if itemType != "" {
		params.Set("IncludeItemTypes", itemType)
	}
	if limit > 0 {
		params.Set("Limit", strconv.Itoa(limit))
	}

	var result ItemsResult
	if err := c.doJSON("GET", "/Users/"+c.userID+"/Items?"+params.Encode(), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetItem(itemID string) (*Item, error) {
	var item Item
	if err := c.doJSON("GET", "/Users/"+c.userID+"/Items/"+itemID, nil, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// --- Search ---

func (c *Client) Search(query string, limit int) (*SearchResult, error) {
	params := url.Values{
		"searchTerm": {query},
		"Limit":      {strconv.Itoa(limit)},
	}
	var result SearchResult
	if err := c.doJSON("GET", "/Search/Hints?"+params.Encode(), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- Metadata ---

func (c *Client) UpdateItem(itemID string, updates map[string]any) error {
	// First get the full item to preserve existing fields
	resp, err := c.do("GET", "/Items/"+itemID, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var item map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return err
	}

	// Apply updates
	for k, v := range updates {
		item[k] = v
	}

	data, _ := json.Marshal(item)
	resp2, err := c.do("POST", "/Items/"+itemID, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	resp2.Body.Close()
	return nil
}

func (c *Client) RefreshMetadata(itemID string) error {
	resp, err := c.do("POST", "/Items/"+itemID+"/Refresh?ReplaceAllMetadata=true&ReplaceAllImages=true", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) IdentifyItem(itemID, name string, year int, providerIDs map[string]string) error {
	body := map[string]any{
		"SearchInfo": map[string]any{
			"Name":        name,
			"Year":        year,
			"ProviderIds": providerIDs,
		},
	}
	data, _ := json.Marshal(body)
	resp, err := c.do("POST", "/Items/"+itemID+"/Identify", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// StreamURL returns a direct stream URL for an item.
func (c *Client) StreamURL(itemID string) string {
	return c.baseURL + "/Videos/" + itemID + "/stream?static=true&api_key=" + c.token
}

// --- Sessions / Playback ---

func (c *Client) GetSessions() ([]SessionInfo, error) {
	var sessions []SessionInfo
	if err := c.doJSON("GET", "/Sessions", nil, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// --- Library Management ---

func (c *Client) ScanLibrary() error {
	resp, err := c.do("POST", "/Library/Refresh", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
