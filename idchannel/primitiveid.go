package idchannel

import (
	"acc/logger"
	"acc/pb"
	"sync"
)

var MAXMESSAGE = 1000

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
	C  chan pb.Message
}

func (p *PIDGroup) GetRootPID(name string) *PrimitiveID {
	pid, ok := p.PIDmap.Load(name)
	if !ok {
		pid = p.createPID(name)
	}
	// type convert
	return pid.(*PrimitiveID)
}

func (p *PIDGroup) GetChildPID(cid string, pid *PrimitiveID) *PrimitiveID {
	n := pid.Id + "." + cid
	return p.GetRootPID(n)
}

func (p *PIDGroup) createPID(name string) *PrimitiveID {
	p.l.Infof("create %s channel", name)
	c := make(chan pb.Message, MAXMESSAGE)
	pid := &PrimitiveID{
		Id: name,
		C:  c,
	}
	p.PIDmap.Store(name, pid)
	return pid
}
