package idchannel

import (
	"acc/logger"
	"acc/pb"
	"strconv"
	"strings"
	"sync"
)

var MAXMESSAGE = 1000

const sep = "."

type PIDGroup struct {
	//PIDmap: string type --> *PrimitiveID type
	PIDmap sync.Map
	l      *logger.Logger
}

func NewPIDGroup(c Config) (*PIDGroup, error) {
	myid, _ := c.GetMyID()
	pidg := PIDGroup{
		PIDmap: sync.Map{},
		l:      logger.NewLoggerWithID("idchannel", myid),
	}
	return &pidg, nil
}

type PrimitiveID struct {
	Id string
	C  chan *pb.Message
}

func (p *PIDGroup) GetRootPID(name string) *PrimitiveID {
	pid, loaded := p.PIDmap.LoadOrStore(name, p.createPID(name))
	if loaded {
		p.l.Infof("create %s channel", name)
	} else {
		p.l.Infof("create %s channel but exist", name)
	}
	// type convert
	return pid.(*PrimitiveID)
}

func (p *PIDGroup) GetChildPID(cid string, pid *PrimitiveID) *PrimitiveID {
	n := pid.Id + sep + cid
	return p.GetRootPID(n)
}

func (p *PIDGroup) GetParentPID(pid *PrimitiveID) *PrimitiveID {
	list := strings.Split(pid.Id, sep)
	n := strings.Join(list[:len(list)-1], sep)
	return p.GetRootPID(n)
}

func (p *PIDGroup) GetInitRoundPID(pid *PrimitiveID) *PrimitiveID {
	n := pid.Id + sep + "1"
	return p.GetRootPID(n)
}

func (p *PIDGroup) GetNextRoundPID(pid *PrimitiveID) *PrimitiveID {
	list := strings.Split(pid.Id, sep)
	lastIndex := len(list) - 1
	r, _ := strconv.Atoi(list[lastIndex])
	list[lastIndex] = strconv.Itoa(r + 1)
	n := strings.Join(list, sep)
	return p.GetRootPID(n)
}

func (p *PIDGroup) createPID(name string) *PrimitiveID {
	c := make(chan *pb.Message, MAXMESSAGE)
	pid := &PrimitiveID{
		Id: name,
		C:  c,
	}
	return pid
}
