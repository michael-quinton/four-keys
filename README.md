# four-keys

CLI that pulls recently-merged GitHub PRs for a repo and reports Lead
Time for Changes (one of the four DORA metrics). Deployment frequency,
change failure rate, and MTTR are stubbed as "not available yet" -
they need a notion of "deployment" this doesn't have a data source for
yet (GH Actions runs, release tags, etc).

## Running it

```bash
export GITHUB_TOKEN=ghp_yourtokenhere   # optional, but unauthenticated is 60 req/hr
go run ./cmd/fourkeys kubernetes/kubernetes
go run ./cmd/fourkeys prometheus/prometheus
```

No argument defaults to `kubernetes/kubernetes`.

```
Repository: prometheus/prometheus

Merged PRs analysed: 24
Average lead time: 2d 4h
Median lead time: 18h
Deployment frequency: not available yet
Change failure rate: not available yet
MTTR: not available yet
```

## Layout

```
cmd/fourkeys/        entry point, wires everything together
internal/github/     GitHub API client + response types
internal/domain/     app's own PullRequest model
internal/metrics/    lead time calcs
internal/report/     turns metrics into the printed report
```

`internal/github` knows GitHub's JSON shape; `internal/domain` is the
small, stable model the rest of the app works with. `domain.FromGitHub`
is where the conversion happens, and also where unmerged PRs get
filtered out.

## Not done

- No persistence - every run re-fetches from GitHub, no trend over time
- No pagination beyond the first 30 closed PRs
- No tests yet (metrics package and FromGitHub are the obvious places to start)
- Toyed with splitting this into separate services over gRPC at one
  point, decided that's overkill for a single binary doing one job
