package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is a thin wrapper around http.Client for the GitHub REST API.
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient builds a client. token can be empty but GitHub's
// unauthenticated rate limit (60/hr) makes that painful in practice.
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		token:      token,
	}
}

// FetchClosedPullRequests gets the most recently updated closed PRs for
// owner/repo. "Closed" covers merged and unmerged - filter on MergedAt
// if you only want the merged ones.
func (c *Client) FetchClosedPullRequests(ctx context.Context, owner, repo string, perPage int) ([]PullRequestResponse, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls?state=closed&per_page=%d&sort=updated&direction=desc",
		owner, repo, perPage,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d for %s", resp.StatusCode, url)
	}

	var prs []PullRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&prs); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return prs, nil
}
