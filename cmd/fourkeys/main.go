// fourkeys fetches recently closed PRs for a repo and prints Lead Time
// for Changes (one of the four DORA metrics). The other three are just
// placeholders for now.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/michael-quinton/four-keys/internal/domain"
	"github.com/michael-quinton/four-keys/internal/github"
	"github.com/michael-quinton/four-keys/internal/metrics"
	"github.com/michael-quinton/four-keys/internal/report"
)

const perPage = 30

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	owner, repo, err := parseTarget(os.Args)
	if err != nil {
		return err
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "warning: GITHUB_TOKEN not set, you'll hit GitHub's rate limit fast")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := github.NewClient(token)

	raw, err := client.FetchClosedPullRequests(ctx, owner, repo, perPage)
	if err != nil {
		return fmt.Errorf("fetching pull requests: %w", err)
	}

	prs, err := mergedOnly(raw)
	if err != nil {
		return err
	}

	rpt := report.Report{
		Repository:      owner + "/" + repo,
		AnalysedCount:   len(prs),
		AverageLeadTime: metrics.AverageLeadTime(prs),
		MedianLeadTime:  metrics.MedianLeadTime(prs),
	}

	fmt.Print(rpt.String())
	return nil
}

// drops PRs that were closed without being merged
func mergedOnly(raw []github.PullRequestResponse) ([]domain.PullRequest, error) {
	prs := make([]domain.PullRequest, 0, len(raw))

	for _, r := range raw {
		pr, ok, err := domain.FromGitHub(r)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

// parseTarget expects args[1] to be "owner/repo", e.g.
// go run ./cmd/fourkeys kubernetes/kubernetes
func parseTarget(args []string) (owner, repo string, err error) {
	target := "kubernetes/kubernetes"
	if len(args) > 1 {
		target = args[1]
	}

	parts := strings.SplitN(target, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo %q, expected owner/repo", target)
	}

	return parts[0], parts[1], nil
}
