package metrics

import (
	"sort"
	"time"

	"github.com/michael-quinton/four-keys/internal/domain"
)

// LeadTime is the time between a PR opening and merging. This is a
// simplified version of DORA's lead time (which is really commit-to-
// production), but we don't have deploy data to work with yet.
func LeadTime(pr domain.PullRequest) time.Duration {
	return pr.MergedAt.Sub(pr.CreatedAt)
}

func AverageLeadTime(prs []domain.PullRequest) time.Duration {
	if len(prs) == 0 {
		return 0
	}

	var total time.Duration
	for _, pr := range prs {
		total += LeadTime(pr)
	}

	return total / time.Duration(len(prs))
}

// MedianLeadTime is less sensitive than the average to a handful of
// slow outlier PRs dragging the number up.
func MedianLeadTime(prs []domain.PullRequest) time.Duration {
	if len(prs) == 0 {
		return 0
	}

	durations := make([]time.Duration, len(prs))
	for i, pr := range prs {
		durations[i] = LeadTime(pr)
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	mid := len(durations) / 2
	if len(durations)%2 == 1 {
		return durations[mid]
	}
	return (durations[mid-1] + durations[mid]) / 2
}
