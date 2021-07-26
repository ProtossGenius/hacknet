package pinfo

// pointInfoMgrByMemory PointInfoMgrItf's memory implement.
type pointInfoMgrByMemory struct {
	infoMap map[string]*PointInfo
}

func (p *pointInfoMgrByMemory) Find(email string) *PointInfo {
	return p.infoMap[email]
}

func (p *pointInfoMgrByMemory) Register(email string, IP string, port int) *PointInfo {
	info := &PointInfo{IP: IP, Port: port, Email: email, Status: AliveStatusLive}
	p.infoMap[email] = info
	return info
}

// Ban ban an email how many sec by what reason.
func (p *pointInfoMgrByMemory) Ban(email string, reason string, sec int) {
	panic("not implemented") // TODO: Implement
}

func (p *pointInfoMgrByMemory) Delete(email string) {
	panic("not implemented") // TODO: Implement
}

func newMemoryMgr() *pointInfoMgrByMemory {
	res := &pointInfoMgrByMemory{}

	return res
}
