package pinfo

// HackerStatus hacker point's status.
type HackerStatus int

const (
	// AliveStatusLive live.
	AliveStatusLive HackerStatus = iota
	// AliveStatusDead dead.
	AliveStatusDead
)

// PointInfo point info.
type PointInfo struct {
	Email  string
	PubKey string
	IP     string
	Port   int
	Status HackerStatus
}

// PointInfoMgrItf point info manager interface.
type PointInfoMgrItf interface {
	Find(email string) PointInfo
	Register(email, IP string, port int) PointInfo
	// Ban ban an email how many sec by what reason.
	Ban(email, reason string, sec int)
	Delete(email string)
}

// NewPointInfoMgr get a point info manager.
func NewPointInfoMgr() PointInfoMgrItf {
	return newMemoryMgr()
}
