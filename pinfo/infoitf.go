package pinfo

import (
	"net"
	"time"
)

// BanInfo ban infos.
type BanInfo struct {
	Reason          string
	NextUseableTime time.Duration
}

// PointInfo point info.
type PointInfo struct {
	Email      string
	PubKey     string
	Status     HackerStatus
	HackerAddr *net.UDPAddr
	BanInf     *BanInfo // nil or ban.
}

// PointInfoMgrItf point info manager interface.
type PointInfoMgrItf interface {
	// HackerJoin .
	HackerJoin(udpAddr *net.UDPAddr, email, pubKey string) *PointInfo
	// FindHacker find hacker.
	FindHacker(email string) *PointInfo

	// BanEmail ban an email how many sec by what reason.
	BanEmail(email, reason string, sec int)
	// BanIp ban an ip how many sec by what reason.
	BanIP(ip, reason string, sec int)
	IsIPCanUse(ip string) (canUse bool)
	Delete(email string)
}

// PointInfoMgrFactory the factory to product PointInfoMgrItf.
type PointInfoMgrFactory func() PointInfoMgrItf

// SetPointInfoMgrFactory set PointInfoMgrItf's factory.
func SetPointInfoMgrFactory(newFactory PointInfoMgrFactory) {
	if newFactory != nil {
		pointInfoMgrFactory = newFactory
	}
}

// NewPointInfoMgr get a point info manager.
func NewPointInfoMgr() PointInfoMgrItf {
	return pointInfoMgrFactory()
}

// CreateRoomStatus .
type CreateRoomStatus int

const (
	// CRSInit ready to send info to client.
	CRSInit CreateRoomStatus = iota
	// CRSWaitPointResposne already send info to client, waiting it's response.
	CRSWaitPointResposne
	// CRSSuccess get success response.
	CRSSuccess
	// CRSRefuse refuse.
	CRSRefuse
)

// P2PRoom when help create p2p we need know eachother's PointInfo.
type P2PRoom struct {
	Creator *PointInfo
	Invitee *PointInfo
	// OffLeaseTime when time > OffLeaseTime, will destroy this room
	OffLeaseTime time.Time
	// CreatorStatus .
	CreatorStatus CreateRoomStatus
	// InviteeStatus .
	InviteeStatus CreateRoomStatus
}

// P2PHelperItf point to point helper, to help point connect another point.
type P2PHelperItf interface {
	// CreateRoom creator invite invitee create p2p connect, if room exist, will return the existed room.
	CreateRoom(creator, invitee *PointInfo) (room *P2PRoom, err error)
	// DestoryRoom if exist will destroy room, or do nothing.
	DestroyRoom(creator, invitee *PointInfo)
}

// P2PHelperFacotry create p2p healper.
type P2PHelperFacotry func() P2PHelperItf

// SetP2PHelperFacotry .
func SetP2PHelperFacotry(factory P2PHelperFacotry) {
	if factory != nil {
		p2PHelperFacotry = factory
	}
}

// NewP2PHelper .
func NewP2PHelper() P2PHelperItf {
	return p2PHelperFacotry()
}
