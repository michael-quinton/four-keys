package report

import (
	"fmt"
	"strings"
	"time"
)

type Report struct {
	Repository      string
	AnalysedCount   int
	AverageLeadTime time.Duration
	MedianLeadTime  time.Duration
}

func (r Report) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Repository: %s\n\n", r.Repository)
	fmt.Fprintf(&b, "Merged PRs analysed: %d\n", r.AnalysedCount)
	fmt.Fprintf(&b, "Average lead time: %s\n", FormatDuration(r.AverageLeadTime))
	fmt.Fprintf(&b, "Median lead time: %s\n", FormatDuration(r.MedianLeadTime))
	fmt.Fprintf(&b, "Deployment frequency: not available yet\n")
	fmt.Fprintf(&b, "Change failure rate: not available yet\n")
	fmt.Fprintf(&b, "MTTR: not available yet\n")

	return b.String()
}

// FormatDuration prints "2d 4h" / "18h 5m" / "45m" instead of Go's
// default "51h23m4.123s" - only the two biggest units, which is plenty
// for a lead time summary.
func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return "0m"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	default:
		return fmt.Sprintf("%dm", minutes)
	}
}
