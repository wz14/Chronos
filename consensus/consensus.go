package consensus

import (
	"acc/config"
	"acc/idchannel"
	"acc/pb"
	"github.com/pkg/errors"
)

type CONSTYPE int

const (
	HBCon   CONSTYPE = 1
	CHOSCon CONSTYPE = 2
)

func Consensus(typ CONSTYPE, value *pb.Message, s config.Start) ([]*pb.Message, []*pb.TS, error) {
	pig := s.Getpig()
	pid := pig.GetRootPID(value.Id)

	acs, err := NewConsen(typ, pid, s)
	if err != nil {
		return nil, nil, errors.Errorf("create Consensus fail: %s", err.Error())
	}
	return acs.Propose(value)
}

func NewConsen(typ CONSTYPE, pid *idchannel.PrimitiveID, s config.Start) (Consen, error) {
	if typ == HBCon {
		return NewHB(pid, s)
	} else if typ == CHOSCon {
		return NewCHOS(pid, s)
	} else {
		return nil, errors.New("no such consensus type")
	}
}

type Consen interface {
	Propose(message *pb.Message) ([]*pb.Message, []*pb.TS, error)
}
