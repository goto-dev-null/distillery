package forgejo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/distillery/pkg/common"
)

const defaultBaseURL = "https://codeberg.org/api/v1"

// NewClient creates a new Forgejo client using the provided http.Client.
// The base URL defaults to Codeberg. Call SetBaseURL to override for other instances.
func NewClient(client *http.Client) *Client {
	return &Client{
		client:  client,
		baseURL: defaultBaseURL,
	}
}

// Client is a minimal HTTP client for the Forgejo/Gitea v1 REST API.
type Client struct {
	client  *http.Client
	baseURL string
	token   string
}

// SetToken sets the API token sent as "Authorization: token <TOKEN>" on every request.
func (c *Client) SetToken(token string) {
	c.token = token
}

// SetBaseURL overrides the default base URL (https://codeberg.org/api/v1).
// The value should include the full API path, e.g. "https://git.example.com/api/v1".
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// GetToken returns the currently configured API token.
func (c *Client) GetToken() string {
	return c.token
}

// GetClient returns the underlying *http.Client.
func (c *Client) GetClient() *http.Client {
	return c.client
}

// doRequest constructs, authenticates, and executes a GET request, returning
// the response. The caller is responsible for closing the response body.
// A non-200 status is treated as an error.
func (c *Client) doRequest(ctx context.Context, url string) (*http.Response, error) {
	logrus.Tracef("GET %s", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", common.NAME, common.AppVersion))
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	return resp, nil
}

// ListReleases returns all releases for the given owner/repo, paging through
// results until the server returns fewer entries than the page size.
// Releases are returned in the order the API provides them (newest first).
func (c *Client) ListReleases(ctx context.Context, owner, repo string) ([]*Release, error) {
	var all []*Release

	const pageSize = 50
	for page := 1; ; page++ {
		url := fmt.Sprintf("%s/repos/%s/%s/releases?limit=%d&page=%d",
			c.baseURL, owner, repo, pageSize, page)

		resp, err := c.doRequest(ctx, url)
		if err != nil {
			return nil, err
		}

		var releases []*Release
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
		_ = resp.Body.Close()

		all = append(all, releases...)

		if len(releases) < pageSize {
			break
		}
	}

	return all, nil
}

// GetLatestRelease returns the most recent non-draft, non-pre-release for the
// given owner/repo. Returns an error (containing "404") if no releases exist.
func (c *Client) GetLatestRelease(ctx context.Context, owner, repo string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", c.baseURL, owner, repo)

	resp, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// GetRelease returns the release matching the given tag (e.g. "v1.2.3").
func (c *Client) GetRelease(ctx context.Context, owner, repo, tag string) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", c.baseURL, owner, repo, tag)

	resp, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}
