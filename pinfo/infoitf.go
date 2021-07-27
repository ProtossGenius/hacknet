package pinfo

import "time"

// HackerStatus hacker point's status.
type HackerStatus int

const (
	// HackerStatusLive live.
	HackerStatusLive HackerStatus = iota
	// HackerStatusDead dead.
	HackerStatusDead
	// HackerStatusBan ban.
	HackerStatusBan
)

// BanInfo ban infos.
type BanInfo struct {
	Reason          string
	NextUseableTime time.Duration
}

// PointInfo point info.
type PointInfo struct {
	Email  string
	PubKey string
	IP     string
	Port   int
	Status HackerStatus
	BanInf *BanInfo
}

// PointInfoMgrItf point info manager interface.
type PointInfoMgrItf interface {
	Find(email string) *PointInfo
	// Register return nil means has exist.
	Register(email, ip, pubKey string, port int, status HackerStatus) *PointInfo
	// BanEmail ban an email how many sec by what reason.
	BanEmail(email, reason string, sec int)
	// BanIp ban an ip how many sec by what reason.
	BanIP(ip, reason string, sec int)
	IsIPCanUse(ip string) (canUse bool)
	Delete(email string)
}

// NewPointInfoMgr get a point info manager.
func NewPointInfoMgr() PointInfoMgrItf {
	return newMemoryMgr()
}
