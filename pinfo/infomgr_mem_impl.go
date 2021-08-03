package pinfo

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

func (p *pointInfoMgrByMemory) Register(email, ip, pubKey string, port int, status HackerStatus) *PointInfo {
	if _, exist := p.infoMap[email]; exist {
		return nil
	}

	info := &PointInfo{IP: ip, Port: port, Email: email, Status: status, BanInf: nil, PubKey: pubKey}
	p.infoMap[email] = info

	return info
}

// Ban ban an email how many sec by what reason.
func (p *pointInfoMgrByMemory) BanEmail(email string, reason string, sec int) {
	if pInfo, exist := p.infoMap[email]; exist {
		pInfo.Status = HackerStatusBan
	} else {
		p.Register(email, "", "", -1, HackerStatusBan)
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
