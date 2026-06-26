// Package github holds the types and HTTP client for talking to the
// GitHub REST API. It doesn't know anything about DORA metrics or lead
// time - just fetching PRs and decoding the JSON.
package github

// PullRequestResponse mirrors the bits of GitHub's PR JSON we actually
// use. Anything without a json tag here just gets dropped on decode.
type PullRequestResponse struct {
	Number int    `json:"number"`
	Title  string `json:"title"`

	CreatedAt string  `json:"created_at"`
	MergedAt  *string `json:"merged_at"` // nil if closed without merging

	MergeCommitSHA string `json:"merge_commit_sha"`

	User User `json:"user"`
	Base Base `json:"base"`
}

type User struct {
	Login string `json:"login"`
}

type Base struct {
	Ref  string `json:"ref"`
	Repo Repo   `json:"repo"`
}

type Repo struct {
	FullName string `json:"full_name"`
}
