package domain

type Result struct {
	RepoPath             string
	Branch               string
	Upstream             string
	Status               Status
	Behind               int
	Ahead                int
	Flags                []string
	Actions              []string
	Details              []string
	Err                  string
	NOUpstreamWasMerged  bool
	NOUpstreamMergeBase  string
	NOUpstreamSuggestion string
}

func (r Result) HasFlag(flag string) bool {
	for _, f := range r.Flags {
		if f == flag {
			return true
		}
	}
	return false
}
