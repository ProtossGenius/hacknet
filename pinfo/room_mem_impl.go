package pinfo

import "time"

// p2pHelperMemImpl p2p helper memory impl version.
type p2pHelperMemImpl struct {
	roomMap map[string]*P2PRoom
}

func (p *p2pHelperMemImpl) getKey(creator, invitee *PointInfo) string {
	return creator.Email + ">>" + invitee.Email
}

const (
	// LeaseTime how long room can quiet.
	LeaseTime = 10 * time.Minute
)

// CreateRoom creator invite invitee create p2p connect, if room exist, will return the existed room.
func (p *p2pHelperMemImpl) CreateRoom(creator *PointInfo, invitee *PointInfo) (room *P2PRoom, err error) {
	key := p.getKey(creator, invitee)
	// TODO: create fail about.
	// TODO: not thread safe.
	if room, exist := p.roomMap[key]; exist && room != nil {
		room.OffLeaseTime.Add(LeaseTime)

		return room, nil
	}

	room = &P2PRoom{
		Creator:       creator,
		Invitee:       invitee,
		OffLeaseTime:  time.Now().Add(LeaseTime),
		CreatorStatus: CRSInit,
		InviteeStatus: CRSInit,
	}
	p.roomMap[key] = room

	return room, nil
}

// DestoryRoom if exist will destroy room, or do nothing.
func (p *p2pHelperMemImpl) DestroyRoom(creator *PointInfo, invitee *PointInfo) {
	delete(p.roomMap, p.getKey(creator, invitee))
}

func newP2pHelperMemoryImpl() P2PHelperItf {
	return &p2pHelperMemImpl{roomMap: make(map[string]*P2PRoom)}
}
