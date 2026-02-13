package domain

type Status string

const (
	StatusNotAGitRepo Status = "NOT_A_GIT_REPO"
	StatusNoRemote    Status = "NO_REMOTE"
	StatusNoUpstream  Status = "NO_UPSTREAM"
	StatusSynced      Status = "SYNCED"
	StatusSyncPending Status = "SYNC_PENDING"
	StatusLate        Status = "LATE"
	StatusDiverged    Status = "DIVERGED"
)
