package pinfo

import "net"

// pointInfoMgrByMemory PointInfoMgrItf's memory implement.
type pointInfoMgrByMemory struct {
	infoMap   map[string]*PointInfo
	ipBans    map[string]*BanInfo
	emailBans map[string]*BanInfo
}

// BanIp ban an ip how many sec by what reason.
func (p *pointInfoMgrByMemory) BanIP(ip string, reason string, sec int) {
}

func (p *pointInfoMgrByMemory) IsIPCanUse(ip string) (canUse bool) {
	return true
}

func (p *pointInfoMgrByMemory) Find(email string) *PointInfo {
	return p.infoMap[email]
}

func (p *pointInfoMgrByMemory) HackerJoin(hackerAddr *net.UDPAddr,
	email, pubKey string) *PointInfo {
	if info, exist := p.infoMap[email]; exist {
		return info
	}

	info := &PointInfo{
		HackerAddr: hackerAddr, Email: email, Status: HackerStatusLive,
		BanInf: nil, PubKey: pubKey,
	}
	p.infoMap[email] = info

	return info
}

// Ban ban an email how many sec by what reason.
func (p *pointInfoMgrByMemory) BanEmail(email string, reason string, sec int) {
	if pInfo, exist := p.infoMap[email]; exist {
		pInfo.Status = HackerStatusBan
	} else {
		pInfo = p.HackerJoin(nil, email, "")
		pInfo.Status = HackerStatusBan
	}
}

func (p *pointInfoMgrByMemory) Delete(email string) {
}

func newMemoryMgr() PointInfoMgrItf {
	res := &pointInfoMgrByMemory{
		infoMap:   make(map[string]*PointInfo),
		ipBans:    make(map[string]*BanInfo),
		emailBans: make(map[string]*BanInfo),
	}

	return res
}
