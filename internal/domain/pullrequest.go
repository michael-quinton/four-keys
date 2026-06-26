// Package domain holds the app's own PullRequest model - smaller than
// GitHub's API response, and what metrics/report actually work with.
package domain

import (
	"fmt"
	"time"

	"github.com/michael-quinton/four-keys/internal/github"
)

type PullRequest struct {
	Number         int
	Title          string
	CreatedAt      time.Time
	MergedAt       time.Time
	MergeCommitSHA string
	Author         string
	Repository     string
}

// FromGitHub converts a raw API response into a PullRequest. Returns
// ok=false for PRs that were closed without being merged - we only care
// about merged work here.
func FromGitHub(r github.PullRequestResponse) (pr PullRequest, ok bool, err error) {
	if r.MergedAt == nil {
		return PullRequest{}, false, nil
	}

	createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		return PullRequest{}, false, fmt.Errorf("parsing created_at for PR #%d: %w", r.Number, err)
	}

	mergedAt, err := time.Parse(time.RFC3339, *r.MergedAt)
	if err != nil {
		return PullRequest{}, false, fmt.Errorf("parsing merged_at for PR #%d: %w", r.Number, err)
	}

	pr = PullRequest{
		Number:         r.Number,
		Title:          r.Title,
		CreatedAt:      createdAt,
		MergedAt:       mergedAt,
		MergeCommitSHA: r.MergeCommitSHA,
		Author:         r.User.Login,
		Repository:     r.Base.Repo.FullName,
	}

	return pr, true, nil
}
